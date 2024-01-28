package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
)

func getPacketN(i int) gopacket.Packet {
	srcMAC, _ := net.ParseMAC(intToMACAddress(int64(i)))
	srcIP := net.ParseIP("192.168.1.1")
	targetIP := net.ParseIP("192.168.1.2")
	packetData, err := generateArpRequestPacket(srcMAC, srcIP, targetIP)
	if err != nil {
		log.Fatal(err)
	}
	return gopacket.NewPacket(packetData.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
}

func intToMACAddress(num int64) string {
	// Convert the integer to a MAC address format
	macAddress := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		(num>>40)&0xFF,
		(num>>32)&0xFF,
		(num>>24)&0xFF,
		(num>>16)&0xFF,
		(num>>8)&0xFF,
		num&0xFF,
	)

	return macAddress
}

func generateArpRequestPacket(srcMAC net.HardwareAddr, srcIP, targetIP net.IP) (gopacket.SerializeBuffer, error) {
	// Create the ARP request packet
	arpRequest := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   srcMAC,
		SourceProtAddress: srcIP,
		DstHwAddress:      net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, // Will be populated by the target's MAC address during transmission
		DstProtAddress:    targetIP,
	}

	ethernetLayer := layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, // Broadcast MAC address
		EthernetType: layers.EthernetTypeARP,
	}

	// Serialize the packet
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
	err := gopacket.SerializeLayers(buf, opts, &ethernetLayer, &arpRequest)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
