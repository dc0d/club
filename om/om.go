package om

import (
	"sort"
)

// generic parameters
type (
	keyT   = string
	valueT = interface{}
)

func compareKey(k1, k2 keyT) int {
	if k1 < k2 {
		return -1
	}
	if k1 == k2 {
		return 0
	}
	return 1
}

// Ordered map
type Ordered struct {
	_map   map[keyT]valueT
	_order index
}

// New .
func New(sourceMap map[keyT]valueT) *Ordered {
	if sourceMap == nil {
		sourceMap = make(map[keyT]valueT)
	}
	order := make(index, len(sourceMap))
	c := 0
	for k := range sourceMap {
		order[c] = k
		c++
	}
	sort.Sort(order)
	res := &Ordered{
		_map:   sourceMap,
		_order: order,
	}
	return res
}

// Get .
func (om *Ordered) Get(k keyT) (v valueT, ok bool) {
	v, ok = om._map[k]
	return
}

// Put .
func (om *Ordered) Put(k keyT, v valueT) {
	om._map[k] = v
	l := len(om._order)
	found := sort.Search(l, func(ix int) bool {
		cm := compareKey(om._order[ix], k)
		return cm >= 0
	})
	if found == l {
		om._order = append(om._order, k)
		return
	}
	if om._order[found] == k {
		return
	}
	om._order = append(om._order, k)
	copy(om._order[found+1:], om._order[found:])
	om._order[found] = k
}

// ItrFn .
func (om *Ordered) ItrFn(desc ...bool) func() (k keyT, v valueT, ok bool) {
	dec := false
	if len(desc) > 0 {
		dec = desc[0]
	}
	var lastIndex int
	if dec {
		lastIndex = len(om._order) - 1
	}
	return func() (k keyT, v valueT, ok bool) {
		if len(om._order) == 0 {
			ok = false
			return
		}
		if dec && lastIndex < 0 {
			ok = false
			return
		}
		if lastIndex >= len(om._order) {
			ok = false
			return
		}
		k = om._order[lastIndex]
		v = om._map[k]
		ok = true
		if dec {
			lastIndex--
		} else {
			lastIndex++
		}
		return
	}
}

// ItrSeek moves to first key equal or greater than seek
func (om *Ordered) ItrSeek(seek keyT, desc ...bool) func() (k keyT, v valueT, ok bool) {
	dec := false
	if len(desc) > 0 {
		dec = desc[0]
	}
	l := len(om._order)
	found := sort.Search(l, func(ix int) bool {
		cm := compareKey(om._order[ix], seek)
		return cm >= 0
	})
	var lastIndex int
	if dec {
		lastIndex = len(om._order) - 1
	}
	if found != l {
		lastIndex = found
	}
	return func() (k keyT, v valueT, ok bool) {
		if len(om._order) == 0 {
			ok = false
			return
		}
		if dec && lastIndex < 0 {
			ok = false
			return
		}
		if lastIndex >= len(om._order) {
			ok = false
			return
		}
		k = om._order[lastIndex]
		v = om._map[k]
		ok = true
		if dec {
			lastIndex--
		} else {
			lastIndex++
		}
		return
	}
}

// Del .
func (om *Ordered) Del(k keyT) {
	delete(om._map, k)
	l := len(om._order)
	if l == 0 {
		return
	}
	found := sort.Search(l, func(ix int) bool {
		cm := compareKey(om._order[ix], k)
		return cm >= 0
	})
	if found == l {
		return
	}
	if om._order[found] != k {
		return
	}
	om._order = append(om._order[:found], om._order[found+1:]...)
}

type index []keyT

func (x index) Len() int           { return len(x) }
func (x index) Less(i, j int) bool { return compareKey(x[i], x[j]) == -1 }
func (x index) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
