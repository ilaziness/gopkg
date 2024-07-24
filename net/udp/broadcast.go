package main

import (
	"log"
	"net"
	"time"
)

// 广播
// 广播都是限制在局域网中的，所以只有局域网有效

func boradcast() {
	go server()
	time.Sleep(time.Second)
	client()

	time.Sleep(time.Second * 2)
}

func server() {
	listener, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 9986})
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("server local: <%s> \n", listener.LocalAddr().String())

	data := make([]byte, 1024)
	for {
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s", err)
		}
		log.Printf("boradcast server receive: <%s> %s\n", remoteAddr, data[:n])
		_, err = listener.WriteToUDP([]byte("world"), remoteAddr)
		if err != nil {
			log.Println(err)
		}
	}
}

func client() {
	// 10.0.2.255 广播地址根据本机子网掩码得出
	ip := net.ParseIP("10.0.2.255")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: 9986}
	// 发送方连接广播地址可以使用ListenPacket,ListenUDP,DialUDP和Dial
	// https://github.com/aler9/howto-udp-broadcast-golang
	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	_, err = conn.WriteToUDP([]byte("hello"), dstAddr)
	if err != nil {
		log.Println(err)
	}
	data := make([]byte, 1024)
	// 设置读取超时时间
	conn.SetReadDeadline(time.Now().Add(time.Second))
	n, _, err := conn.ReadFrom(data)
	if err != nil {
		log.Println(err)
	}
	log.Printf("boradcast client read: <%s> %s\n", dstAddr, data[:n])
}
