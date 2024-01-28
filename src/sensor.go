package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"sync"
	"time"
)

type WireSensor struct {
	receivingInterface string
	sendingInterface   string
	waitTime           int
	store              *PacketStore
	measurements       MeasurementStore
	wg                 sync.WaitGroup
}

func NewWireSensor(receivingInterface, sendingInterface string, storeSize int, waitTime int) *WireSensor {
	return &WireSensor{
		receivingInterface: receivingInterface,
		sendingInterface:   sendingInterface,
		waitTime:           waitTime,
		measurements:       NewMeasureStore(3 * time.Minute),
		store:              NewMessageStore(storeSize),
	}
}

func (ws *WireSensor) run() {
	ws.runPacketSink()
	ws.runPacketSource()
}

func (ws *WireSensor) runPacketSource() {
	go func() {
		ws.PacketSource()
	}()
}

func (ws *WireSensor) PacketSource() {
	ws.wg.Add(1)
	sendingPcapHandle, err := pcap.OpenLive(ws.sendingInterface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Duration(ws.waitTime) * time.Millisecond)

	go func() {
		i := 0
		for {
			<-ticker.C // Wait for the ticker to tick
			p := getPacketN(i)
			p.Metadata().Timestamp = time.Now()
			err = sendingPcapHandle.WritePacketData(p.Data())
			if err != nil {
				log.Fatal(err)
			}

			rp, err := ws.store.addPacket(p)
			if err != nil {
				return
			}
			if rp != nil {
				ws.measurements.addMeasurement(*rp)
			}
			i++
		}
	}()

	// Keep the program running
	select {}
	sendingPcapHandle.Close()
	ws.wg.Done()
}

func (ws *WireSensor) runPacketSink() {
	go func() {
		ws.PacketSink()
	}()
}

func (ws *WireSensor) PacketSink() {
	ws.wg.Add(1)
	receivingPcapHandle, err := pcap.OpenLive(ws.receivingInterface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer receivingPcapHandle.Close()

	packetSource1 := gopacket.NewPacketSource(receivingPcapHandle, receivingPcapHandle.LinkType())

	for packet := range packetSource1.Packets() {
		packet.Metadata().Timestamp = time.Now()
		rp, err := ws.store.removePacket(packet)
		if err != nil {
			return
		}
		if rp != nil {
			ws.measurements.addMeasurement(*rp)
		}
	}
	ws.wg.Done()
}
