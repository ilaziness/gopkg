package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/firebase/genkit/go/plugins/server"
)

// AgentInput - Agent 对话输入
type AgentInput struct {
	Message string `json:"message" jsonschema:"description=用户消息"`
}

// AgentOutput - Agent 输出
type AgentOutput struct {
	Response string `json:"response" jsonschema:"description=Agent 回复"`
}

// WeatherInput - 天气工具输入
type WeatherInput struct {
	Location string `json:"location" jsonschema:"description=地点"`
}

func main() {
	ctx := context.Background()

	ollamaPlugin := &ollama.Ollama{
		ServerAddress: "http://127.0.0.1:11434",
		Timeout:       80,
	}

	g := genkit.Init(ctx,
		genkit.WithPlugins(ollamaPlugin),
	)
	model := ollamaPlugin.DefineModel(
		g,
		ollama.ModelDefinition{
			Name: "gemma4:e2b",
			Type: "generate",
		},
		&ai.ModelOptions{
			Supports: &ai.ModelSupports{
				Multiturn:  true,
				SystemRole: true,
				Tools:      true,
				Media:      false,
			},
		},
	)

	// ====================================================================
	// Agent 模式：带工具的聊天式 AI 助手
	// Agent = LLM + 工具 + 多轮对话能力
	// ====================================================================

	// 定义工具：天气查询
	getWeatherTool := genkit.DefineTool(
		g,
		"getWeather",
		"获取指定地点的当前天气信息。",
		func(ctx *ai.ToolContext, input WeatherInput) (string, error) {
			return fmt.Sprintf("%s: 22°C, 多云", input.Location), nil
		},
	)

	// Agent Flow：简单的 ReAct 模式 Agent
	// 使用系统提示词定义 Agent 的行为角色
	agentFlow := genkit.DefineFlow(g, "agentFlow",
		func(ctx context.Context, input *AgentInput) (*AgentOutput, error) {
			resp, err := genkit.Generate(ctx, g,
				ai.WithModel(model),
				ai.WithSystem("你是一个智能助手（Agent）。你有工具可以使用。"+
					"当你需要实时信息时，使用工具获取数据。"+
					"思考过程：分析需求 -> 使用工具(如果需要) -> 组织回答。使用中文。"),
				ai.WithPrompt(input.Message),
				ai.WithTools(getWeatherTool),
			)
			if err != nil {
				return nil, fmt.Errorf("agent flow failed: %w", err)
			}
			return &AgentOutput{Response: resp.Text()}, nil
		})

	// 多轮对话 Agent（含历史记录）
	chatAgentFlow := genkit.DefineFlow(g, "chatAgentFlow",
		func(ctx context.Context, input *AgentInput) (*AgentOutput, error) {
			// 模拟历史对话
			messages := []*ai.Message{
				ai.NewUserMessage("你好！"),
				ai.NewModelTextMessage("你好！我是你的智能助手，有什么可以帮你的吗？"),
				ai.NewUserMessage(input.Message),
			}

			resp, err := genkit.Generate(ctx, g,
				ai.WithModel(model),
				ai.WithMessages(messages...),
				ai.WithSystem("你是一个有帮助的 AI 助手。使用工具获取实时数据。使用中文。"),
				ai.WithTools(getWeatherTool),
			)
			if err != nil {
				return nil, fmt.Errorf("chat agent failed: %w", err)
			}
			return &AgentOutput{Response: resp.Text()}, nil
		})

	fmt.Println("========== Agent Flow 测试 ==========")
	result, err := agentFlow.Run(ctx, &AgentInput{Message: "北京的天气怎么样？后天适合出门吗？"})
	if err != nil {
		log.Printf("agentFlow error: %v", err)
	} else {
		b, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(b))
	}

	fmt.Println("\n========== 多轮对话 Agent 测试 ==========")
	result2, err := chatAgentFlow.Run(ctx, &AgentInput{Message: "帮我查一下上海的天气"})
	if err != nil {
		log.Printf("chatAgentFlow error: %v", err)
	} else {
		b, _ := json.MarshalIndent(result2, "", "  ")
		fmt.Println(string(b))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /agentFlow", genkit.Handler(agentFlow))
	mux.HandleFunc("POST /chatAgentFlow", genkit.Handler(chatAgentFlow))

	log.Println("Starting server on http://localhost:3400")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}