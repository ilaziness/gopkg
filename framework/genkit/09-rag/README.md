# 09 - 检索增强生成 (RAG)

对应官方文档：[Retrieval-augmented generation (RAG)](https://genkit.dev/docs/go/rag/)

## 概述

RAG 让 LLM 能够利用外部知识库生成更准确、更及时的回复，解决 LLM 知识过时和不存在的问题。

## 核心流程

```
用户问题 → 向量检索 → 检索相关文档 → 注入上下文 → LLM 生成
```

## 关键 API

### 注入文档上下文

```go
resp, err := genkit.Generate(ctx, g,
    ai.WithPrompt("根据以下文档回答问题"),
    ai.WithDocs(ai.NewTextPart("参考文档内容...")),
)
```

### 从向量库检索

```go
docs, err := genkit.Retrieve(ctx, g,
    ai.WithIndex("myIndex"),
    ai.WithRetriever("myRetriever"),
    ai.WithQuery("用户问题"),
)
```

## 支持的向量存储

Pinecone, Chroma, pgvector, LanceDB, Astra DB, Neo4j, AlloyDB, Cloud SQL PostgreSQL, Cloud Firestore, Dev local vector store。

## 运行

```bash
cd framework/genkit/09-rag
go run .