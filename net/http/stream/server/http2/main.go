// 简化版（使用 Go 自动生成证书用于测试）
// 生成测试证书:
// bash: go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost
// powershell: go run "$(go env GOROOT)/src/crypto/tls/generate_cert.go" --host localhost
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func streamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ndjson")
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许跨域

	flusher := w.(http.Flusher)
	for i := 0; i < 10; i++ {
		json.NewEncoder(w).Encode(map[string]string{
			"content": fmt.Sprintf("Token %d (via %s)", i, r.Proto),
		})
		flusher.Flush()
		time.Sleep(400 * time.Millisecond)
	}
}

func main() {
	http.HandleFunc("/stream", streamHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
		w.Write([]byte(r.Proto))
	})
	log.Println("Open https://localhost:8443/stream in browser (accept insecure cert)")
	log.Fatal(http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil))
}
