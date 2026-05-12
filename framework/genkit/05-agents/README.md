# 05 - Agent 模式 (Implementing Agentic Patterns)

对应官方文档：[Implementing Agentic Patterns](https://genkit.dev/docs/go/agentic-patterns/)

## 概述

Agent 模式让 LLM 能够主动使用工具、维护对话状态、进行多步推理来完成任务。Agent = LLM + 工具 + 多轮对话。

## 核心概念

### ReAct 模式

思考过程：**分析需求 → 使用工具(如果需要) → 组织回答**

```go
ai.WithSystem("你是一个智能助手。思考过程：分析需求 -> 使用工具 -> 组织回答。")
```

### 多轮对话

通过 `ai.WithMessages()` 传入历史消息维护对话上下文：

```go
messages := []*ai.Message{
    ai.NewUserMessage("你好！"),
    ai.NewModelTextMessage("你好，有什么可以帮你？"),
    ai.NewUserMessage("查一下天气"),
}
```

### 工具集成

Agent 自动判断何时使用工具获取实时数据。

## 本目录示例

| 示例 | 说明 |
|------|------|
| 基础 Agent | LLM + 工具的基本 ReAct 模式 |
| 多轮对话 Agent | 带历史记录的连续对话 |

## 运行

```bash
cd framework/genkit/05-agents
go run .