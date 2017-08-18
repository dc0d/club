package clubaux

import (
	"context"
	"sync"
	"time"

	"github.com/dc0d/club"
)

//-----------------------------------------------------------------------------

// Registry currently leaks one goroutine
type Registry struct {
	ctx                context.Context
	expirationInterval time.Duration

	rwx          sync.RWMutex
	values       map[interface{}]interface{}
	expiresAt    map[interface{}]time.Time
	expiresAfter map[interface{}]time.Duration
	isSliding    map[interface{}]struct{}
}

// NewRegistry min of expirationInterval is 1 ms, default (if passed 0) is 30 sec, ctx can be nil
func NewRegistry(ctx context.Context, expirationInterval time.Duration) *Registry {
	if expirationInterval < 0 {
		expirationInterval = time.Millisecond
	}
	if expirationInterval == 0 {
		expirationInterval = time.Second * 30
	}
	res := &Registry{
		expirationInterval: expirationInterval,
		values:             make(map[interface{}]interface{}),
		expiresAt:          make(map[interface{}]time.Time),
		expiresAfter:       make(map[interface{}]time.Duration),
		isSliding:          make(map[interface{}]struct{}),
	}
	if ctx != nil {
		res.ctx = ctx
	}
	go res._expireLoop()
	return res
}

// Get .
func (rg *Registry) Get(k interface{}) (interface{}, bool) {
	slide := false
	rg.rwx.RLock()
	v, ok := rg.values[k]
	if ok {
		_, slide = rg.isSliding[k]
	}
	rg.rwx.RUnlock()
	if slide {
		rg.rwx.Lock()
		rg.expiresAt[k] = time.Now().Add(rg.expiresAfter[k])
		rg.rwx.Unlock()
	}
	return v, ok
}

// Put .
func (rg *Registry) Put(k, v interface{}) {
	rg.rwx.Lock()
	rg.values[k] = v
	rg.rwx.Unlock()
}

// PutWithExpiration .
func (rg *Registry) PutWithExpiration(k, v interface{},
	expiresAfter time.Duration,
	isSliding ...bool) {
	rg.rwx.Lock()
	rg.values[k] = v
	rg.expiresAfter[k] = expiresAfter
	rg.expiresAt[k] = time.Now().Add(expiresAfter)
	if len(isSliding) > 0 && isSliding[0] {
		rg.isSliding[k] = struct{}{}
	}
	rg.rwx.Unlock()
}

// errors
var (
	ErrNotFound = club.Errorf("NOT_FOUND")
	ErrCASCond  = club.Errorf("CAS_COND_FAILED")
)

// CAS .
func (rg *Registry) CAS(k, v interface{}, cond func(interface{}) bool) error {
	rg.rwx.Lock()
	defer rg.rwx.Unlock()
	old, ok := rg.values[k]
	if !ok {
		return ErrNotFound
	}
	if !cond(old) {
		return ErrCASCond
	}
	rg.values[k] = v
	return nil
}

// Delete .
func (rg *Registry) Delete(k interface{}) {
	rg.rwx.Lock()
	rg._delete(k)
	rg.rwx.Unlock()
}

//-----------------------------------------------------------------------------

func (rg *Registry) _expireLoop() {
	for {
		var done <-chan struct{}
		if rg.ctx != nil {
			done = rg.ctx.Done() // TODO: should I get this every time?
		}
		select {
		case <-done:
			return
		case <-time.After(rg.expirationInterval):
			rg._expireFunc()
		}
	}
}

func (rg *Registry) _expireFunc() {
	rg.rwx.Lock()
	for k, v := range rg.expiresAt {
		if !time.Now().After(v) {
			continue
		}
		rg._delete(k)
	}
	rg.rwx.Unlock()
}

func (rg *Registry) _delete(k interface{}) {
	delete(rg.expiresAt, k)
	delete(rg.expiresAfter, k)
	delete(rg.isSliding, k)
	delete(rg.values, k)
}

//-----------------------------------------------------------------------------