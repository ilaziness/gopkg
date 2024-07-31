package main

import (
	"expvar"
	"net/http"
	"strconv"
	"sync/atomic"
)

// 自定义http server

// 自定义指标数据类型，需要实现Var接口
type MyVar struct {
	number int64
}

func NewMyVar(name string) *MyVar {
	v := new(MyVar)
	expvar.Publish(name, v)
	return v
}

func (v *MyVar) String() string {
	return strconv.FormatInt(atomic.LoadInt64(&v.number), 10)
}

func (v *MyVar) Add(delta int64) {
	atomic.AddInt64(&v.number, delta)
}

var (
	// 访问指标地址，会输出TotalRequest的数据
	totalReq = expvar.NewInt("TotalRequest")
	number   = NewMyVar("number")
)

func testCustom() {
	// http://127.0.0.1:8881/mydebug/vars
	mux := http.NewServeMux()
	mux.Handle("/mydebug/vars", expvar.Handler())

	// 自定义暴露指标数据
	mux.HandleFunc("/tr", func(w http.ResponseWriter, r *http.Request) {
		totalReq.Add(1)
		number.Add(2)
		w.Write([]byte("total request"))
	})

	http.ListenAndServe(":8881", mux)
}
