// Graph编排 — 使用 cloudwego/eino 的 Graph 示例
package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// EinoGraphExample — 最小可运行的 Graph 示例，仅使用一个 Lambda 节点
// 来模拟模板+模型的处理流程，便于在本地快速运行并查看编排调用。
func EinoGraphExample(ctx context.Context) error {
	// Graph 输入为 map[string]any，输出为 *schema.Message
	graph := compose.NewGraph[map[string]any, *schema.Message]()

	// 单个 lambda 节点模拟模板+模型：接受输入 map -> 返回 *schema.Message
	lambda := func(ctx context.Context, in map[string]any) (*schema.Message, error) {
		q, _ := in["query"].(string)
		resp := fmt.Sprintf("mock reply for query: %s", q)
		return schema.AssistantMessage(resp, nil), nil
	}

	// 加入 lambda 节点并连通 START -> node -> END
	_ = graph.AddLambdaNode("node_model", compose.InvokableLambda(lambda))
	_ = graph.AddEdge(compose.START, "node_model")
	_ = graph.AddEdge("node_model", compose.END)

	// 编译并执行
	compiled, err := graph.Compile(ctx)
	if err != nil {
		return err
	}
	out, err := compiled.Invoke(ctx, map[string]any{"query": "Beijing's weather this weekend"})
	if err != nil {
		return err
	}
	fmt.Printf("graph output: %+v\n", out)
	return nil
}
