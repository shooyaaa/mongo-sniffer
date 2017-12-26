package main

import (
	"fmt"
	"github.com/google/gopacket/pcap"
	"log"
)

func main() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("All Devices: ")

	for _, device := range devices {
		fmt.Println("\nName : ", device.Name)
		if device.Address {
			for _, address := range device.Addresses {
				fmt.Println("IP : ", address.IP)
			}
		}
	}
}
