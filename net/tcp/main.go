package main

import (
	"log"
	"net"
	"time"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//simple()
	//output:
	// 2024/07/25 15:40:44 main.go:30: simple server listen on [::]:8080
	// 2024/07/25 15:40:45 main.go:51: simple server receive: 5 bytes, hello
	// 2024/07/25 15:40:45 main.go:75: simple client receive: 5 bytes, world

	//test1()
	//test2()
	test3()
}

func simple() {
	go simpleServer()
	time.Sleep(time.Second)
	simpleClient()

	time.Sleep(time.Second * 1)
}

func simpleServer() {
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

		// handle conn
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	data := make([]byte, 20)
	n, err := conn.Read(data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("simple server receive: %d bytes, %s\n", n, data[:n])
	conn.Write([]byte("world"))
}

func simpleClient() {
	// net.DialTimeout 设置超时时间
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("hello"))
	if err != nil {
		log.Println(err)
		return
	}
	data := make([]byte, 20)
	n, err := conn.Read(data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("simple client receive: %d bytes, %s\n", n, data[:n])
}
