package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/google/gopacket/pcap"
)

// 网卡信息
type NetInterface struct {
	// 名称
	Name string
	// MAC地址
	HardAddr net.HardwareAddr
	// 本机ip
	LocalIP net.IP
	// 网段内网ip列表
	LanIPs []IP
}

// 内网ip信息
type LanIpInfo struct {
	IP string
	// IP Mac地址
	Mac net.HardwareAddr
	// 主机名
	Hostname string
	// 制造商
	Manufacturer string
}

// 网卡接口列表
var nci []NetInterface

var LanIpInfos *sync.Map = &sync.Map{}

// https://haydz.github.io/2020/07/06/Go-Windows-NIC.html
// https://github.com/google/gopacket/issues/456

// 源项目：https://github.com/timest/goscan
func main() {
	// 指定要扫描的的网卡名称
	var scanIfName string
	var err error
	var ifs []net.Interface
	flag.StringVar(&scanIfName, "i", "", "Network interface name")
	flag.Parse()

	ifs, err = net.Interfaces()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	if len(ifs) == 0 {
		slog.Warn("not found net interface")
		os.Exit(1)
	}
	for _, it := range ifs {
		if scanIfName != "" && it.Name != scanIfName {
			continue
		}
		addrs, err := it.Addrs()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		for _, addr := range addrs {
			ip, ok := addr.(*net.IPNet)
			if !ok || ip.IP.IsLoopback() {
				continue
			}
			if ip.IP.To4() != nil {
				nci = append(nci, NetInterface{
					Name:     it.Name,
					HardAddr: it.HardwareAddr,
					LocalIP:  ip.IP,
					LanIPs:   listLanIps(ip),
				})
			}
		}
	}
	if len(nci) == 0 {
		log.Fatalln("network interface is empty")
	}

	cancelctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 接收arp响应
	ipInfoChan := make(chan LanIpInfo)
	for _, it := range nci {
		deviceName := it.Name
		if runtime.GOOS == "windows" {
			deviceName = findDevice(it.LocalIP)
			if deviceName == "" {
				continue
			}
		}
		slog.Info(fmt.Sprintf("listen arp package: %s", it.Name))
		go listenARP(cancelctx, deviceName, ipInfoChan)
	}

	// 发送arp包
	go func() {
		for _, it := range nci {
			slog.Info(fmt.Sprintf("send arp package, interface name: %s", it.Name))
			for _, ip := range it.LanIPs {
				if it.LocalIP.Equal(net.ParseIP(ip.String())) {
					slog.Info("send arp ignore", "ip", ip.String())
					// 略过本机ip
					continue
				}
				deviceName := it.Name
				if runtime.GOOS == "windows" {
					deviceName = findDevice(it.LocalIP)
					if deviceName == "" {
						continue
					}
				}
				// time.Sleep(time.Millisecond * 50)
				// slog.Info(fmt.Sprintf("send arp package %s", ip))
				sendArpPackage(it.LocalIP, deviceName, it.HardAddr, ip)
			}
		}
		slog.Info("send arp over")
	}()

	timeout := time.Second * 9
	ticker := time.NewTicker(time.Second * 3)
	receiveTime := time.Now()
	//FOR:
	for {
		select {
		case info := <-ipInfoChan:
			receiveTime = time.Now()
			ipinfo, ok := LanIpInfos.Load(info.IP)
			if !ok {
				LanIpInfos.Store(info.IP, info)
				// continue
			}
			if ip, ok := ipinfo.(LanIpInfo); ok {
				slog.Info("ip info", "IP", ip.IP, "MAC", ip.Mac.String(), "Manuf", ip.Manufacturer, "Hostname", ip.Hostname)
			}
		case <-ticker.C:
			if time.Since(receiveTime) > timeout {
				//break FOR
			}
		}
	}
	slog.Info("lan scan completed")
}

func findDevice(ip net.IP) string {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		slog.Error(err.Error())
		return ""
	}
	for _, d := range devices {
		for _, address := range d.Addresses {
			if address.IP.Equal(ip) {
				return d.Name
			}

		}
	}
	return ""
}
