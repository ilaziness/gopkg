package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	manuf "github.com/timest/gomanuf"
)

// listenARP 监听arp响应
// ifName 网络接口名称
func listenARP(ctx context.Context, ifName string, macChan chan LanIpInfo) {
	handle, err := pcap.OpenLive(ifName, 1024, false, 10*time.Second)
	if err != nil {
		log.Fatal("pcap打开失败:", err)
	}
	defer handle.Close()
	handle.SetBPFFilter("arp")
	ps := gopacket.NewPacketSource(handle, handle.LinkType())

	for {
		select {
		case <-ctx.Done():
			return
		case p := <-ps.Packets():
			arp := p.Layer(layers.LayerTypeARP).(*layers.ARP)
			if arp.Operation == 2 {
				mac := net.HardwareAddr(arp.SourceHwAddress)
				m := manuf.Search(mac.String())
				macChan <- LanIpInfo{
					IP:           ParseIP(arp.SourceProtAddress).String(),
					Mac:          mac,
					Manufacturer: m,
				}
				// if strings.Contains(m, "Apple") {
				// 	go sendMdns(ParseIP(arp.SourceProtAddress), mac)
				// } else {
				// 	go sendNbns(ParseIP(arp.SourceProtAddress), mac)
				// }
			}
		}
	}
}

// 往目标ip发送arp包
// localIp: 本机ip，localIfName：本地网络接口名称，localMac：本地网络接口MAC地址, lanIp：内网目标ip
func sendArpPackage(localIp net.IP, localIfName string, localMac net.HardwareAddr, lanIp IP) {
	srcIp := localIp.To4()
	dstIp := net.ParseIP(lanIp.String()).To4()
	if srcIp == nil || dstIp == nil {
		log.Fatal("ip 解析出问题")
	}
	// 以太网首部
	// EthernetType 0x0806  ARP
	ether := &layers.Ethernet{
		SrcMAC:       localMac,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	a := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     uint8(6),
		ProtAddressSize:   uint8(4),
		Operation:         uint16(1), // 0x0001 arp request 0x0002 arp response
		SourceHwAddress:   localMac,
		SourceProtAddress: srcIp,
		DstHwAddress:      net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    dstIp,
	}

	buffer := gopacket.NewSerializeBuffer()
	var opt gopacket.SerializeOptions
	gopacket.SerializeLayers(buffer, opt, ether, a)
	outgoingPacket := buffer.Bytes()

	handle, err := pcap.OpenLive(localIfName, 2048, false, 30*time.Second)
	if err != nil {
		log.Fatal("pcap打开失败:", err)
	}
	defer handle.Close()

	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		log.Fatal("发送arp数据包失败..")
	}
}
