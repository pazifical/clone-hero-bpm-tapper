package internal

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	tapTimes := []float64{
		1, 2, 3, 4, 6, 8, 10, 12,
	}
	chartInfo := ChartInfo{Name: "Name", Artist: "Artist", AverageBeatCount: 4, Factor: 1}

	bpmParts := CalculateBPMParts(tapTimes, chartInfo)
	fmt.Println(bpmParts)

	if bpmParts[0].BPM != 60. {
		t.Errorf("want: %f != got: %f", 60., bpmParts[0].BPM)
	}
	if bpmParts[1].BPM != 30. {
		t.Errorf("want: %f != got: %f", 30., bpmParts[1].BPM)
	}
}
