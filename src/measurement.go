package main

import (
	"time"
)

type measurement struct {
	timestamp time.Time
	reason    MeasurementReason
	duration  time.Duration
}

type MeasurementReason int

const (
	ReasonUnknownPacket = iota
	ReasonAgedPacket
	ReasonResolvedPacket
)

func (mr MeasurementReason) String() string {
	return []string{
		"unknown packet",
		"packet dropped (age)",
		"packet resolved",
	}[mr]
}
