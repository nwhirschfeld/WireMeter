package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createEmptyPacket() gopacket.Packet {
	buffer := gopacket.NewSerializeBuffer()
	ethLayer := &layers.Ethernet{
		SrcMAC:       nil,                     // replace with actual source MAC if needed
		DstMAC:       nil,                     // replace with actual destination MAC if needed
		EthernetType: layers.EthernetTypeIPv4, // or the appropriate EthernetType for your use case
	}
	gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{},
		ethLayer,
	)
	packet := gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeEthernet, gopacket.Default)
	return packet
}

func TestAddMeasurement(t *testing.T) {
	maxAge := 10 * time.Second
	ms := NewMeasureStore(maxAge)

	measurement1 := measurement{
		timestamp: time.Date(1970, 01, 01, 23, 42, 23, 1337, time.UTC),
		reason:    ReasonResolvedPacket,
		duration:  5 * time.Second,
	}

	measurement2 := measurement{
		timestamp: time.Date(1970, 01, 01, 23, 42, 23, 1337, time.UTC),
		reason:    ReasonUnknownPacket,
		duration:  3 * time.Second,
	}

	ms.addMeasurement(measurement1)
	ms.addMeasurement(measurement2)

	assert.Equal(t, 2, len(ms.Measurements))
}

func TestDeleteOldMeasurements(t *testing.T) {
	maxAge := 10 * time.Second
	ms := NewMeasureStore(maxAge)

	measurement1 := measurement{
		timestamp: time.Now().Add(-1 * maxAge),
		reason:    ReasonResolvedPacket,
		duration:  5 * time.Second,
	}

	measurement2 := measurement{
		timestamp: time.Now(),
		reason:    ReasonUnknownPacket,
		duration:  3 * time.Second,
	}

	ms.Measurements = append(ms.Measurements, measurement1, measurement2)

	ms.deleteOldMeasurements()

	assert.Equal(t, 1, len(ms.Measurements))
	assert.Equal(t, measurement2.timestamp, ms.Measurements[0].timestamp)
}

func TestAnalyzeMeasurements(t *testing.T) {
	maxAge := 10 * time.Second
	ms := NewMeasureStore(maxAge)

	measurement1 := measurement{
		timestamp: time.Date(1970, 01, 01, 23, 42, 23, 1337, time.UTC),
		reason:    ReasonResolvedPacket,
		duration:  5 * time.Millisecond,
	}

	measurement2 := measurement{
		timestamp: time.Date(1970, 01, 01, 23, 42, 23, 1337, time.UTC),
		reason:    ReasonUnknownPacket,
		duration:  3 * time.Millisecond,
	}

	ms.Measurements = append(ms.Measurements, measurement1, measurement2)

	timeStrings, avgDurations, minDurations, maxDurations, resolvedPackets, unknownPackets, agedPackets := ms.analyzeMeasurements()

	assert.Equal(t, []string{"00:42:23"}, timeStrings)
	assert.Equal(t, []int64{4}, avgDurations) // Average of 5 and 3
	assert.Equal(t, []int64{3}, minDurations) // Minimum of 5 and 3
	assert.Equal(t, []int64{5}, maxDurations) // Maximum of 5 and 3
	assert.Equal(t, []int64{1}, resolvedPackets)
	assert.Equal(t, []int64{1}, unknownPackets)
	assert.Equal(t, []int64{0}, agedPackets)
}
