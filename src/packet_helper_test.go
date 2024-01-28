package main

import (
	"fmt"
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/assert"
)

func TestGetPacketN(t *testing.T) {
	testCases := []struct {
		input    int
		expected string
	}{
		{1, "00:00:00:00:00:01"},
		{255, "00:00:00:00:00:ff"},
		{0x112233445566, "11:22:33:44:55:66"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input_%d", tc.input), func(t *testing.T) {
			packet := getPacketN(tc.input)
			ethLayer := packet.Layer(layers.LayerTypeEthernet).(*layers.Ethernet)

			assert.Equal(t, tc.expected, ethLayer.SrcMAC.String())
			assert.Equal(t, "ff:ff:ff:ff:ff:ff", ethLayer.DstMAC.String())
			assert.Equal(t, layers.EthernetTypeARP, ethLayer.EthernetType)
		})
	}
}

func TestIntToMACAddress(t *testing.T) {
	testCases := []struct {
		input    int64
		expected string
	}{
		{1, "00:00:00:00:00:01"},
		{255, "00:00:00:00:00:ff"},
		{0x112233445566, "11:22:33:44:55:66"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Input_%d", tc.input), func(t *testing.T) {
			result := intToMACAddress(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGenerateArpRequestPacket(t *testing.T) {
	srcMAC := net.HardwareAddr{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	srcIP := net.IP{192, 168, 1, 1}
	targetIP := net.IP{192, 168, 1, 2}

	packetData, err := generateArpRequestPacket(srcMAC, srcIP, targetIP)
	assert.Nil(t, err)

	packet := gopacket.NewPacket(packetData.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
	ethLayer := packet.Layer(layers.LayerTypeEthernet).(*layers.Ethernet)
	arpLayer := packet.Layer(layers.LayerTypeARP).(*layers.ARP)

	assert.Equal(t, srcMAC.String(), ethLayer.SrcMAC.String())
	assert.Equal(t, "ff:ff:ff:ff:ff:ff", ethLayer.DstMAC.String())
	assert.Equal(t, layers.EthernetTypeARP, ethLayer.EthernetType)

	assert.Equal(t, uint16(layers.ARPRequest), arpLayer.Operation)
	assert.Equal(t, srcMAC.String(), net.HardwareAddr(arpLayer.SourceHwAddress).String())
	assert.Equal(t, srcIP.String(), net.IP(arpLayer.SourceProtAddress).String())
	assert.Equal(t, targetIP.String(), net.IP(arpLayer.DstProtAddress).String())
}
