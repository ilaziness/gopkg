package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/firebase/genkit/go/plugins/server"
)

type ChatInput struct {
	Message string `json:"message" jsonschema:"description=用户消息"`
}
type ChatOutput struct {
	Response string `json:"response" jsonschema:"description=AI 回复"`
}

func main() {
	ctx := context.Background()

	ollamaPlugin := &ollama.Ollama{
		ServerAddress: "http://127.0.0.1:11434",
		Timeout:       80,
	}

	g := genkit.Init(ctx, genkit.WithPlugins(ollamaPlugin))
	model := ollamaPlugin.DefineModel(g, ollama.ModelDefinition{
		Name: "gemma4:e2b", Type: "generate",
	}, &ai.ModelOptions{Supports: &ai.ModelSupports{Multiturn: true, SystemRole: true, Tools: false, Media: false}})

	// 持久化聊天会话 - 模拟维护对话历史
	sessionHistory := make(map[string][]*ai.Message)

	chatFlow := genkit.DefineFlow(g, "chatFlow",
		func(ctx context.Context, input *ChatInput) (*ChatOutput, error) {
			sessionID := "default"
			history := sessionHistory[sessionID]

			// 将用户消息加入历史
			history = append(history, ai.NewUserMessage(input.Message))

			resp, err := genkit.Generate(ctx, g,
				ai.WithModel(model),
				ai.WithMessages(history...),
				ai.WithSystem("你是一个有帮助的聊天助手。记住对话上下文。使用中文。"),
			)
			if err != nil {
				return nil, fmt.Errorf("chat failed: %w", err)
			}

			// 将模型回复加入历史
			history = append(history, ai.NewModelTextMessage(resp.Text()))
			sessionHistory[sessionID] = history

			return &ChatOutput{Response: resp.Text()}, nil
		})

	fmt.Println("========== 持久化聊天会话 ==========")
	fmt.Println("该示例演示了使用 ai.WithMessages() 维护多轮对话上下文。")
	fmt.Println("每轮对话都会将历史消息传入，模型能理解上下文。")

	result, err := chatFlow.Run(ctx, &ChatInput{Message: "你好！我叫小明"})
	if err != nil {
		log.Printf("chatFlow error: %v", err)
	} else {
		fmt.Println("用户: 你好！我叫小明")
		fmt.Printf("AI: %s\n\n", result.Response)
	}

	result2, err := chatFlow.Run(ctx, &ChatInput{Message: "我叫什么名字？"})
	if err != nil {
		log.Printf("chatFlow error: %v", err)
	} else {
		fmt.Println("用户: 我叫什么名字？")
		fmt.Printf("AI: %s\n", result2.Response)
	}

	fmt.Println("\n========== Context (上下文) ==========")
	fmt.Println("ai.WithDocs() 可向模型提供参考文档（RAG 场景）：")
	fmt.Println(`  ai.WithDocs(ai.NewTextPart("参考文档内容"))`)
	fmt.Println("")
	fmt.Println("========== Interrupts (中断) ==========")
	fmt.Println("使用 ai.WithReturnToolRequests(true) 可在工具调用前暂停，"+
		"等待用户确认后再继续执行。")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /chatFlow", genkit.Handler(chatFlow))

	log.Println("Starting server on http://localhost:3400")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}