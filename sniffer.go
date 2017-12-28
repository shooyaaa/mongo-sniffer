package main

import (
	"flag"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

func main() {

	enth := flag.String("i", "", "network card ")
	test := flag.String("test", "", "network card ")
	flag.Parse()
	fmt.Println("enth ", *enth, *test)
	if *enth == "" {
		usage()
	} else {
		sniffer(*enth)
	}

}

var (
	buffer      int32 = 1024
	promiscuous bool  = false
	err         error
	timeout     time.Duration = 30 * time.Second
	handle      *pcap.Handle
)

func sniffer(networkCard string) {
	fmt.Println(networkCard)
	handle, err := pcap.OpenLive(networkCard, buffer, promiscuous, timeout)

	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println(packet)
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
