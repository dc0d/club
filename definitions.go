package club

import (
	"bytes"
	"strings"
	"sync"
)

//-----------------------------------------------------------------------------

// Error constant error
type Error string

func (v Error) Error() string { return string(v) }

//-----------------------------------------------------------------------------

// Errors .
type Errors []error

func (x Errors) Error() string {
	if x == nil {
		return ``
	}
	var buff bytes.Buffer
	for _, ve := range x {
		if ve == nil {
			continue
		}
		buff.WriteString(` [` + ve.Error() + `] `)
	}
	res := strings.TrimSpace(buff.String())
	return res
}

//-----------------------------------------------------------------------------

// Group just a wrapper for sync.WaitGroup for now, might get other parts like
// as in errgroup
type Group struct {
	wg sync.WaitGroup
}

// NewGroup creates *Group
func NewGroup() *Group {
	return &Group{}
}

// Wait waits for group to finish
func (g *Group) Wait() {
	g.wg.Wait()
	return
}

// Go runs f and registers it in sync.WaitGroup
func (g *Group) Go(f func()) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		f()
	}()
}

//-----------------------------------------------------------------------------
