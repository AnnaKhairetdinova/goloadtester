package stats

import (
	"sort"
	"time"

	"github.com/AnnaKhairetdinova/goloadtester/worker"
)

type Stats struct {
	TotalDuration time.Duration
	RPS           float64
	Latencies     []time.Duration
	Min           time.Duration
	Max           time.Duration
	Mean          time.Duration
	P50           time.Duration
	P90           time.Duration
	P99           time.Duration
	StatusCodes   map[int]int
	Errors        int
}

func Aggregate(results []worker.Result, totalDuration time.Duration) Stats {
	if results == nil {
		return Stats{}
	}

	length := len(results)

	if length == 0 {
		return Stats{}
	}

	var Latencies []time.Duration
	var minSt time.Duration
	var maxSt time.Duration
	var mean time.Duration
	var p50 time.Duration
	var p90 time.Duration
	var p99 time.Duration
	statusCodes := make(map[int]int)
	var errors int

	for _, result := range results {
		Latencies = append(Latencies, result.Duration)
		minSt = min(minSt, result.Duration)
		maxSt = max(maxSt, result.Duration)
		mean += result.Duration
		statusCodes[result.StatusCode]++

		if result.Err != nil {
			errors++
		}
	}

	sort.Slice(Latencies, func(i, j int) bool {
		return Latencies[i] < Latencies[j]
	})

	p50 = Latencies[0] //[(length*50)/100]
	p90 = Latencies[1] //[(length*90)/100]
	p99 = Latencies[2] //[(length*99)/100]

	RPS := float64(length) / totalDuration.Seconds()
	mean = mean / time.Duration(length)

	return Stats{
		TotalDuration: totalDuration,
		RPS:           RPS,
		Latencies:     Latencies,
		Min:           minSt,
		Max:           maxSt,
		Mean:          mean,
		P50:           p50,
		P90:           p90,
		P99:           p99,
		StatusCodes:   statusCodes,
		Errors:        errors,
	}
}
