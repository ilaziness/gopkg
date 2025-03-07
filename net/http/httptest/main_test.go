package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 作用：
// 创建模拟 HTTP 服务器：用于测试客户端代码。
// 创建模拟 HTTP 请求：用于测试服务器端代码。
// 捕获 HTTP 响应：方便检查响应内容。

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

// TestServer 测试服务器
func TestServer(t *testing.T) {
	// 创建一个模拟 HTTP 服务器
	server := httptest.NewServer(http.HandlerFunc(helloHandler))
	defer server.Close() // 测试结束后关闭服务器
	println(server.URL)

	// 发送请求到模拟服务器
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		t.Errorf("期望状态码 200，实际状态码 %d", resp.StatusCode)
	}

	// 检查响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("读取响应失败: %v", err)
	}
	expected := "Hello, World!\n"
	if string(body) != expected {
		t.Errorf("期望响应 %q，实际响应 %q", expected, string(body))
	}
}

func helloHandler2(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	println(name)
	fmt.Fprintf(w, "Hello, %s!", name)
}

// TestClient 测试客户端
func TestClient(t *testing.T) {
	// 创建一个模拟请求
	req := httptest.NewRequest("GET", "/hello?name=John", nil)

	// 创建一个 ResponseRecorder 来捕获响应
	rr := httptest.NewRecorder()

	// 调用处理函数
	helloHandler2(rr, req)

	// 检查响应状态码
	if rr.Code != http.StatusOK {
		t.Errorf("期望状态码 200，实际状态码 %d", rr.Code)
	}

	// 检查响应内容
	expected := "Hello, John!"
	if rr.Body.String() != expected {
		t.Errorf("期望响应 %q，实际响应 %q", expected, rr.Body.String())
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 记录日志
		fmt.Println("请求路径:", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// TestMiddleware 测试中间件
func TestMiddleware(t *testing.T) {
	// 创建一个模拟请求
	req := httptest.NewRequest("GET", "/hello", nil)

	// 创建一个 ResponseRecorder 来捕获响应
	rr := httptest.NewRecorder()

	// 创建中间件和处理函数
	handler := loggingMiddleware(http.HandlerFunc(helloHandler))

	// 调用中间件
	handler.ServeHTTP(rr, req)

	// 检查响应状态码
	if rr.Code != http.StatusOK {
		t.Errorf("期望状态码 200，实际状态码 %d", rr.Code)
	}

	// 检查响应内容
	expected := "Hello, World!\n"
	if rr.Body.String() != expected {
		t.Errorf("期望响应 %q，实际响应 %q", expected, rr.Body.String())
	}
}
