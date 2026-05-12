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

	fmt.Println("========== Error Types (错误类型) ==========")
	fmt.Println("Genkit 定义了一套标准错误类型，用于统一处理 AI 工作流中的异常。")
	fmt.Println("")
	fmt.Println("标准错误状态码：")
	fmt.Println("")
	fmt.Println("  core.StatusName 类型：")
	fmt.Println("  ----------------------------")
	fmt.Println("  OK                 - 成功")
	fmt.Println("  CANCELLED          - 操作被取消")
	fmt.Println("  UNKNOWN            - 未知错误")
	fmt.Println("  INVALID_ARGUMENT   - 无效参数")
	fmt.Println("  DEADLINE_EXCEEDED  - 超时")
	fmt.Println("  NOT_FOUND          - 未找到")
	fmt.Println("  ALREADY_EXISTS     - 已存在")
	fmt.Println("  PERMISSION_DENIED  - 权限拒绝")
	fmt.Println("  RESOURCE_EXHAUSTED - 资源耗尽（如配额限制）")
	fmt.Println("  FAILED_PRECONDITION - 前置条件失败")
	fmt.Println("  ABORTED            - 中止")
	fmt.Println("  OUT_OF_RANGE       - 超出范围")
	fmt.Println("  UNIMPLEMENTED      - 未实现")
	fmt.Println("  INTERNAL           - 内部错误")
	fmt.Println("  UNAVAILABLE        - 不可用（临时）")
	fmt.Println("  DATA_LOSS          - 数据丢失")
	fmt.Println("  UNAUTHENTICATED    - 未认证")
	fmt.Println("")
	fmt.Println("错误处理示例：")
	fmt.Println("")
	fmt.Println("  import \"github.com/firebase/genkit/go/core\"")
	fmt.Println("")
	fmt.Println(`  if err != nil {`)
	fmt.Println(`      // 检查错误状态`)
	fmt.Println(`      status := core.ErrorStatus(err)`)
	fmt.Println(`      if status == core.RESOURCE_EXHAUSTED {`)
	fmt.Println(`          // 处理配额不足`)
	fmt.Println(`      }`)
	fmt.Println(`  }`)
	fmt.Println("")
	fmt.Println("这些错误码在中间件（如 Retry、Fallback）中自动处理。")
	fmt.Println("  - Retry 中间件在 RESOURCE_EXHAUSTED 等状态时自动重试")
	fmt.Println("  - Fallback 中间件在这些状态时切换到备用模型")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", genkit.Handler(nil))

	log.Println("Starting server on http://localhost:3400")
	log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
}