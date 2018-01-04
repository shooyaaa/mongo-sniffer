package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

func main() {

	enth := flag.String("i", "", "network card ")
	filter := flag.String("filter", "", "tcpdump rules")
	flag.Parse()
	if *enth != "" {
		sniffer(enth, filter)
	}

}

var (
	buffer      int32 = 1024
	promiscuous bool  = false
	err         error
	timeout     time.Duration = 30 * time.Second
	handle      *pcap.Handle
)

type MyPacket struct {
	ethernetType string
	SrcMac       string
	DstMac       string
	SrcIp        int32
	DstIp        int32
	SrcPort      int32
	DstPort      int32
	AppData      []byte
}

type IpPair struct {
	SrcIp string
	DstIp string
}

type PortPair struct {
	SrcPort layers.TCPPort
	DstPort layers.TCPPort
}

func sniffer(networkCard *string, filter *string) {
	handle, err := pcap.OpenLive(*networkCard, buffer, promiscuous, timeout)

	if *filter != "" {
		fmt.Println("filter is ", *filter)
		handle.SetBPFFilter(*filter)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		//fmt.Println(packet)
		parsePacket(packet)
	}
}

func usage() {
	fmt.Println("mongo-sniffer -i enth")
	listDevices()
}

func listDevices() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	for _, device := range devices {
		if len(device.Addresses) > 0 {
			fmt.Println("\nName : ", device.Name)
			for _, address := range device.Addresses {
				fmt.Println("IP : ", address.IP)
			}
		}
	}
}

func parsePacket(packet gopacket.Packet) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		//fmt.Println("IPv4 layer detectead")
		ip, _ := ipLayer.(*layers.IPv4)
		ipv4Handler(packet, IpPair{ip.SrcIP.String(), ip.DstIP.String()})
	}
}

func ipv4Handler(packet gopacket.Packet, ipPair IpPair) {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		tcpHandler(packet, ipPair, PortPair{tcp.SrcPort, tcp.DstPort})
	}
}

func tcpHandler(packet gopacket.Packet, ipPair IpPair, portPair PortPair) {
	//fmt.Println("port ", portPair)
	appLayer := packet.ApplicationLayer()
	if appLayer != nil {
		buf := bytes.NewReader(appLayer.Payload())
		op := MongoOp{buf, ipPair, portPair}
		op.decode()
	}
}
