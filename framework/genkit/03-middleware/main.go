package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/middleware"
	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/firebase/genkit/go/plugins/server"
)

// ============================================================================
// Schema 定义
// ============================================================================

// RequestInput - 通用的生成请求
type RequestInput struct {
	Topic string `json:"topic" jsonschema:"description=请求主题"`
}

// JokeOutput - 笑话输出
type JokeOutput struct {
	Joke string `json:"joke" jsonschema:"description=生成的笑话"`
}

// ============================================================================
// 自定义中间件示例
// ============================================================================

// LoggingMiddleware - 记录模型调用耗时的自定义中间件
type LoggingMiddleware struct {
	Prefix string `json:"prefix,omitempty"`
}

func (LoggingMiddleware) Name() string { return "mydemo/logger" }

func (l LoggingMiddleware) New(ctx context.Context) (*ai.Hooks, error) {
	return &ai.Hooks{
		WrapModel: func(ctx context.Context, p *ai.ModelParams, next ai.ModelNext) (*ai.ModelResponse, error) {
			start := time.Now()
			resp, err := next(ctx, p)
			log.Printf("%s 模型调用耗时: %s", l.Prefix, time.Since(start))
			return resp, err
		},
	}, nil
}

// CounterMiddleware - 统计模型调用次数的自定义中间件（展示状态共享）
type CounterMiddleware struct{}

func (CounterMiddleware) Name() string { return "mydemo/counter" }

func (CounterMiddleware) New(ctx context.Context) (*ai.Hooks, error) {
	var modelCalls int
	return &ai.Hooks{
		WrapModel: func(ctx context.Context, p *ai.ModelParams, next ai.ModelNext) (*ai.ModelResponse, error) {
			modelCalls++
			return next(ctx, p)
		},
		WrapGenerate: func(ctx context.Context, p *ai.GenerateParams, next ai.GenerateNext) (*ai.ModelResponse, error) {
			resp, err := next(ctx, p)
			log.Printf("迭代 %d: 已进行了 %d 次模型调用", p.Iteration, modelCalls)
			return resp, err
		},
	}, nil
}

func main() {
	ctx := context.Background()

	// 1. 初始化 Genkit
	ollamaPlugin := &ollama.Ollama{
		ServerAddress: "http://127.0.0.1:11434",
		Timeout:       80,
	}

	g := genkit.Init(ctx,
		genkit.WithPlugins(ollamaPlugin, &middleware.Middleware{}),
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

	// ====================================================================
	// 示例 1: Retry 中间件 - 自动重试失败的模型调用
	// 当模型返回临时错误（RESOURCE_EXHAUSTED, UNAVAILABLE 等）时自动重试
	// 使用指数退避 + jitter，默认最多重试 3 次
	// ====================================================================
	retryFlow := genkit.DefineFlow(g, "retryFlow",
		func(ctx context.Context, input *RequestInput) (*JokeOutput, error) {
			resp, err := genkit.Generate(ctx, g,
				ai.WithModel(model),
				ai.WithPrompt("讲一个关于%s的笑话。使用中文。", input.Topic),
				ai.WithUse(&middleware.Retry{
					MaxRetries:      3,     // 最多重试 3 次
					InitialDelayMs:  1000,  // 初始延迟 1 秒
					BackoffFactor:   2,     // 每次翻倍（1s -> 2s -> 4s）
				}),
			)
			if err != nil {
				return nil, fmt.Errorf("retry flow failed: %w", err)
			}
			return &JokeOutput{Joke: resp.Text()}, nil
		})

	// ====================================================================
	// 示例 2: 自定义 Logging 中间件 - 记录模型调用耗时
	// 使用自定义 Logger 中间件，每次模型调用会记录耗时日志
	// ====================================================================
	loggingFlow := genkit.DefineFlow(g, "loggingFlow",
		func(ctx context.Context, input *RequestInput) (*JokeOutput, error) {
			resp, err := genkit.Generate(ctx, g,
				ai.WithModel(model),
				ai.WithPrompt("讲一个关于%s的笑话。使用中文。", input.Topic),
				ai.WithUse(LoggingMiddleware{Prefix: "[性能追踪]"}),
			)
			if err != nil {
				return nil, fmt.Errorf("logging flow failed: %w", err)
			}
			return &JokeOutput{Joke: resp.Text()}, nil
		})

	// ====================================================================
	// 示例 3: 多层中间件组合 - Retry + Counter
	// ai.WithUse(A, B, C) 从左到右组合，最外层先执行
	// 组合顺序：Counter { Retry { actual } }
	// 这表示 Counter 统计的是包含重试在内的完整调用链路
	// ====================================================================
	composedFlow := genkit.DefineFlow(g, "composedFlow",
		func(ctx context.Context, input *RequestInput) (*JokeOutput, error) {
			resp, err := genkit.Generate(ctx, g,
				ai.WithModel(model),
				ai.WithPrompt("讲一个关于%s的笑话。使用中文。", input.Topic),
				ai.WithUse(
					CounterMiddleware{},           // 外层：统计调用
					&middleware.Retry{MaxRetries: 2}, // 内层：自动重试
				),
			)
			if err != nil {
				return nil, fmt.Errorf("composed flow failed: %w", err)
			}
			return &JokeOutput{Joke: resp.Text()}, nil
		})

	// ====================================================================
	// 示例 4: 内联中间件 - ai.MiddlewareFunc
	// 适用于不需要命名类型或 Dev UI 可见性的临时中间件
	// ====================================================================
	inlineMwFlow := genkit.DefineFlow(g, "inlineMwFlow",
		func(ctx context.Context, input *RequestInput) (*JokeOutput, error) {
			resp, err := genkit.Generate(ctx, g,
				ai.WithModel(model),
				ai.WithPrompt("讲一个关于%s的笑话。使用中文。", input.Topic),
				ai.WithUse(ai.MiddlewareFunc(func(ctx context.Context) (*ai.Hooks, error) {
					return &ai.Hooks{
						WrapModel: func(ctx context.Context, p *ai.ModelParams, next ai.ModelNext) (*ai.ModelResponse, error) {
							log.Printf("模型请求消息数: %d", len(p.Request.Messages))
							return next(ctx, p)
						},
					}, nil
				})),
			)
			if err != nil {
				return nil, fmt.Errorf("inline mw failed: %w", err)
			}
			return &JokeOutput{Joke: resp.Text()}, nil
		})

	// ====================================================================
	// 本地测试
	// ====================================================================

	testCases := []*RequestInput{
		{Topic: "程序员"},
		{Topic: "猫"},
	}

	for _, tc := range testCases {
		fmt.Printf("\n========== Flow 测试 (主题: %s) ==========\n", tc.Topic)

		// 测试 RetryFlow
		result, err := retryFlow.Run(ctx, tc)
		if err != nil {
			log.Printf("retryFlow error: %v", err)
		} else {
			b, _ := json.MarshalIndent(result, "", "  ")
			fmt.Printf("retryFlow 结果: %s\n", string(b))
		}

		// 测试 inlineMwFlow
		result2, err := inlineMwFlow.Run(ctx, tc)
		if err != nil {
			log.Printf("inlineMwFlow error: %v", err)
		} else {
			b, _ := json.MarshalIndent(result2, "", "  ")
			fmt.Printf("inlineMwFlow 结果: %s\n", string(b))
		}
	}

	// ====================================================================
	// HTTP 部署
	// ====================================================================
	mux := http.NewServeMux()
	mux.HandleFunc("POST /retryFlow", genkit.Handler(retryFlow))
	mux.HandleFunc("POST /loggingFlow", genkit.Handler(loggingFlow))
	mux.HandleFunc("POST /composedFlow", genkit.Handler(composedFlow))
	mux.HandleFunc("POST /inlineMwFlow", genkit.Handler(inlineMwFlow))

	fmt.Println("\n========== 已部署的 Flow 端点 ==========")
	for _, flow := range genkit.ListFlows(g) {
		fmt.Printf("  POST /%s\n", flow.Name())
	}

	log.Println("Starting server on http://localhost:3400")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}