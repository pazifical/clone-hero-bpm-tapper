package internal

import "fmt"

type BPMPart struct {
	T0       float64
	T1       float64
	BPM      float64
	Position int
}

type ChartInfo struct {
	Name             string  `json:"name"`
	Artist           string  `json:"artist"`
	AverageBeatCount int     `json:"average_beat_count"`
	Factor           float64 `json:"factor"`
}

func CalculateBPMParts(tapTimes []float64, chartInfo ChartInfo) []BPMPart {
	bpmParts := make([]BPMPart, 0)
	i := 0
	n := chartInfo.AverageBeatCount
	for i < len(tapTimes)-n {
		currentTapTimes := tapTimes[i:(i + n)]
		fmt.Println(currentTapTimes)

		starts := currentTapTimes[0 : len(currentTapTimes)-1]
		ends := currentTapTimes[1:]

		dtSum := 0.
		for j := 0; j < len(starts); j++ {
			dtSum += ends[j] - starts[j]
		}
		avg := dtSum / float64(len(starts))
		bpm := 60. / avg
		bpmParts = append(bpmParts, BPMPart{
			T0:       tapTimes[i],
			T1:       tapTimes[i+n-1],
			BPM:      bpm,
			Position: 192 * (i + 1),
		})

		i += (n - 1)
	}
	return bpmParts
}