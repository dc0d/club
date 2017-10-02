package metrics

import (
	"fmt"
	"sync"
	"time"
)

//-----------------------------------------------------------------------------

type options struct {
	maxAge time.Duration
	format func(float64) string
}

func newOptions() *options {
	return &options{
		maxAge: time.Second * 120,
	}
}

// Option .
type Option func(*options)

// MaxAge option
func MaxAge(maxAge time.Duration) Option {
	return func(opt *options) {
		opt.maxAge = maxAge
	}
}

// Format option
func Format(format func(float64) string) Option {
	return func(opt *options) {
		opt.format = format
	}
}

//-----------------------------------------------------------------------------

// Average .
type Average struct {
	mx      sync.RWMutex
	count   int64
	av      float64
	options *options
}

// NewAverage .
func NewAverage(options ...Option) *Average {
	res := &Average{options: newOptions()}
	for _, vf := range options {
		vf(res.options)
	}
	return res
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
	if av.options.format != nil {
		str = av.options.format(av.av)
	} else {
		str = fmt.Sprintf("%v", av.av)
	}
	return
}

//-----------------------------------------------------------------------------

// TimedAverage .
type TimedAverage struct {
	mx      sync.Mutex
	values  map[time.Time][]float64
	options *options
}

// NewTimedAverage .
func NewTimedAverage(options ...Option) *TimedAverage {
	res := &TimedAverage{
		values:  make(map[time.Time][]float64),
		options: newOptions(),
	}
	for _, vf := range options {
		vf(res.options)
	}
	if res.options.maxAge <= 0 {
		res.options.maxAge = 120
	}
	return res
}

// Set .
func (opd *TimedAverage) Set(v float64) {
	opd.mx.Lock()
	defer opd.mx.Unlock()
	k := time.Now()
	opd.values[k] = append(opd.values[k], v)

	var toDelete []time.Time
	for k := range opd.values {
		if time.Since(k) > opd.options.maxAge {
			toDelete = append(toDelete, k)
		}
	}
	for _, v := range toDelete {
		delete(opd.values, v)
	}
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
		if time.Since(k) > opd.options.maxAge {
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
	if opd.options.format != nil {
		str = opd.options.format(av)
	} else {
		str = fmt.Sprintf("%v", av)
	}

	str = fmt.Sprintf("op/sec %.2f average ", N/time.Since(min).Seconds()) + str
	return
}

//-----------------------------------------------------------------------------
