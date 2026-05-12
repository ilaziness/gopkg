# 10 - Evaluation + Local Observability

对应官方文档：[Evaluation](https://genkit.dev/docs/go/evaluation/), [Local observability and metrics](https://genkit.dev/docs/go/local-observability/)

## Evaluation (评估)

Genkit 提供评估框架来测试 AI 工作流的质量。

- **Evaluator**: 衡量输出的维度（相关性、事实准确性、安全性等）
- **Dataset**: 测试用例集合
- **EvaluationRun**: 执行并记录评估结果

```bash
genkit eval:run flowName --input input.json
```

## Local Observability (本地可观测性)

Genkit 内置本地可观测性：

- **Trace**: 自动记录 Flow 执行的完整调用链
- **Flow Steps**: 使用 `genkit.Run()` 创建独立追踪步骤
- **Developer UI**: `genkit start -- go run .` → http://localhost:4000

```go
result, err := genkit.Run(ctx, "my-step", func() (T, error) {
    // 此代码块会在 trace 中显示为独立步骤
    return value, nil
})
```

## 运行

```bash
cd framework/genkit/10-observability
go run .