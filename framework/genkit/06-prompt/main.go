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

	g := genkit.Init(ctx,
		genkit.WithPlugins(ollamaPlugin),
		genkit.WithPromptDir("./prompts"),
	)
	model := ollamaPlugin.DefineModel(
		g,
		ollama.ModelDefinition{
			Name: "gemma4:e2b",
			Type: "generate",
		},
		nil,
	)

	fmt.Println("========== Dotprompt 功能演示 ==========")
	fmt.Println("Dotprompt 允许将提示词定义在 .prompt 文件中，")
	fmt.Println("与代码分离，便于迭代和管理。")
	fmt.Println("")
	fmt.Println("创建 prompts/ 目录并在其中放置 .prompt 文件即可自动加载。")
	fmt.Println("")
	fmt.Println("示例 .prompt 文件格式：")
	fmt.Println("---")
	fmt.Println("model: ollama/gemma4:e2b")
	fmt.Println("config:")
	fmt.Println("  temperature: 0.9")
	fmt.Println("input:")
	fmt.Println("  schema:")
	fmt.Println("    location: string")
	fmt.Println("    style?: string")
	fmt.Println("---")
	fmt.Println("你是一个热情的 AI 助手，当前在 {{location}} 工作。")
	fmt.Println("用 {{style}} 风格问候客人。")
	fmt.Println("")
	fmt.Println("在代码中加载和使用：")
	fmt.Println(`  prompt := genkit.LookupPrompt(g, "hello")`)
	fmt.Println(`  resp, err := prompt.Execute(ctx, ai.WithInput(...))`)
	fmt.Println("")
	fmt.Println("使用 LookupDataPrompt 获得强类型支持：")
	fmt.Println(`  prompt := genkit.LookupDataPrompt[In, Out](g, "hello")`)
	fmt.Println(`  result, resp, err := prompt.Execute(ctx, input)`)
	fmt.Println("")
	fmt.Println("使用 Developer UI 可可视化编辑并导出 .prompt 文件。")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", genkit.Handler(nil))

	log.Println("Starting server on http://localhost:3400")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}