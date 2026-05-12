package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/ollama"
	"github.com/firebase/genkit/go/plugins/server"
)

func main() {
	ctx := context.Background()

	ollamaPlugin := &ollama.Ollama{
		ServerAddress: "http://127.0.0.1:11434",
		Timeout:       80,
	}

	g := genkit.Init(ctx, genkit.WithPlugins(ollamaPlugin))
	ollamaPlugin.DefineModel(g, ollama.ModelDefinition{
		Name: "gemma4:e2b", Type: "generate",
	}, nil)

	fmt.Println("========== Evaluation (评估) ==========")
	fmt.Println("Genkit 提供评估框架来测试 AI 工作流的质量。")
	fmt.Println("")
	fmt.Println("核心概念：")
	fmt.Println("  1. 评估器 (Evaluator): 衡量输出的特定维度")
	fmt.Println("     - 内置评估器: 相关性、事实准确性、安全性等")
	fmt.Println("     - 自定义评估器: 根据需求定制")
	fmt.Println("  2. 数据集 (Dataset): 测试用例集合")
	fmt.Println("  3. 运行评估 (EvaluationRun): 执行并记录结果")
	fmt.Println("")
	fmt.Println("评估命令：")
	fmt.Println("  genkit eval:run flowName --input input.json")
	fmt.Println("  genkit eval:list  # 查看评估结果")
	fmt.Println("")
	fmt.Println("========== Local Observability (本地可观测性) ==========")
	fmt.Println("Genkit 内置本地可观测性支持，无需额外配置。")
	fmt.Println("")
	fmt.Println("功能特性：")
	fmt.Println("  1. Trace (追踪)")
	fmt.Println("     - 自动记录每次 Flow 执行的完整调用链")
	fmt.Println("     - 每个 genkit.Run() 步骤都是独立追踪节点")
	fmt.Println("     - 包含输入、输出、耗时信息")
	fmt.Println("  2. Developer UI 可视化")
	fmt.Println("     - 运行 genkit start -- go run .")
	fmt.Println("     - 访问 http://localhost:4000")
	fmt.Println("     - 在 Inspect 标签页查看追踪详情")
	fmt.Println("  3. Flow Steps (步骤追踪)")
	fmt.Println("     - genkit.Generate()/Embed()/Retrieve() 自动追踪")
	fmt.Println("     - genkit.Run() 包装自定义代码为追踪步骤")
	fmt.Println("")
	fmt.Println("示例：")
	fmt.Println(`  result, err := genkit.Run(ctx, "my-step", func() (T, error) {`)
	fmt.Println(`      // 此代码块会在 trace 中显示为独立步骤`)
	fmt.Println(`      return value, nil`)
	fmt.Println(`  })`)
	fmt.Println("")
	fmt.Println("更多信息：")
	fmt.Println("  https://genkit.dev/docs/go/evaluation/")
	fmt.Println("  https://genkit.dev/docs/go/local-observability/")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", genkit.Handler(nil))

	log.Println("Starting server on http://localhost:3400")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}