package om

import (
	"sort"
)

// generic parameters
type (
	keyT   = string // has to be comparable
	valueT = interface{}
	keysT  = index
)

// generic zeros
// var (
// 	zeroKey keyT
// )

// generic parameter (being a slice<T>)
type index []keyT

func (x index) Len() int           { return len(x) }
func (x index) Less(i, j int) bool { return x[i] < x[j] }
func (x index) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// check implementation
var _ sort.Interface = keysT{}

// Ordered map
type Ordered struct {
	_map   map[keyT]valueT
	_order keysT
}

// New .
func New(sourceMap map[keyT]valueT) *Ordered {
	if sourceMap == nil {
		sourceMap = make(map[keyT]valueT)
	}
	order := make(keysT, len(sourceMap))
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
	om._order = append(om._order, k) // TODO: find and insert? (benchmark)
	sort.Sort(om._order)
}

// ItrFn .
func (om *Ordered) ItrFn() func() (k keyT, v valueT, ok bool) {
	var lastIndex int
	return func() (k keyT, v valueT, ok bool) {
		if len(om._order) == 0 {
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
		lastIndex++
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
		return om._order[ix] >= k
	})
	if found == l {
		return
	}
	if om._order[found] != k {
		return
	}
	om._order = append(om._order[:found], om._order[found+1:]...)
	sort.Sort(om._order)
}
