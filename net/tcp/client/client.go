package main

import (
	"log"
	"net"
	"time"
)

func main() {
	addr := "127.0.0.1:9000"
	clientNum := 10
	for n := 0; n < clientNum; n++ {
		go client(n, addr)
	}

	time.Sleep(time.Second * 30)
}

func client(n int, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	deadline := time.Now().Add(time.Second * 10)
	for {
		if time.Now().After(deadline) {
			log.Printf("client %d exit\n", n)
			return
		}
		time.Sleep(time.Millisecond * 500)
		_, err = conn.Write([]byte("client hello"))
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("client %d sent\n", n)
		data := make([]byte, 20)
		nr, err := conn.Read(data)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("client %d receive: %s\n", n, data[:nr])
	}
}
