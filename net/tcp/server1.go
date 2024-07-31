package main

import (
	"log"
	"net"
	"time"
)

func test1() {
	go server1()
	time.Sleep(time.Second)

	// client1()
	// time.Sleep(time.Second)

	// 服务端Accept来不急处理，连接缓冲已满的情况下客户端连接服务器，服务端Accept前sleep 5秒来模拟
	// windows的表现是客户端连接Dial会直接返回refused错误
	// ubuntu linux表现是Dial会阻塞，直到服务器可连接后会继续连接，不会返回错误
	// 服务端accept一个连接，客户端就新增一个连接成功
	// outout:
	// ..........
	// accept new client connect
	// connect server ok 4129
	// accept new client connect
	// connect server ok 4130
	// accept new client connect
	// connect server ok 4131
	// ...........
	batchClient1()
	time.Sleep(time.Second * 60)

	// 服务端未启动或不可用
	// output:
	// 2024/07/25 16:01:16 server1.go:41: dial tcp 127.0.0.1:8080: connectex: No connection could be made because the target machine actively refused it.

}

func server1() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()
	log.Println("simple server listen on", ln.Addr())

	for {
		time.Sleep(time.Second * 5)
		_, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("accept new client connect")

		// handle conn
		//go handleConn(conn)
	}
}

func client1() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	log.Println("connect server ok")
}

func batchClient1() {
	for i := 1; i <= 6000; i++ {
		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		log.Println("connect server ok", i)
	}
}
