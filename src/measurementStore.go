package main

import (
	"bytes"
	"github.com/erkkah/margaid"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MeasurementStore struct {
	Measurements []measurement
	maxAge       time.Duration
	mu           sync.Mutex
}

func NewMeasureStore(maxAge time.Duration) MeasurementStore {
	return MeasurementStore{Measurements: []measurement{}, maxAge: maxAge}
}

func (mS *MeasurementStore) addMeasurement(m measurement) {
	mS.mu.Lock()
	defer mS.mu.Unlock()

	mS.Measurements = append(mS.Measurements, m)
	mS.deleteOldMeasurements()
}

func (mS *MeasurementStore) deleteOldMeasurements() {
	now := time.Now()
	for i, m := range mS.Measurements {
		if now.Sub(m.timestamp) < mS.maxAge {
			mS.Measurements = mS.Measurements[i:]
			break
		}
	}
}

func (mS MeasurementStore) analyzeMeasurements() ([]string, []int64, []int64, []int64, []int64, []int64, []int64) {
	timestamps := []time.Time{}
	averageDurations := []time.Duration{}
	minDurations := []time.Duration{}
	maxDurations := []time.Duration{}
	resolvedPackets := []int64{}
	unknownPackets := []int64{}
	agedPackets := []int64{}

	duration_cnt := 0
	var duration_sum time.Duration
	for _, m := range mS.Measurements {
		t := roundToSecond(m.timestamp)
		if !containsTimestamp(timestamps, t) {
			// timestamps
			timestamps = append(timestamps, t)
			// durations
			averageDurations = append(averageDurations, m.duration)
			minDurations = append(minDurations, m.duration)
			maxDurations = append(maxDurations, m.duration)
			duration_cnt = 1
			duration_sum = m.duration
			// types
			resolvedPackets = append(resolvedPackets, 0)
			unknownPackets = append(unknownPackets, 0)
			agedPackets = append(agedPackets, 0)
		} else {
			// durations
			duration_cnt += 1
			duration_sum = duration_sum + m.duration
			averageDurations[len(averageDurations)-1] = duration_sum / time.Duration(duration_cnt)
			if m.duration > maxDurations[len(maxDurations)-1] {
				maxDurations[len(maxDurations)-1] = m.duration
			}
			if m.duration < minDurations[len(minDurations)-1] {
				minDurations[len(minDurations)-1] = m.duration
			}
		}
		// types
		if m.reason == ReasonResolvedPacket {
			resolvedPackets[len(resolvedPackets)-1] += 1
		}
		if m.reason == ReasonUnknownPacket {
			unknownPackets[len(unknownPackets)-1] += 1
		}
		if m.reason == ReasonAgedPacket {
			agedPackets[len(agedPackets)-1] += 1
		}
	}

	return convertToTimeStrings(timestamps), durationsToMilliseconds(averageDurations), durationsToMilliseconds(minDurations), durationsToMilliseconds(maxDurations), resolvedPackets, unknownPackets, agedPackets
}

func (mS MeasurementStore) exportCSV() string {
	var builder strings.Builder
	header := []string{"Timestamp", "AvgDuration", "MinDuration", "MaxDuration", "ResolvedPackets", "UnknownPackets", "AgedPackets"}
	builder.WriteString(strings.Join(header, ",") + "\n")

	timestamps, avgDurations, minDurations, maxDurations, resolvedPackets, unknownPackets, agedPackets := mS.analyzeMeasurements()

	for i, t := range timestamps {
		record := []string{
			t,
			strconv.FormatInt(avgDurations[i], 10),
			strconv.FormatInt(minDurations[i], 10),
			strconv.FormatInt(maxDurations[i], 10),
			strconv.FormatInt(resolvedPackets[i], 10),
			strconv.FormatInt(unknownPackets[i], 10),
			strconv.FormatInt(agedPackets[i], 10),
		}
		builder.WriteString(strings.Join(record, ",") + "\n")
	}

	return builder.String()
}

func (mS MeasurementStore) exportSVG() (string, string) {
	timestamps, avgDurations, minDurations, maxDurations, resolvedPackets, unknownPackets, agedPackets := mS.analyzeMeasurements()
	//length := 0.01 * float64(len(timestamps)-1)
	//yDurationMod := float64(100) / float64(slices.Max(maxDurations))

	//timestampSeries := margaid.NewSeries()
	avgDurationSeries := margaid.NewSeries(margaid.Titled("average"))
	minDurationSeries := margaid.NewSeries(margaid.Titled("min"))
	maxDurationSeries := margaid.NewSeries(margaid.Titled("max"))
	resPacketSeries := margaid.NewSeries(margaid.Titled("resolved"))
	unkPacketSeries := margaid.NewSeries(margaid.Titled("unknown"))
	agePacketSeries := margaid.NewSeries(margaid.Titled("aged"))
	for i, _ := range timestamps {
		avgDurationSeries.Add(margaid.MakeValue(float64(i), float64(avgDurations[i])))
		minDurationSeries.Add(margaid.MakeValue(float64(i), float64(minDurations[i])))
		maxDurationSeries.Add(margaid.MakeValue(float64(i), float64(maxDurations[i])))
		resPacketSeries.Add(margaid.MakeValue(float64(i), float64(resolvedPackets[i])))
		unkPacketSeries.Add(margaid.MakeValue(float64(i), float64(unknownPackets[i])))
		agePacketSeries.Add(margaid.MakeValue(float64(i), float64(agedPackets[i])))
	}

	// Create the runtimeDiagram object:
	runtimeDiagram := margaid.New(1200, 600,
		margaid.WithBackgroundColor("white"),
		margaid.WithAutorange(margaid.XAxis, maxDurationSeries),
		margaid.WithRange(margaid.YAxis, 0, float64(1+slices.Max(maxDurations))),
		margaid.WithColorScheme(90),
	)
	lossMax := slices.Max([]int64{slices.Max(resolvedPackets), slices.Max(unknownPackets), slices.Max(agedPackets)})
	lossDiagram := margaid.New(1200, 600,
		margaid.WithBackgroundColor("white"),
		margaid.WithAutorange(margaid.XAxis, maxDurationSeries),
		margaid.WithRange(margaid.YAxis, 0, 1.1*float64(lossMax)),
		margaid.WithColorScheme(90),
	)

	// Plot the series
	runtimeDiagram.Smooth(avgDurationSeries, margaid.UsingAxes(margaid.XAxis, margaid.YAxis), margaid.UsingStrokeWidth(3))
	runtimeDiagram.Smooth(minDurationSeries, margaid.UsingAxes(margaid.XAxis, margaid.YAxis), margaid.UsingStrokeWidth(3))
	runtimeDiagram.Smooth(maxDurationSeries, margaid.UsingAxes(margaid.XAxis, margaid.YAxis), margaid.UsingStrokeWidth(3))
	lossDiagram.Smooth(resPacketSeries, margaid.UsingAxes(margaid.XAxis, margaid.YAxis), margaid.UsingStrokeWidth(3))
	lossDiagram.Smooth(unkPacketSeries, margaid.UsingAxes(margaid.XAxis, margaid.YAxis), margaid.UsingStrokeWidth(3))
	lossDiagram.Smooth(agePacketSeries, margaid.UsingAxes(margaid.XAxis, margaid.YAxis), margaid.UsingStrokeWidth(3))

	// Add a frame and X axis
	runtimeDiagram.Frame()
	runtimeDiagram.Legend(margaid.RightBottom)
	runtimeDiagram.Title("Packet Runtime")
	runtimeDiagram.Axis(avgDurationSeries, margaid.XAxis, runtimeDiagram.ValueTicker('f', 0, 30), false, "s")
	runtimeDiagram.Axis(avgDurationSeries, margaid.YAxis, runtimeDiagram.ValueTicker('f', 0, 5), false, "ms")
	lossDiagram.Frame()
	lossDiagram.Legend(margaid.RightBottom)
	lossDiagram.Title("Measurement Type")
	lossDiagram.Axis(avgDurationSeries, margaid.XAxis, runtimeDiagram.ValueTicker('f', 0, 30), false, "s")
	lossDiagram.Axis(avgDurationSeries, margaid.YAxis, runtimeDiagram.ValueTicker('f', 0, 25), false, "count")
	// Render to stdout
	runtimeBuffer := new(bytes.Buffer)
	runtimeDiagram.Render(runtimeBuffer)
	lossBuffer := new(bytes.Buffer)
	lossDiagram.Render(lossBuffer)
	return runtimeBuffer.String(), lossBuffer.String()
}

func durationsToMilliseconds(durations []time.Duration) []int64 {
	var result []int64

	for _, duration := range durations {
		millis := duration.Milliseconds()
		result = append(result, millis)
	}

	return result
}

func roundToSecond(t time.Time) time.Time {
	return time.Unix(t.Unix(), 0)
}

func containsTimestamp(timestamps []time.Time, timestamp time.Time) bool {
	if len(timestamps) == 0 {
		return false
	}
	return timestamps[len(timestamps)-1] == timestamp
}

func convertToTimeStrings(timestamps []time.Time) []string {
	var timeStrings []string

	for _, t := range timestamps {
		// Use Format("15:04:05") to represent hours, minutes, and seconds
		timeString := t.Format("15:04:05")
		timeStrings = append(timeStrings, timeString)
	}

	return timeStrings
}
