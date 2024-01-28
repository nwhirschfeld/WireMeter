package main

import (
	"container/ring"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestNewMessageStore(t *testing.T) {
	size := 5
	store := NewMessageStore(size)

	assert.NotNil(t, store.Buffer)
	assert.Equal(t, size, store.size)
}

func isRingBufferUndefined(r *ring.Ring) bool {
	for i := 0; i < r.Len(); i++ {
		assert.Nil(nil, r.Value)
		r = r.Next()
	}
	return true
}

func TestPacketStore_RemovePacket(t *testing.T) {
	size := 5
	store := NewMessageStore(size)

	packet1 := createTestPacket("00:11:22:33:44:55")
	packet2 := createTestPacket("00:11:22:33:44:55")

	store.addPacket(packet1)

	meas, err := store.removePacket(packet2)

	assert.Nil(t, err)
	assert.Equal(t, MeasurementReason(ReasonResolvedPacket), meas.reason)
	assert.Equal(t, true, isRingBufferUndefined(store.Buffer))
}

func TestPacketStore_AddPacket(t *testing.T) {
	size := 5
	store := NewMessageStore(size)

	packet1 := createTestPacket("00:11:22:33:44:55")
	packet2 := createTestPacket("11:22:33:44:55:66")

	meas, err := store.addPacket(packet1)
	assert.Nil(t, err)
	assert.Nil(t, meas)

	meas, err = store.addPacket(packet2)
	assert.Nil(t, err)
	assert.Nil(t, meas)

	// Add more test cases as needed
}

func createTestPacket(srcMACString string) gopacket.Packet {
	srcMAC, _ := net.ParseMAC(srcMACString)
	srcIP := net.ParseIP("192.168.1.1")
	targetIP := net.ParseIP("192.168.1.2")
	packetData, _ := generateArpRequestPacket(srcMAC, srcIP, targetIP)
	return gopacket.NewPacket(packetData.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
}

func createTestPacketWithTimestamp(srcMAC string, timestamp time.Time) gopacket.Packet {
	packet := createTestPacket(srcMAC)
	packet.Metadata().Timestamp = timestamp
	return packet
}

func parseMAC(mac string) net.HardwareAddr {
	parsedMAC, _ := net.ParseMAC(mac)
	return parsedMAC
}
