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

	fmt.Println("========== 检索增强生成 (RAG) ==========")
	fmt.Println("RAG 让 LLM 能够利用外部知识库生成更准确、更及时的回复。")
	fmt.Println("")
	fmt.Println("核心流程：")
	fmt.Println("  1. 文档分块 (Document Chunking)")
	fmt.Println("  2. 向量嵌入 (Embedding)")
	fmt.Println("  3. 向量存储 (Vector Store)")
	fmt.Println("  4. 检索 (Retrieval)")
	fmt.Println("  5. 生成 (Generation)")
	fmt.Println("")
	fmt.Println("关键 API：")
	fmt.Println("")
	fmt.Println("ai.WithDocs() - 向模型提供参考文档：")
	fmt.Println(`  resp, err := genkit.Generate(ctx, g,`)
	fmt.Println(`      ai.WithPrompt("根据以下文档回答问题"),`)
	fmt.Println(`      ai.WithDocs(ai.NewTextPart("文档内容...")),`)
	fmt.Println(`  )`)
	fmt.Println("")
	fmt.Println("ai.Retrieve() / genkit.Retrieve() - 从向量库检索：")
	fmt.Println(`  docs, err := genkit.Retrieve(ctx, g,`)
	fmt.Println(`      ai.WithIndex("myIndex"),`)
	fmt.Println(`      ai.WithRetriever("myRetriever"),`)
	fmt.Println(`      ai.WithQuery("用户问题"),`)
	fmt.Println(`  )`)
	fmt.Println("")
	fmt.Println("支持的向量存储提供商：")
	fmt.Println("  - Pinecone, Chroma, pgvector, LanceDB")
	fmt.Println("  - Astra DB, Neo4j, AlloyDB, Cloud SQL")
	fmt.Println("  - Dev local vector store (开发用)")
	fmt.Println("")
	fmt.Println("更多信息请参阅：")
	fmt.Println("  https://genkit.dev/docs/go/rag/")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", genkit.Handler(nil))

	log.Println("Starting server on http://localhost:3400")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}