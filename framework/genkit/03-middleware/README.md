# 03 - 中间件 (Middleware)

对应官方文档：[Middleware](https://genkit.dev/docs/go/middleware/)

## 概述

Genkit 中间件提供了一种在 `generate()` 调用中修改模型行为的方式。中间件可用于重试失败请求、回退到不同模型、注入工具和上下文等。

## 核心概念

### Middleware 接口

中间件需满足 `ai.Middleware` 接口：

```go
type Middleware interface {
    Name() string
    New(ctx context.Context) (*ai.Hooks, error)
}
```

### Hooks（钩子）

`New()` 返回 `*ai.Hooks`，包含三种钩子类型：

| 钩子 | 触发时机 | 用途 |
|------|----------|------|
| `WrapGenerate` | 每次工具循环迭代 | 查看整个对话（重写、系统提示注入） |
| `WrapModel` | 每次模型 API 调用 | 模型调用逻辑（重试、回退、缓存） |
| `WrapTool` | 每次工具执行 | 单次工具执行（审批、沙箱、日志） |

## 内置中间件

### Retry - 自动重试

失败时自动重试，使用指数退避 + jitter：

```go
ai.WithUse(&middleware.Retry{
    MaxRetries:     3,
    InitialDelayMs: 1000,
    BackoffFactor:  2,
})
```

### Fallback - 模型回退

主模型失败时自动切换到备用模型：

```go
ai.WithUse(&middleware.Fallback{
    Models: []ai.ModelRef{...},
    Statuses: []core.StatusName{core.RESOURCE_EXHAUSTED},
})
```

### ToolApproval - 工具审批

限制工具执行到允许列表，未批准的触发中断：

```go
ai.WithUse(&middleware.ToolApproval{
    AllowedTools: []string{},
})
```

### Skills - 技能注入

扫描目录中的 SKILL.md 文件注入到系统提示词：

```go
ai.WithUse(&middleware.Skills{SkillPaths: []string{"./skills"}})
```

### Filesystem - 文件系统访问

赋予模型文件系统操作能力（根目录限制）：

```go
ai.WithUse(&middleware.Filesystem{
    RootDir: "./workspace",
    AllowWriteAccess: true,
})
```

## 构建自定义中间件

```go
type LoggingMiddleware struct {
    Prefix string `json:"prefix,omitempty"`
}

func (LoggingMiddleware) Name() string { return "mine/logger" }

func (l LoggingMiddleware) New(ctx context.Context) (*ai.Hooks, error) {
    return &ai.Hooks{
        WrapModel: func(ctx context.Context, p *ai.ModelParams, next ai.ModelNext) (*ai.ModelResponse, error) {
            start := time.Now()
            resp, err := next(ctx, p)
            log.Printf("%s 模型调用耗时: %s", l.Prefix, time.Since(start))
            return resp, err
        },
    }, nil
}
```

## 组合顺序

`ai.WithUse(A, B, C)` 从左到右组合，最外层先执行：

```go
// 链式：Retry { Fallback { model } }
ai.WithUse(
    &middleware.Retry{MaxRetries: 3},       // 外层
    &middleware.Fallback{Models: models},   // 内层
)
```

## 本目录示例

| 示例 | 说明 |
|------|------|
| Retry 中间件 | 自动重试失败的模型调用 |
| Logging 中间件 | 自定义中间件记录模型耗时 |
| 多层组合 | Counter + Retry 组合使用 |
| 内联中间件 | 使用 `ai.MiddlewareFunc` 快速创建 |

## 运行

```bash
cd framework/genkit/03-middleware
go run .