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

// ============================================================================
// Schema 定义 - 结构化输入/输出类型
// ============================================================================

// MenuItemInput - 菜谱生成请求
type MenuItemInput struct {
	Theme   string `json:"theme" jsonschema:"description=餐厅主题"`
	Cuisine string `json:"cuisine,omitempty" jsonschema:"description=菜系类型"`
}

// MenuItem - 菜品结构（作为 GenerateData[T] 的泛型参数）
type MenuItem struct {
	Name        string   `json:"name" jsonschema:"description=菜品名称"`
	Description string   `json:"description" jsonschema:"description=菜品描述"`
	Ingredients []string `json:"ingredients" jsonschema:"description=食材列表"`
	Price       string   `json:"price" jsonschema:"description=价格"`
	SpiceLevel  string   `json:"spiceLevel,omitempty" jsonschema:"description=辣度等级"`
}

func main() {
	ctx := context.Background()

	// 1. 初始化 Genkit 和 Ollama 插件
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
				Tools:      false,
				Media:      false,
			},
		},
	)

	// 示例 1: 基础文本生成 - genkit.Generate()
	// 最简单的使用方式：指定模型 + 提示词，获取纯文本响应
	basicGenerateFlow := genkit.DefineFlow(g, "basicGenerateFlow",
		func(ctx context.Context, input *MenuItemInput) (string, error) {
			resp, err := genkit.Generate(ctx, g,
				ai.WithModel(model),
				ai.WithSystem("你是一位创意美食顾问。使用简体中文回答。"),
				ai.WithPrompt("为一个%s主题的餐厅发明一道创意菜品，菜系：%s。", input.Theme, input.Cuisine),
			)
			if err != nil {
				return "", fmt.Errorf("generate failed: %w", err)
			}
			return resp.Text(), nil
		})

	// 示例 2: 结构化输出 - genkit.GenerateData[T]()
	// 指定 Go struct 作为泛型参数，模型输出自动解析为结构化数据
	// 返回值: (*T, *ai.ModelResponse, error)
	structuredOutputFlow := genkit.DefineFlow(g, "structuredOutputFlow",
		func(ctx context.Context, input *MenuItemInput) (*MenuItem, error) {
			item, _, err := genkit.GenerateData[MenuItem](ctx, g,
				ai.WithModel(model),
				ai.WithSystem("你是一位创意美食顾问。严格按照要求的 JSON 格式输出。使用简体中文。"),
				ai.WithPrompt("为一个%s主题的餐厅设计一道菜品，菜系：%s。返回名称、描述、食材列表、价格。",
					input.Theme, input.Cuisine),
			)
			if err != nil {
				return nil, fmt.Errorf("generate data failed: %w", err)
			}
			return item, nil
		})

	// 示例 3: 系统提示词 - ai.WithSystem()
	// 设置模型的行为角色、语气和输出格式
	systemPromptFlow := genkit.DefineFlow(g, "systemPromptFlow",
		func(ctx context.Context, input *MenuItemInput) (*MenuItem, error) {
			item, _, err := genkit.GenerateData[MenuItem](ctx, g,
				ai.WithModel(model),
				ai.WithSystem("你是米其林三星餐厅的主厨，回复要专业但带幽默感，使用英文。"+
					"输出严格的 JSON 格式包含 name, description, ingredients, price 字段。"),
				ai.WithPrompt("Create a signature dish for a %s themed restaurant.", input.Theme),
			)
			if err != nil {
				return nil, fmt.Errorf("system prompt demo failed: %w", err)
			}
			return item, nil
		})

	// 示例 4: 模型参数控制 - ai.WithConfig()
	// 通过配置控制输出的创造性和长度
	modelParamsFlow := genkit.DefineFlow(g, "modelParamsFlow",
		func(ctx context.Context, input *MenuItemInput) (*MenuItem, error) {
			item, _, err := genkit.GenerateData[MenuItem](ctx, g,
				ai.WithModel(model),
				ai.WithSystem("严格输出 JSON 格式。使用简体中文。"),
				ai.WithPrompt("为一个%s主题的餐厅设计一道菜品。", input.Theme),
				ai.WithConfig(&ai.GenerationCommonConfig{
					Temperature:     0.8,
					MaxOutputTokens: 500,
					StopSequences:   []string{"</end>"},
				}),
			)
			if err != nil {
				return nil, fmt.Errorf("model params demo failed: %w", err)
			}
			return item, nil
		})

	// 示例 5: 流式生成 - genkit.GenerateStream()
	// 返回迭代器，逐块输出生成内容
	streamingGenerateFlow := genkit.DefineFlow(g, "streamingGenerateFlow",
		func(ctx context.Context, input *MenuItemInput) (string, error) {
			var fullText string
			for result, err := range genkit.GenerateStream(ctx, g,
				ai.WithModel(model),
				ai.WithSystem("使用中文回答。"),
				ai.WithPrompt("写一段关于%s主题餐厅的创意菜品描述，200字左右。", input.Theme),
			) {
				if err != nil {
					return "", fmt.Errorf("stream error: %w", err)
				}
				if result.Done {
					fullText = result.Response.Text()
					fmt.Println("\n[流式生成完成]")
				} else {
					fmt.Print(result.Chunk.Text())
				}
			}
			return fullText, nil
		})

	// 示例 6: 结构化流式输出 - genkit.GenerateDataStream[T]()
	structuredStreamFlow := genkit.DefineFlow(g, "structuredStreamFlow",
		func(ctx context.Context, input *MenuItemInput) (*MenuItem, error) {
			for result, err := range genkit.GenerateDataStream[*MenuItem](ctx, g,
				ai.WithModel(model),
				ai.WithSystem("输出 JSON 格式，包含 name, description, ingredients, price。使用简体中文。"),
				ai.WithPrompt("为一个%s主题的餐厅设计一道美食。", input.Theme),
			) {
				if err != nil {
					return nil, fmt.Errorf("structured stream error: %w", err)
				}
				if result.Done {
					return result.Output, nil
				}
				fmt.Printf("收到部分结构化数据: %+v\n", result.Chunk)
			}
			return nil, fmt.Errorf("unexpected end of stream")
		})

	// 示例 7: 便捷函数 genkit.GenerateText()
	// 直接返回纯文本字符串，适合简单场景
	simpleTextFlow := genkit.DefineFlow(g, "simpleTextFlow",
		func(ctx context.Context, input *MenuItemInput) (string, error) {
			text, err := genkit.GenerateText(ctx, g,
				ai.WithModel(model),
				ai.WithPrompt("用一句话描述%s主题餐厅的招牌菜品。", input.Theme),
			)
			if err != nil {
				return "", fmt.Errorf("generate text failed: %w", err)
			}
			return text, nil
		})

	// 本地测试：运行所有 Flow 验证功能
	fmt.Println("========== 示例 1: 基础文本生成 ==========")
	text, err := basicGenerateFlow.Run(ctx, &MenuItemInput{Theme: "太空", Cuisine: "分子料理"})
	if err != nil {
		log.Printf("basicGenerateFlow error: %v", err)
	} else {
		fmt.Println(text)
	}

	fmt.Println("\n========== 示例 2: 结构化输出 ==========")
	item, err := structuredOutputFlow.Run(ctx, &MenuItemInput{Theme: "森林", Cuisine: "素食"})
	if err != nil {
		log.Printf("structuredOutputFlow error: %v", err)
	} else {
		b, _ := json.MarshalIndent(item, "", "  ")
		fmt.Println(string(b))
	}

	fmt.Println("\n========== 示例 3: 系统提示词 ==========")
	item2, err := systemPromptFlow.Run(ctx, &MenuItemInput{Theme: "Underwater"})
	if err != nil {
		log.Printf("systemPromptFlow error: %v", err)
	} else {
		b, _ := json.MarshalIndent(item2, "", "  ")
		fmt.Println(string(b))
	}

	fmt.Println("\n========== 示例 7: 简短便捷文本 ==========")
	short, err := simpleTextFlow.Run(ctx, &MenuItemInput{Theme: "赛博朋克"})
	if err != nil {
		log.Printf("simpleTextFlow error: %v", err)
	} else {
		fmt.Println(short)
	}

	// 部署：启动 HTTP 服务
	mux := http.NewServeMux()
	mux.HandleFunc("POST /basicGenerateFlow", genkit.Handler(basicGenerateFlow))
	mux.HandleFunc("POST /structuredOutputFlow", genkit.Handler(structuredOutputFlow))
	mux.HandleFunc("POST /systemPromptFlow", genkit.Handler(systemPromptFlow))
	mux.HandleFunc("POST /modelParamsFlow", genkit.Handler(modelParamsFlow))
	mux.HandleFunc("POST /streamingGenerateFlow", genkit.Handler(streamingGenerateFlow))
	mux.HandleFunc("POST /structuredStreamFlow", genkit.Handler(structuredStreamFlow))
	mux.HandleFunc("POST /simpleTextFlow", genkit.Handler(simpleTextFlow))

	fmt.Println("\n========== 已部署的 Flow 端点 ==========")
	for _, flow := range genkit.ListFlows(g) {
		fmt.Printf("  POST /%s\n", flow.Name())
	}
	fmt.Println(`测试命令示例:`)
	fmt.Println(`  curl -X POST "http://localhost:3400/structuredOutputFlow" -H "Content-Type: application/json" -d '{"data": {"theme": "海洋", "cuisine": "海鲜"}}'`)

	log.Println("Starting server on http://localhost:3400")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}