package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/firebase/genkit/go/plugins/server"
)

// ============================================================================
// Schema 定义
// ============================================================================

// MenuSuggestionInput - 菜单建议 Flow 输入
type MenuSuggestionInput struct {
	Theme string `json:"theme" jsonschema:"description=餐厅主题风格"`
}

// MenuItem - 菜单项
type MenuItem struct {
	Name        string `json:"name" jsonschema:"description=菜品名称"`
	Description string `json:"description" jsonschema:"description=菜品描述"`
	Price       string `json:"price" jsonschema:"description=价格"`
}

// FormattedMenuOutput - 格式化后的菜单输出
type FormattedMenuOutput struct {
	FormattedMenuItem string `json:"formattedMenuItem" jsonschema:"description=格式化的菜单项"`
}

// ReviewInput - 评论分析 Flow 输入
type ReviewInput struct {
	ReviewText string `json:"reviewText" jsonschema:"description=用户评价文本"`
}

// ReviewOutput - 评论分析输出
type ReviewOutput struct {
	Sentiment  string `json:"sentiment" jsonschema:"description=情感倾向: 正面/负面/中性"`
	Rating     int    `json:"rating" jsonschema:"description=评分 1-5"`
	Summary    string `json:"summary" jsonschema:"description=评价摘要"`
	ReplySuggestion string `json:"replySuggestion,omitempty" jsonschema:"description=回复建议"`
}

// JokeInput - 笑话生成 Flow 输入
type JokeInput struct {
	Topic string `json:"topic" jsonschema:"description=笑话主题"`
}

// ============================================================================
// main
// ============================================================================

func main() {
	ctx := context.Background()

	// ---------- 1. 初始化 Genkit 和模型 ----------
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


	// ---------- 2. 定义各种 Flow ----------

	// --- 2a. 基本 Flow：菜单建议 ---
	// 演示 genkit.DefineFlow + genkit.GenerateData[T]
	menuSuggestionFlow := genkit.DefineFlow(g, "menuSuggestionFlow",
		func(ctx context.Context, input *MenuSuggestionInput) (*MenuItem, error) {
			item, _, err := genkit.GenerateData[MenuItem](ctx, g,
				ai.WithSystem("使用简体中文回答。"),
				ai.WithPrompt("为一个%s主题的餐厅发明一道新菜品，包含名称、描述和价格。", input.Theme),
				ai.WithModel(model),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to generate menu item: %w", err)
			}
			return item, nil
		})

	// --- 2b. 结构化输出 Flow：格式化菜单建议 ---
	// 演示 Flow 的输出 schema 可以不同于内部模型调用的 schema
	menuSuggestionMarkdownFlow := genkit.DefineFlow(g, "menuSuggestionMarkdownFlow",
		func(ctx context.Context, input *MenuSuggestionInput) (*FormattedMenuOutput, error) {
			// 第一步：生成结构化的菜单项
			item, _, err := genkit.GenerateData[MenuItem](ctx, g,
				ai.WithSystem("使用简体中文回答。"),
				ai.WithPrompt("为一个%s主题的餐厅发明一道新菜品。", input.Theme),
				ai.WithModel(model),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to generate menu item: %w", err)
			}

			// 第二步：将结构化数据格式化为 Markdown 文本
			return &FormattedMenuOutput{
				FormattedMenuItem: fmt.Sprintf("**%s** - %s\n价格: %s", item.Name, item.Description, item.Price),
			}, nil
		})

	// --- 2c. Streaming Flow：笑话生成（迭代器模式）---
	// 演示 genkit.DefineStreamingFlow + genkit.GenerateStream
	jokeStreamFlow := genkit.DefineStreamingFlow(g, "jokeStreamFlow",
		func(ctx context.Context, input *JokeInput, sendChunk core.StreamCallback[string]) (string, error) {
			prompt := fmt.Sprintf("给我讲一个关于%s的笑话，要幽默风趣。", input.Topic)

			stream := genkit.GenerateStream(ctx, g,
				ai.WithModel(model),
				ai.WithPrompt(prompt),
			)

			var fullText string
			for result, err := range stream {
				if err != nil {
					return "", fmt.Errorf("stream error: %w", err)
				}
				if result.Done {
					fullText = result.Response.Text()
				} else {
					// 将每个 chunk 发送到 Flow 的输出流
					if err := sendChunk(ctx, result.Chunk.Text()); err != nil {
						return "", fmt.Errorf("send chunk failed: %w", err)
					}
				}
			}
			return fullText, nil
		})

	// --- 2d. 带 Flow Steps 的评论分析 Flow ---
	// 演示 genkit.Run() 创建可追踪的步骤
	reviewAnalysisFlow := genkit.DefineFlow(g, "reviewAnalysisFlow",
		func(ctx context.Context, input *ReviewInput) (*ReviewOutput, error) {
			// Step 1: 情感分析（使用 genkit.Run 包装，在 trace 中可见）
			sentiment, err := genkit.Run(ctx, "analyze-sentiment", func() (string, error) {
				resp, err := genkit.Generate(ctx, g,
					ai.WithModel(model),
					ai.WithSystem("只返回一个词：正面、负面或中性。"),
					ai.WithPrompt(fmt.Sprintf("分析以下评价的情感倾向：\n%s", input.ReviewText)),
				)
				if err != nil {
					return "", err
				}
				return resp.Text(), nil
			})
			if err != nil {
				return nil, fmt.Errorf("sentiment analysis failed: %w", err)
			}

			// Step 2: 提取评分
			rating, err := genkit.Run(ctx, "extract-rating", func() (int, error) {
				resp, err := genkit.Generate(ctx, g,
					ai.WithModel(model),
					ai.WithSystem("根据评价内容给出 1-5 的评分。只返回数字。"),
					ai.WithPrompt(fmt.Sprintf("评价：%s", input.ReviewText)),
				)
				if err != nil {
					return 0, err
				}
				var rating int
				fmt.Sscanf(resp.Text(), "%d", &rating)
				if rating < 1 {
					rating = 1
				}
				if rating > 5 {
					rating = 5
				}
				return rating, nil
			})
			if err != nil {
				return nil, fmt.Errorf("rating extraction failed: %w", err)
			}

			// Step 3: 生成摘要和回复建议
			partial, _, err := genkit.GenerateData[ReviewOutput](ctx, g,
				ai.WithModel(model),
				ai.WithSystem("根据用户评价生成摘要和回复建议。使用简体中文。"),
				ai.WithPrompt(fmt.Sprintf(
					"评价：%s\n情感：%s\n评分：%d",
					input.ReviewText, sentiment, rating,
				)),
			)
			if err != nil {
				return nil, fmt.Errorf("summary generation failed: %w", err)
			}

			// 填充 Step 1 和 Step 2 的结果
			partial.Sentiment = sentiment
			partial.Rating = rating

			return partial, nil
		})

	// ---------- 3. 本地运行测试 ----------

	// 测试 1: 菜单建议 Flow
	fmt.Println("===== 测试 1: 菜单建议 Flow =====")
	menuItem, err := menuSuggestionFlow.Run(ctx, &MenuSuggestionInput{Theme: "科幻"})
	if err != nil {
		log.Printf("menuSuggestionFlow error: %v", err)
	} else {
		b, _ := json.MarshalIndent(menuItem, "", "  ")
		fmt.Printf("菜单项:\n%s\n\n", string(b))
	}

	// 测试 2: 格式化菜单 Flow
	fmt.Println("===== 测试 2: 格式化菜单 Flow =====")
	formatted, err := menuSuggestionMarkdownFlow.Run(ctx, &MenuSuggestionInput{Theme: "太空"})
	if err != nil {
		log.Printf("menuSuggestionMarkdownFlow error: %v", err)
	} else {
		fmt.Printf("格式化输出:\n%s\n\n", formatted.FormattedMenuItem)
	}

	// 测试 3: 笑话 Streaming Flow
	fmt.Println("===== 测试 3: 笑话 Streaming Flow =====")
	streamCh := jokeStreamFlow.Stream(ctx, &JokeInput{Topic: "程序员"})
	fmt.Print("笑话内容: ")
	for result := range streamCh {
		if result.Done {
			fmt.Printf("\n(流式输出完成)\n\n")
		} else {
			fmt.Print(result.Stream)
		}
	}

	// 测试 4: 评论分析 Flow（带 Flow Steps）
	fmt.Println("===== 测试 4: 评论分析 Flow =====")
	reviewResult, err := reviewAnalysisFlow.Run(ctx, &ReviewInput{
		ReviewText: "菜品味道不错，但上菜速度太慢了，等了快一个小时。环境倒是很好。",
	})
	if err != nil {
		log.Printf("reviewAnalysisFlow error: %v", err)
	} else {
		b, _ := json.MarshalIndent(reviewResult, "", "  ")
		fmt.Printf("评论分析:\n%s\n\n", string(b))
	}

	// ---------- 4. 启动 HTTP 服务部署所有 Flow ----------
	fmt.Println("===== 部署所有 Flow =====")

	mux := http.NewServeMux()

	// 自动注册所有已定义的 Flow 为 HTTP 端点
	for _, flow := range genkit.ListFlows(g) {
		path := fmt.Sprintf("POST /%s", flow.Name())
		mux.HandleFunc(path, genkit.Handler(flow))
		fmt.Printf("  -> %s\n", path)
	}

	fmt.Println()
	log.Println("Starting server on http://localhost:3400")
	log.Println("Deployed flow endpoints:")

	fmt.Println("\n可用测试命令:")
	fmt.Println(`  curl -X POST "http://localhost:3400/menuSuggestionFlow" ` +
		`-H "Content-Type: application/json" -d '{"data": {"theme": "法国菜"}}'`)
	fmt.Println(`  curl -X POST "http://localhost:3400/reviewAnalysisFlow" ` +
		`-H "Content-Type: application/json" -d '{"data": {"reviewText": "非常好吃！"}}'`)
	fmt.Println(`  curl -X POST "http://localhost:3400/jokeStreamFlow" ` +
		`-H "Content-Type: application/json" -H "Accept: text/event-stream" ` +
		`-d '{"data": {"topic": "程序员"}}'`)

	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}