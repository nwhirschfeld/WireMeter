package main

import (
	"encoding/binary"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
)

func getPacketN(i uint64) gopacket.Packet {
	packetData, err := generateRawEthRequestPacket(i)
	if err != nil {
		log.Fatal(err)
	}
	return gopacket.NewPacket(packetData.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
}

func generateRawEthRequestPacket(counter uint64) (gopacket.SerializeBuffer, error) {
	ethernetLayer := layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x42, 0x23, 0x69, 0x42, 0x23, 0x69},
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, // Broadcast MAC address
		EthernetType: layers.EthernetType(0x1337),
	}

	b := make([]byte, binary.MaxVarintLen64)
	binary.LittleEndian.PutUint64(b, uint64(counter))
	rawData := gopacket.Payload(b)

	// Serialize the packet
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{ComputeChecksums: true}
	err := gopacket.SerializeLayers(buf, opts, &ethernetLayer, rawData)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
