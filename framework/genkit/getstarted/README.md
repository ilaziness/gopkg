# Genkit Get Started

本目录演示如何使用 Genkit 在 Go 中创建第一个 AI 应用程序。

## 主要内容

1. 先决条件
   - 需要 Go 1.24 或更高版本。

2. 初始化项目
   - 在项目根目录中运行 `go mod init <module>`。
   - 创建 `main.go` 作为应用入口。

3. 安装 Genkit
   - 安装 Genkit CLI：`curl -sL cli.genkit.dev | bash`
   - 安装 Go 包：`go get github.com/firebase/genkit/go`

4. 配置模型 API Key
   - 本指南默认使用 Gemini API。
   - 在系统环境中设置 `GEMINI_API_KEY`。

5. 创建第一个应用
   - 定义输入和输出结构体作为 schema。
   - 使用 `genkit.Init` 初始化 Genkit，配置插件和默认模型。
   - 使用 `genkit.DefineFlow` 定义 Flow，封装模型调用逻辑。
   - 通过 `genkit.GenerateData[Recipe]` 生成结构化输出。
   - 调用 `recipeGeneratorFlow.Run` 做一次测试。

### 为什么使用 Flow？

- 类型安全的输入和输出：使用结构体 schema 明确数据类型。
- 与 Developer UI 集成：可视化测试和调试 Flow。
- 易于部署为 API：Flow 可以直接作为 HTTP 接口暴露。
- 内置追踪与可观测性：方便监控性能与调试模型行为。

6. 运行应用
   - 执行 `go run .` 来运行程序。
   - 程序会生成一个示例菜谱，并启动 HTTP 服务。

7. HTTP 测试
   - 在另一个终端使用 `curl` 向 `http://localhost:3400/recipeGeneratorFlow` 发送 POST 请求。
   - 请求体格式为 `{"data": {"ingredient": "tomato", "dietaryRestrictions": "vegan"}}`。

8. 使用 Developer UI
   - 启动 Dev UI：`genkit start -- go run .`
   - 在 `http://localhost:4000` 中可视化调试和运行 Flow。

## 目录结构

- `main.go`：示例应用入口。
- `README.md`：本文件。

## 后续学习

- [Genkit 开发工具](https://genkit.dev/docs/go/devtools/)
- [生成内容](https://genkit.dev/docs/go/models/)
- [创建 Flow](https://genkit.dev/docs/go/flows/)
- [工具调用](https://genkit.dev/docs/go/tool-calling/)
- [Dotprompt 提示管理](https://genkit.dev/docs/go/dotprompt/)
