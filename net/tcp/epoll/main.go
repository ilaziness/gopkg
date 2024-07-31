package main

import (
	"expvar"
	"log"
	"net"
	"time"
)

// epoll模式，不用标准库的net
// 只支持linux类系统

var connNum = expvar.NewInt("conns_number")
var reqNum = expvar.NewInt("req_num")
var qps = expvar.NewInt("qps")

var epoller *epoll

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
	}()

	epoller, err = MkEpoll()
	if err != nil {
		panic(err)
	}

	go start()
	for {
		conn, e := ln.Accept()
		if e != nil {
			log.Println(err)
			return
		}

		if err := epoller.Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
		connNum.Add(1)
	}
}

func start() {
	var buf = make([]byte, 20)
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
			if _, err := conn.Read(buf); err != nil {
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}
				// 客户端连接关闭，会走到这里
				log.Println("conn close")
				conn.Close()
			}
			conn.Write([]byte("server hello"))

			reqNum.Add(1)
		}
	}
}
