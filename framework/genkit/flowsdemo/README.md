# Genkit Flows 详解

本目录总结了 [Genkit Flows 官方文档](https://genkit.dev/docs/go/flows/) 的核心概念，并提供了一个完整的示例项目。

## 什么是 Flow？

Flow 是 Genkit 中封装 AI 逻辑的特殊函数，它提供了以下能力：

| 特性 | 说明 |
|------|------|
| **类型安全的输入/输出** | 使用 Go struct 定义 schema，编译期和运行时双重验证 |
| **流式支持** | 支持流式生成部分响应数据 |
| **Developer UI 集成** | 通过可视化界面测试和调试 Flow |
| **易于部署** | 可直接作为 HTTP 端点部署到任意平台 |

Flow 是轻量级的，它们像普通函数一样编写，几乎无抽象。

## 文档核心要点

### 1. 核心函数概念

#### genkit.Init()
- **作用**：初始化 Genkit 运行时环境
- **参数**：ctx（context.Context），可选的插件列表、模型选项等
- **返回值**：`*genkit.Genkit` 实例，所有 Flow 都依赖此实例注册

```go
g := genkit.Init(ctx, genkit.WithPlugins(&googlegenai.GoogleAI{}))
```

#### genkit.DefineFlow()
- **作用**：定义一个新的 Flow（非流式）
- **参数**：g（`*genkit.Genkit`），flowName（string），回调函数 `func(ctx, input) (output, error)`
- **返回值**：`*genkit.Flow[I, O]`，可通过 `.Run()` 调用
- **输入/输出类型安全**：通过 Go 泛型 `[I, O]` 自动推断

```go
flow := genkit.DefineFlow(g, "myFlow",
    func(ctx context.Context, input MyInput) (MyOutput, error) {
        // ... AI 逻辑
    })
```

#### genkit.DefineStreamingFlow()
- **作用**：定义支持流式输出的 Flow
- **参数**：同 `DefineFlow`，但回调函数额外接收 `sendChunk` 参数
- **sendChunk 类型**：`core.StreamCallback[T]`，一个回调函数，用于在生成过程中推送中间结果

```go
flow := genkit.DefineStreamingFlow(g, "myStreamFlow",
    func(ctx context.Context, input MyInput, sendChunk core.StreamCallback[string]) (MyOutput, error) {
        // sendChunk(ctx, chunkData) 用于推送流式数据
    })
```

#### flow.Run()
- **作用**：执行 Flow，阻塞等待最终结果
- **返回**：(output, error)

```go
result, err := flow.Run(ctx, inputData)
```

#### flow.Stream()
- **作用**：以流式方式执行 Flow，返回一个 channel
- **返回**：`<-chan *core.StreamingFlowValue[I, O]`
- **channel 中的值**：
  - `result.Done == true`：`result.Output` 包含最终结果
  - `result.Done == false`：`result.Stream` 包含中间流式数据块

```go
ch := flow.Stream(ctx, inputData)
for r := range ch {
    if r.Done { fmt.Println(r.Output) } else { fmt.Print(r.Stream) }
}
```

#### genkit.Generate()
- **作用**：执行模型生成调用，返回原始文本响应
- **返回值**：`*ai.ModelResponse`，通过 `.Text()` 获取生成的文本

```go
resp, err := genkit.Generate(ctx, g, ai.WithPrompt("..."))
text := resp.Text()
```

#### genkit.GenerateData[T]()
- **作用**：执行模型生成并自动解析为结构化类型 T
- **泛型参数 T**：Go struct 类型，模型输出会被自动解析为该类型
- **返回值**：`(*T, *ai.ModelResponse, error)` - 解析后的结构体指针、完整响应、错误

```go
item, _, err := genkit.GenerateData[MenuItem](ctx, g, ai.WithPrompt("..."))
// item 是 *MenuItem，可直接访问 item.Name, item.Description
```

#### genkit.GenerateStream()
- **作用**：流式模型生成，返回迭代器
- **迭代器值**：
  - `result.Done == true`：`result.Response` 包含完整响应
  - `result.Done == false`：`result.Chunk` 包含当前文本块

```go
for result, err := range genkit.GenerateStream(ctx, g, ai.WithPrompt("...")) {
    if result.Done { fmt.Println(result.Response.Text()) }
}
```

#### genkit.GenerateDataStream[T]()
- **作用**：流式生成并自动解析为结构化类型 T
- **迭代器值**：
  - `result.Done == true`：`result.Output` 包含最终强类型输出
  - `result.Done == false`：`result.Chunk` 也是强类型的

```go
for result, err := range genkit.GenerateDataStream[*MenuItem](ctx, g, ...) {
    if result.Done { return result.Output, nil }  // *MenuItem
    sendChunk(ctx, result.Chunk)                   // 也是 *MenuItem
}
```

#### genkit.Run()
- **作用**：将一个代码块包装为独立的追踪步骤（在 Developer UI 中可见）
- **参数**：stepName（string），回调函数
- **用途**：将非 Genkit 原生操作（数据库查询、外部 API 调用等）纳入追踪

```go
result, err := genkit.Run(ctx, "step-name", func() (T, error) {
    // 任何代码，返回值类型 T
    return value, nil
})
```

#### genkit.Handler()
- **作用**：将 Flow 转换为 `http.HandlerFunc`，用于部署为 HTTP 端点
- **请求格式**：`POST {"data": {...}}`

```go
mux.HandleFunc("POST /myFlow", genkit.Handler(myFlow))
```

#### genkit.ListFlows()
- **作用**：列出 Genkit 实例中所有已注册的 Flow
- **用途**：批量注册 HTTP 路由

```go
for _, flow := range genkit.ListFlows(g) {
    mux.HandleFunc("POST /" + flow.Name(), genkit.Handler(flow))
}
```

#### ai.WithPrompt()
- **作用**：设置用户提示词，支持 `fmt.Sprintf` 风格的格式化参数

```go
ai.WithPrompt("Invent a menu item for a %s themed restaurant.", theme)
```

#### ai.WithSystem()
- **作用**：设置系统提示词，定义模型的行为角色

```go
ai.WithSystem("使用简体中文回答。")
```

#### ai.WithModel()
- **作用**：指定使用的模型

```go
ai.WithModel(model)
```

#### ai.WithStreaming()
- **作用**：在 Generate 调用中注册流式回调，接收逐块生成的响应

```go
ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
    return callback(ctx, chunk.Text())
})
```

#### ai.WithDocs()
- **作用**：为模型提供参考文档（RAG 场景）

```go
ai.WithDocs(ai.NewTextPart(menuContent))
```

#### core.StreamCallback[T]
- **类型定义**：`type StreamCallback[T] func(context.Context, T) error`
- **作用**：流式回调函数类型，推送到 Flow 的输出流

### 2. 定义和调用 Flow

Flow 最简单的形式就是包装一个函数：

```go
// 方式 A: 使用 genkit.Generate() 获取原始文本响应
menuSuggestionFlow := genkit.DefineFlow(g, "menuSuggestionFlow",
    func(ctx context.Context, input MenuSuggestionInput) (MenuSuggestionOutput, error) {
        resp, err := genkit.Generate(ctx, g,
            ai.WithPrompt("Invent a menu item for a %s themed restaurant.", input.Theme),
        )
        if err != nil {
            return MenuSuggestionOutput{}, err
        }
        return MenuSuggestionOutput{MenuItem: resp.Text()}, nil
    })

// 方式 B: 使用 genkit.GenerateData[T]() 获取结构化输出
menuSuggestionFlow := genkit.DefineFlow(g, "menuSuggestionFlow",
    func(ctx context.Context, input MenuSuggestionInput) (MenuItem, error) {
        item, _, err := genkit.GenerateData[MenuItem](ctx, g,
            ai.WithPrompt("Invent a menu item for a %s themed restaurant.", input.Theme),
        )
        return item, err
    })
```

调用方式：
```go
output, err := menuSuggestionFlow.Run(ctx, input)
```

### 3. 结构化 Schema 最佳实践

- 使用 Go struct + JSON tags 定义 schema
- **优先使用 struct 而非原始类型**，因为：
  - Developer UI 中显示更友好的输入字段（带标签）
  - 易于未来扩展（添加新字段不破坏现有客户端）
- Flow 的 schema **不必**与内部模型调用的 schema 一致。Flow 可以对模型输出做后处理再返回不同的类型：

```go
type MenuSuggestionInput struct {
    Theme string `json:"theme"`
}

type MenuItem struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}

type FormattedMenuOutput struct {
    FormattedMenuItem string `json:"formattedMenuItem"`
}

// Flow 的输出类型 FormattedMenuOutput 不同于内部 GenerateData 的 MenuItem
// 这证明了 Flow 可以有自己的输入/输出 schema，独立于其中调用的模型
menuSuggestionMarkdownFlow := genkit.DefineFlow(g, "menuSuggestionMarkdownFlow",
    func(ctx context.Context, input MenuSuggestionInput) (FormattedMenuOutput, error) {
        item, _, err := genkit.GenerateData[MenuItem](ctx, g,
            ai.WithPrompt("Invent a menu item for a %s themed restaurant.", input.Theme),
        )
        if err != nil {
            return FormattedMenuOutput{}, err
        }
        return FormattedMenuOutput{
            FormattedMenuItem: fmt.Sprintf("**%s**: %s", item.Name, item.Description),
        }, nil
    })
```

### 4. Streaming Flow

使用 `genkit.DefineStreamingFlow()` 定义流式 Flow。流式 Flow 的回调函数增加了一个 `sendChunk` 参数，用于在生成完成前推送中间结果。

> **Durable streaming（可选）**：对于长时间运行的流或网络可靠性要求高的场景，可以使用 [durable streaming](https://genkit.dev/docs/go/durable-streaming/)，允许客户端断线重连并回放内容。Go 版目前为实验性功能。

#### 迭代器模式（推荐）

`genkit.GenerateStream()` 返回一个 Go 迭代器（iterator），天然与 `for range` 语法集成。迭代器的每个值包含当前块内容或完成标志：

```go
// 文本流
menuSuggestionFlow := genkit.DefineStreamingFlow(g, "menuSuggestionFlow",
    func(ctx context.Context, theme string, sendChunk core.StreamCallback[string]) (string, error) {
        stream := genkit.GenerateStream(ctx, g, ai.WithPrompt(...))
        for result, err := range stream {
            if err != nil { return "", err }
            if result.Done {
                // result.Response 包含完整的生成响应
                return result.Response.Text(), nil
            }
            // result.Chunk 包含当前文本块，通过 sendChunk 推送给调用方
            sendChunk(ctx, result.Chunk.Text())
        }
        return "", nil
    })
```

#### 结构化流输出

`genkit.GenerateDataStream[T]()` 是结构化版本的流式调用，迭代值中的 `Chunk` 和 `Output` 都是强类型 `*T`：

```go
menuSuggestionFlow := genkit.DefineStreamingFlow(g, "menuSuggestionFlow",
    func(ctx context.Context, theme string, sendChunk core.StreamCallback[*MenuItem]) (*MenuItem, error) {
        stream := genkit.GenerateDataStream[*MenuItem](ctx, g, ai.WithPrompt(...))
        for result, err := range stream {
            if result.Done {
                // result.Output 是 *MenuItem
                return result.Output, nil
            }
            // result.Chunk 也是 *MenuItem
            sendChunk(ctx, result.Chunk)
        }
        return nil, nil
    })
```

#### 回调模式

如果你需要将流式回调与其他 `Generate` 选项结合（例如同时设置 `WithPrompt`、`WithSystem`、`WithModel`），可以使用 `ai.WithStreaming()`：

```go
menuSuggestionFlow := genkit.DefineStreamingFlow(g, "menuSuggestionFlow",
    func(ctx context.Context, input MenuSuggestionInput, callback core.StreamCallback[string]) (Menu, error) {
        item, _, err := genkit.GenerateData[MenuItem](ctx, g,
            ai.WithPrompt("Invent a menu item for a %s themed restaurant.", input.Theme),
            ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
                // chunk.Text() 是当前生成的文本片段
                // 通过 callback 推送到 Flow 的输出流
                return callback(ctx, chunk.Text())
            }),
        )
        if err != nil { return Menu{}, err }
        return Menu{Theme: input.Theme, Items: []MenuItem{item}}, nil
    })
```

注意：`StreamCallback[string]` 中的 `string` 指定流式值的类型，它不必与 Flow 的最终返回类型 `Menu` 相同。

#### Channel 模式（实验性）

对于更符合 Go 风格的编程，可以使用 experimental 的 channel API。函数接收一个 `chan<- string`（只写通道）替代回调：

```go
import x "github.com/firebase/genkit/go/genkit/x"

jokeFlow := x.DefineStreamingFlow(g, "jokeFlow",
    func(ctx context.Context, topic string, streamCh chan<- string) (string, error) {
        for chunk, err := range genkit.GenerateStream(ctx, g,
            ai.WithPrompt("Tell me a joke about %s.", topic),
        ) {
            if err != nil { return "", err }
            if chunk.Done { return chunk.Response.Text(), nil }
            select {
            case streamCh <- chunk.Chunk.Text():
            case <-ctx.Done():
                return "", ctx.Err()
            }
        }
        return "", nil
    })
```

Channel 模式要点：
- 函数接收 `chan<- string`（只写通道），替代回调
- 通道由框架管理，函数返回后自动关闭
- 函数**不应当关闭**该通道
- 使用 `select` + `ctx.Done()` 优雅处理取消

### 5. 调用 Streaming Flow

Streaming Flow 可以像普通 Flow 一样通过 `Run()` 调用（只获取最终结果），也可以通过 `Stream()` 获取流式数据：

```go
// 非流式调用（只获取最终结果）
output, err := flow.Run(ctx, input)

// 流式调用
streamCh := flow.Stream(ctx, input) // 注意：只返回一个值（channel），没有 error
for result := range streamCh {
    if result.Done {
        // result.Output 包含最终输出
        fmt.Println(result.Output.Theme)
    } else {
        // result.Stream 包含流式数据块
        fmt.Print(result.Stream)
    }
}
```

### 6. 命令行运行

```bash
# 运行普通 Flow
genkit flow:run menuSuggestionFlow '{"theme": "French"}'

# 运行 Streaming Flow（显示流式输出）
genkit flow:run menuSuggestionFlow '{"theme": "French"}' -s
```

### 7. 调试 Flow

启动 Developer UI：

```bash
genkit start -- go run .
```

访问 `http://localhost:4000` 可视化调试。

> **注意**：Developer UI 依赖 Go 应用保持运行。如果 Genkit 不是更大应用的一部分，可以在 `main()` 最后添加 `select {}` 防止退出，以便在 UI 中检查。`select {}` 是一个永远阻塞的空选择，让 main goroutine 不会退出。

从 **Run** 标签页可以运行项目中定义的任何 Flow。运行后可通过 **View trace** 或 **Inspect** 标签页查看调用追踪。

### 8. Flow Steps（追踪）

Genkit 的基础操作（`genkit.Generate()`、`genkit.Embed()`、`genkit.Retrieve()`）会自动显示为追踪中的独立步骤，每个步骤包含输入、输出和耗时信息。

使用 `genkit.Run()` 将其他代码包装为独立的追踪步骤（例如第三方库调用、数据库查询等关键代码段）。这可以让 Developer UI 的 trace 视图中看到完整的执行流程：

```go
menu, err := genkit.Run(ctx, "retrieve-daily-menu", func() (string, error) {
    // 任何非 Genkit 原生操作（数据库查询、API 调用等）
    return menu, nil
})
```

trace 视图中会显示一个名为 "retrieve-daily-menu" 的步骤节点。

### 9. 部署 Flow

**net/http 方式**：
```go
mux := http.NewServeMux()
mux.HandleFunc("POST /menuSuggestionFlow", genkit.Handler(menuSuggestionFlow))
log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
```

`server.Start()` 是 Genkit 的辅助函数，自动处理生命周期管理和中断信号。也可以使用标准的 `http.ListenAndServe`。

**批量注册所有 Flows**（推荐，避免遗漏）：
```go
for _, flow := range genkit.ListFlows(g) {
    mux.HandleFunc("POST /" + flow.Name(), genkit.Handler(flow))
}
log.Fatal(server.Start(ctx, "127.0.0.1:3400", mux))
```

**Gin 框架**：
```go
router := gin.Default()
for _, flow := range genkit.ListFlows(g) {
    router.POST("/"+flow.Name(), func(c *gin.Context) {
        genkit.Handler(flow)(c.Writer, c.Request)
    })
}
log.Fatal(router.Run(":3400"))
```

### 10. 调用已部署的 Flow

```bash
# 普通请求（请求体格式：{"data": {input fields}}）
curl -X POST "http://localhost:3400/menuSuggestionFlow" \
  -H "Content-Type: application/json" \
  -d '{"data": {"theme": "banana"}}'

# 流式响应（添加 Accept: text/event-stream 头）
curl -X POST "http://localhost:3400/menuSuggestionFlow" \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{"data": {"theme": "banana"}}'
```

也可以使用 Genkit 的 Web 客户端库从 Web 应用调用，详见 [Accessing flows from the client](https://genkit.dev/docs/client/)。

更多部署指南：[Deploy with Cloud Run](https://genkit.dev/docs/go/deployment/cloud-run/)

## 本目录示例说明

本目录中的示例项目演示了以下 Flow 特性：

1. **基本 Flow** - 定义和运行简单的生成 Flow（`GenerateData[T]`）
2. **结构化输出 Flow** - Flow 输出 schema 不同于内部模型调用的 schema
3. **Streaming Flow** - 使用迭代器模式实现流式响应
4. **Flow Steps** - 使用 `genkit.Run()` 创建可追踪的执行步骤
5. **多 Flow 部署** - 使用 `genkit.ListFlows()` 批量注册路由

## 运行方式

```bash
# 当前示例使用 Ollama 插件，需运行本地 Ollama 服务
# 或修改为其他模型插件

cd framework/genkit/flowsdemo
go run .