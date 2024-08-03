package main

import (
	"log"
	"net"
)

// workerpool模式
// epool模式加上worker pool

var epoller *epoll
var workerPool *pool

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Println(err)
		return
	}

	workerPool = newPool(1000)
	workerPool.start()
	defer workerPool.close()

	epoller, err = MkEpoll()
	if err != nil {
		panic(err)
	}

	go start()

	for {
		conn, e := ln.Accept()
		if e != nil {
			log.Printf("accept err: %v", e)
			return
		}

		if err := epoller.Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
	}
}

func start() {
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

			workerPool.pushTask(conn)
		}
	}
}

func hancleConn(conn net.Conn) {
	data := make([]byte, 20)
	// read 已关闭无数据的连接会io.EOF错误
	_, err := conn.Read(data)
	if err != nil {
		conn.Close()
		return
	}
	// log.Printf("client: %s \n", conn.RemoteAddr())
	// 200ms模拟处理时间
	// time.Sleep(time.Millisecond * 200)
	conn.Write([]byte("server hello"))
}
