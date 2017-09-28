package metrics

import (
	"fmt"
	"sync"
	"time"
)

//-----------------------------------------------------------------------------

// Average .
type Average struct {
	mx     sync.RWMutex
	count  int64
	av     float64
	format func(float64) string
}

// NewAverage .
func NewAverage(format func(float64) string) *Average {
	return &Average{format: format}
}

// Next .
func (av *Average) Next(v float64) {
	av.mx.Lock()
	defer av.mx.Unlock()

	newCount := av.count + 1

	av.av = (av.av*float64(av.count) + v) / float64(newCount)
	av.count = newCount
}

func (av *Average) String() (str string) {
	av.mx.RLock()
	defer av.mx.RUnlock()
	if av.format != nil {
		str = av.format(av.av)
	} else {
		str = fmt.Sprintf("%v", av.av)
	}
	return
}

//-----------------------------------------------------------------------------

// TimedAverage .
type TimedAverage struct {
	mx         sync.Mutex
	timeWindow time.Duration
	values     map[time.Time][]float64
	format     func(float64) string
}

// NewTimedAverage .
func NewTimedAverage(timeWindow time.Duration, format func(float64) string) *TimedAverage {
	if timeWindow <= 0 {
		timeWindow = time.Second * 90
	}
	res := &TimedAverage{
		timeWindow: timeWindow,
		values:     make(map[time.Time][]float64),
		format:     format,
	}
	return res
}

// Set .
func (opd *TimedAverage) Set(v float64) {
	opd.mx.Lock()
	defer opd.mx.Unlock()
	k := time.Now()
	opd.values[k] = append(opd.values[k], v)
}

func (opd *TimedAverage) String() (str string) {
	opd.mx.Lock()
	defer opd.mx.Unlock()

	if len(opd.values) == 0 {
		return "N/A"
	}
	var toDelete []time.Time
	count := 0
	var sum float64
	var min = time.Now().Add(time.Hour)

	for k, v := range opd.values {
		if time.Since(k) > opd.timeWindow {
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
		delete(opd.values, v)
	}

	N := float64(count)

	var av = sum / float64(count)
	if opd.format != nil {
		str = opd.format(av)
	} else {
		str = fmt.Sprintf("%v", av)
	}

	str = fmt.Sprintf("op/sec %.2f average ", N/time.Since(min).Seconds()) + str
	return
}

//-----------------------------------------------------------------------------
