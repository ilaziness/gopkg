package main

// 标准库的server

import (
	_ "expvar"
	"net/http"
)

// 路径: /debug/vars
// http://127.0.0.1:8888/debug/vars
func testStd() {
	http.ListenAndServe(":8888", nil)
}
