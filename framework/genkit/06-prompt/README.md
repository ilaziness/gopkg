# 06 - Dotprompt 提示管理 (Managing prompts with Dotprompt)

对应官方文档：[Managing prompts with Dotprompt](https://genkit.dev/docs/go/dotprompt/)

## 概述

Dotprompt 是 Genkit 的提示词管理框架，将提示词定义在 `.prompt` 文件中，与代码分离。

## 核心概念

### .prompt 文件

```yaml
---
model: googleai/gemini-2.5-flash
config:
  temperature: 0.9
input:
  schema:
    location: string
    style?: string
output:
  schema:
    message: string
---
你是一个热情的 AI 助手，当前在 {{location}} 工作。
用 {{style}} 风格问候客人。
```

### 从代码加载

```go
// 普通加载
prompt := genkit.LookupPrompt(g, "hello")
resp, err := prompt.Execute(ctx, ai.WithInput(map[string]any{"location": "北京"}))

// 强类型加载
prompt := genkit.LookupDataPrompt[Input, *Output](g, "hello")
result, resp, err := prompt.Execute(ctx, input)
```

### 特性

- **提示模板**：Handlebars 模板语法 `{{variable}}`
- **多消息提示**：支持 system/user/assistant 多条消息
- **Schema 定义**：Picoschema / JSON Schema / 代码引用
- **工具调用**：prompt 文件中可直接引用工具
- **Partials**：可复用的提示片段
- **变体**：同一 prompt 多个版本

## 运行

```bash
cd framework/genkit/06-prompt
go run .