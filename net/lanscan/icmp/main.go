package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	// 获取局域网的子网掩码和网关
	gateway, subnet, err := getLocalNetworkInfo()
	if err != nil {
		fmt.Println("Error getting local network info:", err)
		return
	}
	fmt.Printf("Local network: %s/%s\n", gateway, subnet)

	// 扫描局域网中的主机
	scanNetwork(gateway, subnet)
}

// getLocalNetworkInfo 获取局域网的网关和子网掩码
func getLocalNetworkInfo() (string, string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", "", err
	}
	for _, iface := range ifaces {
		if iface.Name != "以太网" {
			continue
		}
		// 排除虚拟网卡
		if isVirtualInterface(iface) {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue // 接口未启用
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // 跳过环回接口
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", "", err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					size, _ := ipnet.Mask.Size()
					return ipnet.IP.String(), strconv.Itoa(size), nil
				}
			}
		}
	}

	return "", "", fmt.Errorf("no suitable network interface found")
}

// isVirtualInterface 判断是否为虚拟网卡
func isVirtualInterface(iface net.Interface) bool {
	// 根据接口名称判断是否为虚拟网卡
	virtualPrefixes := []string{"veth", "docker", "virbr", "lxc", "vboxnet"}
	for _, prefix := range virtualPrefixes {
		if iface.Name == prefix || strings.HasPrefix(iface.Name, prefix) {
			return true
		}
	}
	// 根据标志位判断是否为虚拟网卡
	if iface.Flags&net.FlagPointToPoint != 0 {
		return true
	}
	return false
}

// scanNetwork 扫描局域网中的主机
func scanNetwork(gateway string, subnet string) {
	// 创建ICMP套接字
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println("Error creating ICMP socket:", err)
		return
	}
	defer conn.Close()

	// 构造ICMP回显请求消息
	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("Hello, World!"),
		},
	}
	wm, _ := message.Marshal(nil)

	// 解析子网
	_, ipnet, err := net.ParseCIDR(fmt.Sprintf("%s/%s", gateway, subnet))
	if err != nil {
		fmt.Println("Error parsing subnet:", err)
		return
	}

	// 遍历子网中的IP地址并发送ICMP请求
	for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		if ip.String() == gateway {
			continue // 跳过网关
		}

		fmt.Printf("Scanning %s...\n", ip)
		_, err := conn.WriteTo(wm, &net.IPAddr{IP: ip})
		if err != nil {
			fmt.Printf("Error sending ICMP request to %s: %v\n", ip, err)
			continue
		}

		// 设置超时时间
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		buffer := make([]byte, 1024)
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				continue // 超时，未收到响应
			}
			fmt.Printf("Error receiving ICMP response from %s: %v\n", ip, err)
			continue
		}

		// 解析ICMP响应
		rm, err := icmp.ParseMessage(1, buffer[:n])
		if err != nil {
			fmt.Printf("Error parsing ICMP response from %s: %v\n", ip, err)
			continue
		}

		if rm.Type == ipv4.ICMPTypeEchoReply {
			fmt.Printf("Host %s is up\n", addr.String())
		}
	}
}

// incrementIP 增加IP地址
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
