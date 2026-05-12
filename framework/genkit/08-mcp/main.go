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

	fmt.Println("========== Model Context Protocol (MCP) ==========")
	fmt.Println("MCP 是 Genkit 支持的标准协议，允许 AI 模型与外部工具和数据源交互。")
	fmt.Println("")
	fmt.Println("核心概念：")
	fmt.Println("  1. MCP Server: 提供工具和资源的服务器")
	fmt.Println("  2. MCP Client: 消费工具和资源的客户端")
	fmt.Println("  3. Genkit MCP Server: 将 Genkit 功能暴露为 MCP 服务")
	fmt.Println("")
	fmt.Println("Genkit MCP Server 用法：")
	fmt.Println("  genkit mcp start -- go run .")
	fmt.Println("")
	fmt.Println("这会启动一个 MCP 服务器，将 Genkit 的 Flow 暴露为 MCP 工具。")
	fmt.Println("其他支持 MCP 的客户端（如 Claude Desktop）可以调用这些工具。")
	fmt.Println("")
	fmt.Println("更多信息请参阅：")
	fmt.Println("  - MCP 文档: https://genkit.dev/docs/go/model-context-protocol/")
	fmt.Println("  - MCP Server: https://genkit.dev/docs/go/mcp-server/")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", genkit.Handler(nil))

	log.Println("Starting server on http://localhost:3400")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}