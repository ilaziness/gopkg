package main

import (
	"log"
	"net"
	"time"
)

func simple() {
	go simpleServer()
	time.Sleep(time.Second)
	simpleClient()
	time.Sleep(time.Second * 1)
}

func simpleServer() {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 9981})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("server local: <%s>\n", listener.LocalAddr())

	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			log.Println(err)
		}
		log.Printf("server client: <%s> %s\n", remoteAddr, data[:n])
		_, err = listener.WriteToUDP([]byte("server hello"), remoteAddr)
		if err != nil {
			log.Println(err)
		}
	}
}

func simpleClient() {
	sip := net.ParseIP("127.0.0.1")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: sip, Port: 9981}

	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	conn.Write([]byte("client hello"))
	data := make([]byte, 256)
	n, err := conn.Read(data)
	if err != nil {
		log.Println(err)
	}
	log.Printf("client remote: <%s>\n", conn.RemoteAddr())
	log.Printf("client local: <%s>\n", conn.LocalAddr())
	log.Printf("client server message: <%s> \n", data[:n])
}
