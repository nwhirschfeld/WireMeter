package main

import (
	"container/ring"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/assert"
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

	packet1 := createTestPacket(1337)
	packet2 := createTestPacket(1337)

	store.addPacket(1337, packet1.Metadata().Timestamp)

	meas, err := store.removePacket(packet2)

	assert.Nil(t, err)
	assert.Equal(t, MeasurementReason(ReasonResolvedPacket), meas.reason)
	assert.Equal(t, true, isRingBufferUndefined(store.Buffer))
}

func TestPacketStore_AddPacket(t *testing.T) {
	size := 5
	store := NewMessageStore(size)

	packet1 := createTestPacket(1337)
	packet2 := createTestPacket(1337)

	meas, err := store.addPacket(1337, packet1.Metadata().Timestamp)
	assert.Nil(t, err)
	assert.Nil(t, meas)

	meas, err = store.addPacket(1337, packet2.Metadata().Timestamp)
	assert.Nil(t, err)
	assert.Nil(t, meas)

	// Add more test cases as needed
}

func createTestPacket(counter uint64) gopacket.Packet {
	packetData, _ := generateRawEthRequestPacket(counter)
	return gopacket.NewPacket(packetData.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
}

func createTestPacketWithTimestamp(counter uint64, timestamp time.Time) gopacket.Packet {
	packet := createTestPacket(counter)
	packet.Metadata().Timestamp = timestamp
	return packet
}
