package internal

import (
	"fmt"
	"testing"
)

func Test44WithAverage4(t *testing.T) {
	tapTimes := []float64{
		1, 2, 3, 4, 6, 8, 10, 12,
	}
	chartInfo := ChartInfo{Name: "Name", Artist: "Artist", AverageBeatCount: 4, BeatsPerBar: 4}

	bpmParts := CalculateBPMParts(tapTimes, chartInfo)
	fmt.Println(bpmParts)

	if bpmParts[0].BPM != 60. {
		t.Errorf("want: %f != got: %f", 60., bpmParts[0].BPM)
	}
	if bpmParts[1].BPM != 30. {
		t.Errorf("want: %f != got: %f", 30., bpmParts[1].BPM)
	}
	if float64(bpmParts[0].Position) != 0 {
		t.Errorf("want: %d != got: %d", 0, bpmParts[0].Position)
	}
	if float64(bpmParts[1].Position) != 192*4 {
		t.Errorf("want: %d != got: %d", 192*4, bpmParts[1].Position)
	}
}
func Test44WithAverage8(t *testing.T) {
	tapTimes := []float64{
		1, 2, 3, 4, 5, 6, 7, 8, 10, 12, 14, 16, 18, 20, 22,
	}
	chartInfo := ChartInfo{Name: "Name", Artist: "Artist", AverageBeatCount: 8, BeatsPerBar: 4}

	bpmParts := CalculateBPMParts(tapTimes, chartInfo)
	fmt.Println(bpmParts)

	if bpmParts[0].BPM != 60. {
		t.Errorf("want: %f != got: %f", 60., bpmParts[0].BPM)
	}
	if bpmParts[1].BPM != 30. {
		t.Errorf("want: %f != got: %f", 30., bpmParts[1].BPM)
	}

	if float64(bpmParts[0].Position) != 0 {
		t.Errorf("want: %d != got: %d", 0, bpmParts[0].Position)
	}
	if float64(bpmParts[1].Position) != 192*8 {
		t.Errorf("want: %d != got: %d", 192*8, bpmParts[1].Position)
	}
}
