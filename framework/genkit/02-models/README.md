# 02 - 使用 AI 模型生成内容 (Generating content)

对应官方文档：[Generating content with AI models](https://genkit.dev/docs/go/models/)

## 概述

Genkit 提供统一的接口 `genkit.Generate()` 与各种 AI 模型进行交互。配置好模型插件后，所有模型都通过同一 API 调用，方便组合多个模型或在应用中切换模型。

## 核心函数

### genkit.Generate() - 基础生成

最核心的函数，接收模型、提示词等参数，返回 `*ai.ModelResponse`。

```go
resp, err := genkit.Generate(ctx, g,
    ai.WithModel(model),
    ai.WithPrompt("你的提示词"),
)
text := resp.Text() // 获取生成的文本
```

### genkit.GenerateText() - 便捷文本生成

直接返回纯文本字符串，适合简单场景：

```go
text, err := genkit.GenerateText(ctx, g,
    ai.WithPrompt("你的提示词"),
)
```

### genkit.GenerateData[T]() - 结构化输出

指定 Go struct 作为泛型参数，模型输出自动解析为结构化数据：

```go
type Item struct {
    Name  string `json:"name"`
    Price string `json:"price"`
}

item, _, err := genkit.GenerateData[Item](ctx, g,
    ai.WithPrompt("生成一个菜品"),
)
// item 是 *Item 类型，直接访问字段
```

### genkit.GenerateStream() - 流式生成

返回 Go 迭代器，逐块输出生成内容，适合实时场景：

```go
for result, err := range genkit.GenerateStream(ctx, g, ...) {
    if result.Done {
        fullText = result.Response.Text()
    } else {
        fmt.Print(result.Chunk.Text())
    }
}
```

### genkit.GenerateDataStream[T]() - 结构化流式

流式生成并自动解析为结构化类型，中间值和最终输出都是强类型：

```go
for result, err := range genkit.GenerateDataStream[*Item](ctx, g, ...) {
    if result.Done {
        return result.Output, nil  // *Item 类型
    }
    sendChunk(ctx, result.Chunk)   // 也是 *Item 类型
}
```

## 关键概念

### 系统提示词 (System Prompt)

使用 `ai.WithSystem()` 设置模型的行为角色：

```go
ai.WithSystem("你是米其林三星主厨，用英文回复。")
```

### 模型参数控制

使用 `ai.WithConfig()` 配置生成参数：

- **Temperature**: 控制创造性（0=确定性，>1=创造性）
- **MaxOutputTokens**: 限制输出长度（以 token 计，英文单词≈2-4 tokens）
- **StopSequences**: 停止序列（遇到这些字符停止生成）
- **TopP**: 按累积概率筛选候选词（0.0-1.0）
- **TopK**: 只考虑概率最高的 K 个候选词

```go
ai.WithConfig(&ai.GenerationCommonConfig{
    Temperature:     0.8,
    MaxOutputTokens: 500,
    StopSequences:   []string{"</end>"},
})
```

### 模型标识符

模型用 `providerid/modelid` 格式指定，例如：
- `googleai/gemini-2.5-flash`
- `ollama/gemma4:e2b`

可通过 `genkit.WithDefaultModel()` 设置默认模型。

## 代码示例说明

本目录示例演示了 7 种生成模式：

| 示例 | 函数 | 场景 |
|------|------|------|
| 1 | `genkit.Generate()` | 基础文本生成 |
| 2 | `genkit.GenerateData[T]()` | 结构化数据输出 |
| 3 | `ai.WithSystem()` | 系统提示词设定角色 |
| 4 | `ai.WithConfig()` | 模型参数调优 |
| 5 | `genkit.GenerateStream()` | 流式文本输出 |
| 6 | `genkit.GenerateDataStream[T]()` | 流式结构化输出 |
| 7 | `genkit.GenerateText()` | 简短便捷文本 |

## 运行方式

```bash
# 本地 Ollama 服务需运行
cd framework/genkit/02-models
go run .
```

## HTTP 部署端点

所有 Flow 注册为 HTTP 端点：

```bash
curl -X POST "http://localhost:3400/structuredOutputFlow" \
  -H "Content-Type: application/json" \
  -d '{"data": {"theme": "海洋", "cuisine": "海鲜"}}'