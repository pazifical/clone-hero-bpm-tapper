package internal

import (
	"fmt"
	"testing"
)

func Test44WithAverage4(t *testing.T) {
	tapTimes := []float64{
		0.0, 0.5, 1.0, 1.5, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0,
	}
	chartInfo := ChartInfo{Name: "Name", Artist: "Artist", AverageBeatCount: 4, BeatsPerBar: 4}

	bpmParts := CalculateBPMParts(tapTimes, chartInfo)
	fmt.Println(bpmParts)

	want := 120.
	got := bpmParts[0].BPM
	if want != got {
		t.Errorf("want: %f != got: %f", want, got)
	}
	want = 60.
	got = bpmParts[1].BPM
	if want != got {
		t.Errorf("want: %f != got: %f", want, got)
	}
	wantPos := 0
	gotPos := bpmParts[0].Position
	if gotPos != wantPos {
		t.Errorf("want: %d != got: %d", wantPos, gotPos)
	}
	wantPos = 192 * 4
	gotPos = bpmParts[1].Position
	if gotPos != wantPos {
		t.Errorf("want: %d != got: %d", wantPos, gotPos)
	}
}
func Test44WithAverage8(t *testing.T) {
	tapTimes := []float64{
		0.0, 0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 5., 6., 7., 8., 9., 10., 11., 12.,
	}
	chartInfo := ChartInfo{Name: "Name", Artist: "Artist", AverageBeatCount: 8, BeatsPerBar: 4}

	bpmParts := CalculateBPMParts(tapTimes, chartInfo)
	fmt.Println(bpmParts)

	want := 120.
	got := bpmParts[0].BPM
	if want != got {
		t.Errorf("want: %f != got: %f", want, got)
	}
	want = 60.
	got = bpmParts[1].BPM
	if want != got {
		t.Errorf("want: %f != got: %f", want, got)
	}
	wantPos := 0
	gotPos := bpmParts[0].Position
	if gotPos != wantPos {
		t.Errorf("want: %d != got: %d", wantPos, gotPos)
	}
	wantPos = 192 * 8
	gotPos = bpmParts[1].Position
	if gotPos != wantPos {
		t.Errorf("want: %d != got: %d", wantPos, gotPos)
	}
}
