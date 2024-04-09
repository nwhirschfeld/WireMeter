package main

import (
	"container/ring"
	"encoding/binary"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"sync"
	"time"
)

type PacketStore struct {
	Buffer *ring.Ring
	size   int
	mu     sync.Mutex
}

type PacketStoreElement struct {
	counter       uint64
	sendTimestamp time.Time
}

func NewMessageStore(size int) *PacketStore {
	return &PacketStore{
		Buffer: ring.New(size),
		size:   size,
	}
}

func isMagicPacket(index uint64, packet2 gopacket.Packet) bool {
	ethLayer2 := packet2.Layer(layers.LayerTypeEthernet)
	if ethLayer2 == nil {
		return false
	}
	ethPacket2, _ := ethLayer2.(*layers.Ethernet)
	if ethPacket2.EthernetType != layers.EthernetType(0x1337) {
		return false
	}
	if len(ethPacket2.Payload) < binary.MaxVarintLen64 {
		return false
	}
	dataIndex := uint64(binary.LittleEndian.Uint64(ethPacket2.Payload))
	return index == dataIndex
}

func (mS *PacketStore) removePacket(p gopacket.Packet) (*measurement, error) {
	mS.mu.Lock()
	defer mS.mu.Unlock()

	for i := 0; i < mS.size; i++ {
		if mS.Buffer.Value != nil {
			e := mS.Buffer.Value.(PacketStoreElement)
			if isMagicPacket(e.counter, p) {
				meas := measurement{timestamp: e.sendTimestamp, reason: ReasonResolvedPacket, duration: p.Metadata().Timestamp.Sub(e.sendTimestamp)}
				mS.Buffer.Value = nil
				return &meas, nil
			}
		}
		mS.Buffer = mS.Buffer.Next()
	}
	return &measurement{timestamp: p.Metadata().Timestamp, reason: ReasonUnknownPacket}, nil
}

func (mS *PacketStore) addPacket(counter uint64, ts time.Time) (*measurement, error) {
	element := PacketStoreElement{
		counter:       counter,
		sendTimestamp: ts,
	}
	mS.mu.Lock()
	defer mS.mu.Unlock()

	if mS.Buffer.Value == nil {
		mS.Buffer.Value = element
		return nil, nil
	}

	oldestTimestamp := mS.Buffer.Value.(PacketStoreElement).sendTimestamp
	for i := 0; i < mS.size; i++ {
		if mS.Buffer.Value == nil {
			mS.Buffer.Value = element
			return nil, nil
		}
		ts := mS.Buffer.Value.(PacketStoreElement).sendTimestamp
		if ts.Before(oldestTimestamp) {
			oldestTimestamp = ts
		}
		mS.Buffer = mS.Buffer.Next()
	}

	for i := 0; i < mS.size; i++ {
		e := mS.Buffer.Value.(PacketStoreElement)
		if e.sendTimestamp.Equal(oldestTimestamp) {
			meas := measurement{timestamp: e.sendTimestamp, reason: ReasonAgedPacket}
			mS.Buffer.Value = element
			return &meas, nil
		}
		mS.Buffer = mS.Buffer.Next()
	}

	return nil, nil
}
