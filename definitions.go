package club

import (
	"fmt"
	"sync"
	"time"
)

//-----------------------------------------------------------------------------

// ByteSize represents size in bytes
type ByteSize float64

// ByteSize values
const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB

	// TB
	// PB
	// EB
	// ZB
	// YB
)

func (bs ByteSize) String() string {
	n := bs
	switch {
	case n > GB:
		return fmt.Sprintf("%.2f GB", n/GB)
	case n > MB:
		return fmt.Sprintf("%.2f MB", n/MB)
	case n > KB:
		return fmt.Sprintf("%.2f KB", n/KB)
	default:
		return fmt.Sprintf("%.2f B", n)
	}
}

//-----------------------------------------------------------------------------

// OpDuration for using with expvar
type OpDuration struct {
	mx         sync.Mutex
	historyAge time.Duration
	durations  map[time.Time][]time.Duration
}

// NewOpDuration .
func NewOpDuration(historyAge time.Duration) *OpDuration {
	if historyAge <= 0 {
		historyAge = time.Second * 90
	}
	res := &OpDuration{
		historyAge: historyAge,
		durations:  make(map[time.Time][]time.Duration),
	}
	return res
}

// Set .
func (opd *OpDuration) Set(d time.Duration) {
	opd.mx.Lock()
	defer opd.mx.Unlock()
	k := time.Now()
	opd.durations[k] = append(opd.durations[k], d)
}

func (opd *OpDuration) String() (str string) {
	opd.mx.Lock()
	defer opd.mx.Unlock()

	if len(opd.durations) == 0 {
		return "N/A"
	}
	var toDelete []time.Time
	count := 0
	var sum time.Duration
	var min = time.Now().Add(time.Hour)

	for k, v := range opd.durations {
		if time.Since(k) > opd.historyAge {
			toDelete = append(toDelete, k)
			continue
		}
		for _, vd := range v {
			count++
			sum += vd
		}
		if k.Before(min) {
			min = k
		}
	}
	for _, v := range toDelete {
		delete(opd.durations, v)
	}

	N := float64(count)

	return fmt.Sprintf("op/sec %.2f time/op %v", N/time.Since(min).Seconds(), sum/time.Duration(count))
}


//-----------------------------------------------------------------------------
