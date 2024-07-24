package main

import (
	"log"
	"net"
	"time"
)

// 等价的客户端和服务器
// 两个服务器通信的例子，互为客户端和服务器，我们可以将发送的一方称之为源地址，发送的目的地一方称之为目标地址。

func sameCS() {
	addr1 := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9983}
	addr2 := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9984}

	go func() {
		listener1, err := net.ListenUDP("udp", addr1)
		if err != nil {
			log.Println(err)
			return
		}
		go read(listener1)
		time.Sleep(time.Second * 1)
		listener1.WriteToUDP([]byte("ping to #2: "+addr2.String()), addr2)
	}()

	go func() {
		listener2, err := net.ListenUDP("udp", addr2)
		if err != nil {
			log.Println(err)
			return
		}
		go read(listener2)
		time.Sleep(time.Second * 1)
		listener2.WriteToUDP([]byte("ping to #1: "+addr1.String()), addr1)
	}()

	time.Sleep(time.Second * 2)
}

func read(conn *net.UDPConn) {
	for {
		data := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s \n", err)
		}
		log.Printf("receive %s from <%s>\n", data[:n], remoteAddr)
	}
}
