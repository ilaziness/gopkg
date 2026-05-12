# 11 - 错误类型 (Error Types)

对应官方文档：[Error types](https://genkit.dev/docs/go/error-types/)

## 概述

Genkit 定义了一套标准错误类型（gRPC 风格状态码），用于统一处理 AI 工作流中的异常。

## 状态码列表

| 状态码 | 说明 |
|--------|------|
| OK | 成功 |
| CANCELLED | 操作被取消 |
| UNKNOWN | 未知错误 |
| INVALID_ARGUMENT | 无效参数 |
| DEADLINE_EXCEEDED | 超时 |
| NOT_FOUND | 未找到 |
| ALREADY_EXISTS | 已存在 |
| PERMISSION_DENIED | 权限拒绝 |
| RESOURCE_EXHAUSTED | 资源耗尽（配额限制） |
| FAILED_PRECONDITION | 前置条件失败 |
| ABORTED | 中止 |
| OUT_OF_RANGE | 超出范围 |
| UNIMPLEMENTED | 未实现 |
| INTERNAL | 内部错误 |
| UNAVAILABLE | 服务不可用（临时） |
| DATA_LOSS | 数据丢失 |
| UNAUTHENTICATED | 未认证 |

## 中间件集成

这些错误码在 Retry 和 Fallback 中间件中自动处理：

- **Retry**: 在 `RESOURCE_EXHAUSTED`、`UNAVAILABLE` 等状态时自动重试
- **Fallback**: 在这些状态时切换到备用模型

## 运行

```bash
cd framework/genkit/11-errors
go run .