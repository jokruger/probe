package probe

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

type ProbeResult struct {
	TotalTime time.Duration
	CallCount int
}

var (
	reportLock sync.Mutex
	results    = make(map[string]ProbeResult)
)

type Probe struct {
	name  string
	start time.Time
}

func Start(name string) Probe {
	return Probe{name: name, start: time.Now()}
}

func (p Probe) Stop() {
	duration := time.Since(p.start)

	reportLock.Lock()
	defer reportLock.Unlock()

	r := results[p.name]
	r.TotalTime += duration
	r.CallCount++
	results[p.name] = r
}

func PrintReport() {
	reportLock.Lock()
	defer reportLock.Unlock()

	if len(results) == 0 {
		fmt.Println("\nNo probes to report.")
		return
	}

	type result struct {
		name        string
		callCount   int
		totalTime   time.Duration
		avgTime     time.Duration
		callsPerSec float64
	}

	resultsList := make([]result, 0, len(results))
	for name, res := range results {
		r := result{
			name:      name,
			callCount: res.CallCount,
			totalTime: res.TotalTime,
		}

		if res.CallCount > 0 {
			r.avgTime = res.TotalTime / time.Duration(res.CallCount)
			if res.TotalTime.Seconds() > 0 {
				r.callsPerSec = float64(res.CallCount) / res.TotalTime.Seconds()
			}
		}

		resultsList = append(resultsList, r)
	}
	slices.SortFunc(resultsList, func(a, b result) int {
		return int(b.totalTime.Milliseconds() - a.totalTime.Milliseconds())
	})

	fmt.Println("\nAggregated Execution Time Report:")
	fmt.Println("--------------------------------------------------------------------------")
	fmt.Printf("%-20s %10s %15s %15s %10s\n", "Name", "Calls", "Total Time", "Avg Time", "Calls/sec")
	fmt.Println("--------------------------------------------------------------------------")
	for _, r := range resultsList {
		fmt.Printf("%-20s %10d %15v %15v %10.2f\n", r.name, r.callCount, r.totalTime, r.avgTime, r.callsPerSec)
	}
	fmt.Println("--------------------------------------------------------------------------")
}
