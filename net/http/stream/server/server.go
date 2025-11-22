// server.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type StreamMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Index   int    `json:"index,omitempty"`
}

type MCPMessage struct {
	Type      string      `json:"type"` // "message", "tool_call", "tool_result", "done"
	Role      string      `json:"role,omitempty"`
	Content   string      `json:"content,omitempty"`
	ToolName  string      `json:"tool_name,omitempty"`
	Arguments interface{} `json:"arguments,omitempty"`
	Result    string      `json:"result,omitempty"`
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应头：NDJSON + 流式传输
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 支持跨域

	// 确保 header 被立即发送（某些中间件可能缓冲）
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// 模拟逐步生成数据（如 LLM token）
	// 监听客户端断开连接
	ctx := r.Context()
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			log.Println("Client disconnected")
			return
		default:
		}
		msg := StreamMessage{
			Type:    "message",
			Content: fmt.Sprintf("Token %d", i),
			Index:   i,
		}

		// 序列化为 JSON 并写入一行
		if err := json.NewEncoder(w).Encode(msg); err != nil {
			log.Printf("Error encoding message: %v", err)
			return
		}

		// 强制 flush，确保客户端立即收到
		if f, ok := w.(http.Flusher); ok {
			log.Printf("flush %v", msg)
			f.Flush()
		}

		time.Sleep(1 * time.Second) // 模拟生成延迟
	}
}

func stream2Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 支持跨域
	w.Header().Set("Content-Type", "text/plain")
	// 不要设置 Content-Length！

	// 确保响应支持 Flush
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(MCPMessage{
		Type: "message", Role: "assistant", Content: "The weather in Shanghai is 22°C and sunny.",
	})
	flusher.Flush() // 👈 关键：强制将缓冲区数据立即发送
}

func mcpStreamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 支持跨域
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	messages := []MCPMessage{
		{Type: "message", Role: "assistant", Content: "I'll check the weather for you."},
		{Type: "tool_call", ToolName: "get_weather", Arguments: map[string]string{"city": "Shanghai"}},
		{Type: "tool_result", ToolName: "get_weather", Result: "22°C, sunny"},
		{Type: "message", Role: "assistant", Content: "The weather in Shanghai is 22°C and sunny."},
		{Type: "done"},
	}

	for _, msg := range messages {
		json.NewEncoder(w).Encode(msg)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(800 * time.Millisecond)
	}
}

func main() {
	http.HandleFunc("/stream", streamHandler)
	http.HandleFunc("/stream2", stream2Handler)
	http.HandleFunc("/mcp-stream", mcpStreamHandler)
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
