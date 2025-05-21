package main

import (
	"fmt"
	"net/http"
	"time"
)

// https://www.ruanyifeng.com/blog/2017/05/server-sent_events.html

func main() {
	http.HandleFunc("/events", sseHandler)
	http.HandleFunc("/long", longHold)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// 设置 SSE 所需的 HTTP 头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 支持跨域

	// 模拟实时数据推送
	c := 0
	for c < 10 {
		// 每一条消息用\n\n分割
		// 一条消息的每一行格式：[field]: value\n
		// field取值：
		//   data：消息内容
		//   id：消息ID
		//   retry：重试时间
		//   event：事件名称，默认是message，也就是客户端的onmessage事件，可以自定义然后用addEventListener来监听
		fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))
		flusher.Flush() // 立即将数据发送给客户端
		time.Sleep(1 * time.Second)
		c++
	}
	go sendLater(w)
	fmt.Println("sseHandler done")
}

// sendLater 模拟连接数据输出完毕后，再发送数据
// 数据输出完了后，连接已经关闭，不能再写入数据，会panic
func sendLater(w http.ResponseWriter) {
	// flusher, _ := w.(http.Flusher)
	time.Sleep(20 * time.Second)
	fmt.Println("sseHandler send later")
	_, err := fmt.Fprintf(w, "send later data: %s\n\n", time.Now().Format(time.RFC3339))
	if err != nil {
		fmt.Println(err)
	}
	// 数据输出完了后，连接已经关闭，不能再写入数据，会panic
	//flusher.Flush()
}

// longHold 模拟长连接，每10秒发送一次数据
func longHold(w http.ResponseWriter, r *http.Request) {
	// 设置 SSE 所需的 HTTP 头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 支持跨域

	flusher, _ := w.(http.Flusher)

	// 重试间隔，1000毫秒
	// 只有意外中断才会重试连接，服务器关闭或者重启不会自动重新连接
	_, err := fmt.Fprintf(w, "retry: 1000\n")
	if err != nil {
		fmt.Println(err)
	}
	_, err = fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))
	if err != nil {
		fmt.Println(err)
	}

	flusher.Flush()

	for {
		_, err = fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))
		if err != nil {
			fmt.Printf("Client disconnected: %v\n", err)
			break
		}
		flusher.Flush()
		time.Sleep(10 * time.Second)
	}
}
