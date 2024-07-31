package main

import (
	"expvar"
	"log"
	"net"
	"time"
)

// 一个连接一个goroutine处理的模式
// 或
// read 一个goroutine, write一个goroutine

var connNum = expvar.NewInt("conns_number")
var reqNum = expvar.NewInt("req_num")
var qps = expvar.NewInt("qps")

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln(err)
	}

	// 计算qps
	go func() {
		var lastReqTotal int64
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for t := range ticker.C {
			_ = t
			total := reqNum.Value()
			qps.Set(total - lastReqTotal)
			lastReqTotal = total

			log.Printf("conn total: %d, qps: %d\n", connNum.Value(), qps.Value())
		}
		// for {
		// 	select {
		// 	case <-ticker.C:
		// 		total := reqNum.Value()
		// 		qps.Set(total - lastReqTotal)
		// 		lastReqTotal = total
		// 	}
		// }
	}()

	for {
		conn, e := ln.Accept()
		if e != nil {
			log.Println(err)
			return
		}

		go handleConn(conn)
		connNum.Add(1)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	data := make([]byte, 20)
	for {
		// read 已关闭无数据的连接会io.EOF错误
		_, err := conn.Read(data)
		if err != nil {
			// log.Println(err)
			return
		}
		// 200ms模拟处理时间
		// time.Sleep(time.Millisecond * 200)
		conn.Write([]byte("server hello"))
		reqNum.Add(1)
	}
}
