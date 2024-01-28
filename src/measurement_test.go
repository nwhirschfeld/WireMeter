package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMeasurementReasonString(t *testing.T) {
	assert.Equal(t, "unknown packet", MeasurementReason(ReasonUnknownPacket).String())
	assert.Equal(t, "packet dropped (age)", MeasurementReason(ReasonAgedPacket).String())
	assert.Equal(t, "packet resolved", MeasurementReason(ReasonResolvedPacket).String())
}
