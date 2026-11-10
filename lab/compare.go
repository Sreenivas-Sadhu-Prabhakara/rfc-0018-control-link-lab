package lab

import "fmt"

const (
	Host = "host"
	IMP  = "imp"
)

type Scenario struct {
	Arrivals []Arrival `json:"arrivals"`
}

type Arrival struct {
	Tick  int    `json:"tick"`
	Class string `json:"class"`
}

type Result struct {
	Arrival
	Message    int `json:"message"`
	SharedWait int `json:"shared_wait"`
	SplitWait  int `json:"split_wait"`
}

type Summary struct {
	Class     string `json:"class"`
	Count     int    `json:"count"`
	SharedMax int    `json:"shared_max_wait"`
	SplitMax  int    `json:"split_max_wait"`
}

type Report struct {
	Classes []Summary `json:"classes"`
	Results []Result  `json:"results"`
}

// lane tracks the next service slot in one FIFO queue.
type lane struct {
	tick int
	used int
}

func (l *lane) wait(arrival, slots int) int {
	if l.tick < arrival {
		l.tick = arrival
		l.used = 0
	}
	wait := l.tick - arrival
	l.used++
	if l.used == slots {
		l.tick++
		l.used = 0
	}
	return wait
}

// Compare gives both policies two service slots per tick. Shared uses one
// FIFO with two slots. Split uses one FIFO and one slot for each class.
func Compare(s Scenario) (Report, error) {
	if len(s.Arrivals) == 0 || len(s.Arrivals) > 500 {
		return Report{}, fmt.Errorf("expected 1 to 500 arrivals")
	}

	report := Report{
		Classes: []Summary{{Class: Host}, {Class: IMP}},
		Results: make([]Result, 0, len(s.Arrivals)),
	}
	var shared lane
	var reserved [2]lane
	lastTick := 0
	for i, arrival := range s.Arrivals {
		if arrival.Tick < lastTick || arrival.Tick > 10000 {
			return Report{}, fmt.Errorf("arrival %d: ticks must be nondecreasing and within 0 to 10000", i+1)
		}
		lastTick = arrival.Tick
		var classIndex int
		switch arrival.Class {
		case Host:
			classIndex = 0
		case IMP:
			classIndex = 1
		default:
			return Report{}, fmt.Errorf("arrival %d: class must be host or imp", i+1)
		}
		result := Result{
			Arrival:    arrival,
			Message:    i + 1,
			SharedWait: shared.wait(arrival.Tick, 2),
			SplitWait:  reserved[classIndex].wait(arrival.Tick, 1),
		}
		summary := &report.Classes[classIndex]
		summary.Count++
		summary.SharedMax = max(summary.SharedMax, result.SharedWait)
		summary.SplitMax = max(summary.SplitMax, result.SplitWait)
		report.Results = append(report.Results, result)
	}
	return report, nil
}
