package main

import (
	"log"
	"net"
	"time"
)

func test3() {
	// 写阻塞，服务端没有读取，客户端一直写直到缓冲区满
	// 客户端写缓冲满了之后不可写，会一直阻塞，等读取完一定数据（需要读取多少windows、linux都不一定）后会再次写入

	go server31()
	time.Sleep(time.Second)
	client31()
	time.Sleep(time.Second * 20)
}

func server31() {
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
		// 设置socket属性，需要用类型断言net.TCPConn
		// conn.(*net.TCPConn).SetKeepAlive(true)

		for {
			time.Sleep(time.Second * 1)
			data := make([]byte, 20)
			n, err := conn.Read(data)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("server receive: %d bytes, %s\n", n, data[:n])
		}
	}
}

func client31() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		log.Println("client closed")
	}()
	defer conn.Close()

	total := 0
	// conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
	for i := 1; i <= 150000; i++ {
		// write 满了后会阻塞
		n, err := conn.Write([]byte("hellohellohellohello"))
		if err != nil {
			log.Println("write error:", err)
			continue
		}
		total += n
		log.Printf("client31 write total %d bytes, %d", total, i)
	}
}
