// server_h3.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/http3/qlog"
)

func streamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	for i := 0; i < 5; i++ {
		msg := map[string]interface{}{
			"token": fmt.Sprintf("h3-word%d", i),
			"proto": r.Proto, // 应为 "HTTP/3.0"
			"index": i,
		}
		json.NewEncoder(w).Encode(msg)
		flusher.Flush()
		time.Sleep(600 * time.Millisecond)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok 5\n"))
		w.Write([]byte(r.Proto))
	})
	mux.HandleFunc("/stream", streamHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("../.."))))

	certPath := "cert.pem"

	keyPath := "key.pem"

	// chrome: chrome.exe --origin-to-force-quic-on=127.0.0.1:4433 https://127.0.0.1:4433
	go func() {
		server := &http3.Server{
			Handler: mux,
			Addr:    ":443",
			// TLSConfig:  http3.ConfigureTLSConfig(&tls.Config{}),
			QUICConfig: &quic.Config{
				// qlog是QUIC的调试日志格式，用于调试和分析QUIC的连接行为，https://quic-go.net/docs/quic/qlog/
				Tracer: qlog.DefaultConnectionTracer,
			},
		}

		// 默认的Tracer实现会把日志输出到QLOGDIR环境变量指定的目录，目录如果不存在会自动创建
		// 日志文件在线分析：https://qvis.quictools.info/
		os.Setenv("QLOGDIR", "./qlog")

		log.Println("HTTP/3 server running on :443")
		log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
	}()

	// 如果使用443端口，不能直接用127.0.0.1，需要用127.0.0.1:443
	// chrome chrome.exe --origin-to-force-quic-on=127.0.0.1:4431 https://127.0.0.1:4431
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, fmt.Sprintf("ok, http proto=%s", r.Proto))
		})
		// If mux is nil, the http.DefaultServeMux is used.
		log.Println("HTTP/3 server running on :4431")
		log.Fatal(http3.ListenAndServeQUIC(":4431", certPath, keyPath, mux))
	}()

	// 同时启动 HTTP/2 (TCP)
	log.Println("HTTP/2 on TCP :8443")
	log.Fatal(http.ListenAndServeTLS(":8443", certPath, keyPath, mux))
}
