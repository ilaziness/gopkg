// client.go
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type StreamMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Index   int    `json:"index,omitempty"`
}

func main() {
	//resp, err := http.Get("http://localhost:8080/stream")
	resp, err := http.Get("http://localhost:8080/mcp-stream")
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server returned error: %s", resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		var msg StreamMessage
		if err := json.Unmarshal(line, &msg); err != nil {
			log.Printf("Failed to parse line: %s, error: %v", string(line), err)
			continue
		}
		fmt.Printf("Received: %+v\n", msg)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading stream:", err)
	}
}
