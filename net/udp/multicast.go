package main

import (
	"log"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

// 多播(组播)
// 组播理论上可以用在广域网，实际需要路由设备支持，运营商一般会关闭，所以实际广域网做不了组播

func multicast() {
	general()
	stdlib()
}

// 通用多播
func general() {
	// 网络接口, linux ifconfig
	en4, err := net.InterfaceByName("enp0s3")
	// windows
	// en4, err := net.InterfaceByName("以太网")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(en4)
	// 多播组
	// 224.0.0.250 固定的组播地址
	// 组播地址是iana的保留地址
	group := net.IPv4(224, 0, 0, 250)

	// 侦听，绑定端口
	c, err := net.ListenPacket("udp4", "0.0.0.0:1024")
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	// 应用加入多播组
	p := ipv4.NewPacketConn(c)
	if err := p.JoinGroup(en4, &net.UDPAddr{IP: group}); err != nil {
		log.Println(err)
		return
	}

	// 更多控制
	// 不支持windows
	if err := p.SetControlMessage(ipv4.FlagDst, true); err != nil {
		log.Println(err)
		return
	}

	//接收数据包
	go func() {
		b := make([]byte, 1500)
		for {
			n, cm, src, err := p.ReadFrom(b)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("received1: %s from <%s>\n", b[:n], src)
			// 需要设置SetControlMessage
			if cm.Dst.IsMulticast() {
				// 检查包是否同一个组的包
				if !cm.Dst.Equal(group) {
					log.Println("Unknown group")
					continue
				}
				log.Printf("received: %s from <%s>\n", b[:n], src)
				_, err = p.WriteTo([]byte("world"), cm, src)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	// 发送组播数据包，上面接收能打印出来
	dst := &net.UDPAddr{IP: group, Port: 1024}
	if err := p.SetMulticastInterface(en4); err != nil {
		log.Println(err)
		return
	}
	p.SetMulticastTTL(5)
	if _, err := p.WriteTo([]byte("hello"), nil, dst); err != nil {
		log.Println(err)
		return
	}

	generalClient()
	time.Sleep(time.Second * 2)
}

func generalClient() {
	ip := net.ParseIP("224.0.0.250")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: 1024}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	conn.Write([]byte("hello2 1024 server"))
	log.Printf("generalClient <%s>\n", conn.RemoteAddr())
}

// 标准库多播
// 和上面比更少控制选项
// 标准库实现仅仅适用于简单的小应用
func stdlib() {
	// 224.0.0.250多播地址
	addr, err := net.ResolveUDPAddr("udp", "224.0.0.250:9985")
	if err != nil {
		log.Println(err)
	}
	listener, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Local: <%s> \n", listener.LocalAddr().String())

	go func() {
		data := make([]byte, 1024)
		for {
			n, remoteAddr, err := listener.ReadFromUDP(data)
			if err != nil {
				log.Printf("error during read: %s", err)
			}
			log.Printf("stdlib receive <%s> %s\n", remoteAddr, data[:n])
		}
	}()
	stdlibClient()
	time.Sleep(time.Second * 2)
}

// 客户端
func stdlibClient() {
	ip := net.ParseIP("224.0.0.250")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: 9985}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	conn.Write([]byte("hello"))
	log.Printf("stdlibClient <%s>\n", conn.RemoteAddr())
}
