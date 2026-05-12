# 08 - Model Context Protocol (MCP)

对应官方文档：[Model Context Protocol](https://genkit.dev/docs/go/model-context-protocol/), [Genkit MCP server](https://genkit.dev/docs/go/mcp-server/)

## 概述

MCP 是 Genkit 支持的开放协议，允许 AI 模型与外部工具和数据源交互。类似于 USB-C 为设备提供标准连接，MCP 为 AI 应用提供标准化的工具集成方式。

## 核心组件

### MCP Server
提供工具和资源的服务器端。Genkit 可以将其 Flow 暴露为 MCP 工具。

### MCP Client
消费工具和资源的客户端。Genkit 应用可以作为 MCP 客户端调用外部 MCP 服务。

### Genkit MCP Server
通过 `genkit mcp start -- go run .` 启动，将所有 Flow 暴露为 MCP 工具。

## 使用方式

```bash
# 启动 Genkit MCP Server
genkit mcp start -- go run .
```

然后其他 MCP 客户端（如 Claude Desktop）可以调用这些工具。

## 运行

```bash
cd framework/genkit/08-mcp
go run .