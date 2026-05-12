# 04 - 工具调用 (Tool Calling)

对应官方文档：[Tool calling](https://genkit.dev/docs/go/tool-calling/)

## 概述

工具调用（又称函数调用）让 LLM 能够向应用发起结构化请求，从而获取实时数据、执行计算或触发操作。

## 核心概念

### 工具定义

使用 `genkit.DefineTool()` 定义工具：

```go
getWeatherTool := genkit.DefineTool(
    g,
    "getWeather",
    "获取指定地点的当前天气信息。",
    func(ctx *ai.ToolContext, input WeatherInput) (string, error) {
        return "晴, 22°C", nil
    },
)
```

### 自动工具调用

LLM 自动决定何时使用工具，Genkit 自动执行工具并继续对话：

```go
resp, err := genkit.Generate(ctx, g,
    ai.WithPrompt("北京的天气怎么样？"),
    ai.WithTools(getWeatherTool),
)
```

### 显式工具调用

需要完全控制时使用 `WithReturnToolRequests(true)`：

```go
resp, err := genkit.Generate(ctx, g,
    ai.WithTools(getWeatherTool),
    ai.WithReturnToolRequests(true),
)
// 手动处理每个工具请求
```

## 适用场景

| 场景 | 示例 |
|------|------|
| 实时数据 | 股票价格、天气查询 |
| 确定性操作 | 数学计算、模板文本生成 |
| 触发动作 | 开关灯、预订餐厅 |

## 本目录示例

| 示例 | 说明 |
|------|------|
| 天气查询 | 模型自动调用 getWeather 工具 |
| 菜单问答 | 餐厅助手使用菜单工具解答 |
| 显式控制 | 手动处理工具调用循环 |

## 运行

```bash
cd framework/genkit/04-tools
go run .