# 07 - 聊天会话 + Context + Interrupts

对应官方文档：[Creating persistent chat sessions](https://genkit.dev/docs/go/chat/), [Passing information through context](https://genkit.dev/docs/go/context/), [Pause generation using interrupts](https://genkit.dev/docs/go/interrupts/)

## 概念概览

### 聊天会话 (Chat Sessions)

使用 `ai.WithMessages()` 传递消息历史，维护多轮对话上下文。

```go
history = append(history, ai.NewUserMessage("你好！"))
history = append(history, ai.NewModelTextMessage("你好！有什么可以帮你？"))

resp, err := genkit.Generate(ctx, g,
    ai.WithMessages(history...),
    ai.WithPrompt("今天的天气？"),
)
```

### Context (上下文注入)

使用 `ai.WithDocs()` 向模型提供参考文档（RAG 场景）：

```go
ai.WithDocs(ai.NewTextPart("参考文档内容"))
```

### Interrupts (生成中断)

`ai.WithReturnToolRequests(true)` 可暂停生成，等用户确认后继续。

## 本目录示例

- **持久化聊天会话**: 模拟维护多轮对话历史
- 概念说明: Context 注入和 Interrupts 机制

## 运行

```bash
cd framework/genkit/07-chat
go run .