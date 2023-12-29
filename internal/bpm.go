package internal

type BPMPart struct {
	T0       float64
	T1       float64
	BPM      float64
	Position int
}

type ChartInfo struct {
	Name             string `json:"name"`
	Artist           string `json:"artist"`
	AverageBeatCount int    `json:"average_beat_count"`
	BeatsPerBar      int    `json:"beats_per_bar"`
}

func CalculateBPMParts(tapTimes []float64, chartInfo ChartInfo) []BPMPart {
	bpmParts := make([]BPMPart, 0)
	i := 0
	r := 0
	n := chartInfo.AverageBeatCount
	factor := chartInfo.AverageBeatCount / chartInfo.BeatsPerBar
	for i < len(tapTimes)-n {
		currentTapTimes := tapTimes[i:(i + n + 1)]

		starts := currentTapTimes[0 : len(currentTapTimes)-1]
		ends := currentTapTimes[1:]

		dtSum := 0.
		for j := 0; j < n; j++ {
			dtSum += ends[j] - starts[j]
		}
		avg := dtSum / float64(n)
		bpm := 60. / avg
		bpmParts = append(bpmParts, BPMPart{
			T0:       starts[0],
			T1:       ends[len(ends)-1],
			BPM:      bpm,
			Position: 192 * r * chartInfo.BeatsPerBar * factor,
		})

		i += n
		r += 1
	}
	return bpmParts
}
