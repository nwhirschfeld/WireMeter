package main

import (
	"container/ring"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"net"
	"sync"
	"time"
)

type PacketStore struct {
	Buffer *ring.Ring
	size   int
	mu     sync.Mutex
}

type PacketStoreElement struct {
	srcMac        net.HardwareAddr
	sendTimestamp time.Time
}

func NewMessageStore(size int) *PacketStore {
	return &PacketStore{
		Buffer: ring.New(size),
		size:   size,
	}
}

func hasPacketSrcMac(mac net.HardwareAddr, packet2 gopacket.Packet) bool {
	ethLayer2 := packet2.Layer(layers.LayerTypeEthernet)
	if ethLayer2 == nil {
		return false
	}
	ethPacket2, _ := ethLayer2.(*layers.Ethernet)
	return mac.String() == ethPacket2.SrcMAC.String()
}

func (mS *PacketStore) removePacket(p gopacket.Packet) (*measurement, error) {
	mS.mu.Lock()
	defer mS.mu.Unlock()

	for i := 0; i < mS.size; i++ {
		if mS.Buffer.Value != nil {
			e := mS.Buffer.Value.(PacketStoreElement)
			if hasPacketSrcMac(e.srcMac, p) {
				meas := measurement{timestamp: e.sendTimestamp, reason: ReasonResolvedPacket, duration: p.Metadata().Timestamp.Sub(e.sendTimestamp)}
				mS.Buffer.Value = nil
				return &meas, nil
			}
		}
		mS.Buffer = mS.Buffer.Next()
	}
	return &measurement{timestamp: p.Metadata().Timestamp, reason: ReasonUnknownPacket}, nil
}

func (mS *PacketStore) addPacket(p gopacket.Packet) (*measurement, error) {
	element := PacketStoreElement{
		srcMac:        p.Layer(layers.LayerTypeEthernet).(*layers.Ethernet).SrcMAC,
		sendTimestamp: p.Metadata().Timestamp,
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
