package main

import (
	"log"
	"net"
	"time"
)

func test2() {
	// 1. 客户端已连接一直没有发送数据，服务端read会阻塞，windows会返回EOF错误
	// 2. 服务端read数据比读取缓冲长，每次读取填满缓冲后返回数据，读取全部数据需要循环读取
	// output:
	// 2024/07/25 16:41:19 server2.go:52: server receive: 5 bytes, hello
	// 2024/07/25 16:41:19 server2.go:52: server receive: 5 bytes,  i am
	// 2024/07/25 16:41:19 server2.go:52: server receive: 5 bytes,  clie
	// 2024/07/25 16:41:19 server2.go:52: server receive: 3 bytes, nt2

	// 3. 客户端已关闭
	// 3.1 有数据服务端读，可以把数据全部读出，最后读取完成再读取会返回EOF错误
	// output:
	// 2024/07/25 17:26:02 server2.go:33: simple server listen on [::]:8080
	// 2024/07/25 17:26:03 server2.go:41: accept new client connect
	// 2024/07/25 17:26:04 server2.go:71: client closed
	// 2024/07/25 17:26:05 server2.go:59: server receive: 5 bytes, hello
	// 2024/07/25 17:26:05 server2.go:59: server receive: 5 bytes,  i am
	// 2024/07/25 17:26:05 server2.go:59: server receive: 5 bytes,  clie
	// 2024/07/25 17:26:05 server2.go:59: server receive: 3 bytes, nt2
	// 2024/07/25 17:26:05 server2.go:56: EOF
	// 3.2 无数据读，返回EOF错误
	go server2()
	time.Sleep(time.Second)
	client2()
	time.Sleep(time.Second * 5)
}

func server2() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()
	log.Println("simple server listen on", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("accept new client connect")

		// data := make([]byte, 20)
		// n, err := conn.Read(data)
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// log.Printf("server receive: %d bytes, %s\n", n, data[:n])

		time.Sleep(time.Second * 2)
		for {
			data := make([]byte, 5)
			n, err := conn.Read(data)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("server receive: %d bytes, %s\n", n, data[:n])
		}
	}
}

func client2() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		log.Println("client closed")
	}()
	defer conn.Close()

	_, err = conn.Write([]byte("hello i am client2"))
	if err != nil {
		log.Println(err)
		return
	}
	//time.Sleep(time.Second)
}
