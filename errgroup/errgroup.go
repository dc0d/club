// Package errgroup adapted from https://github.com/golang/sync/blob/master/errgroup/errgroup.go
// tests are the same; their license is BSD.
package errgroup

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"context"
	"sync"

	"github.com/dc0d/errgo/sentinel"
)

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid and does not cancel on error.
type Group struct {
	cancel func()

	wg sync.WaitGroup

	errOnce sync.Once
	err     error
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// Go calls the given function in a new goroutine.
//
// The first call to return a non-nil error cancels the group; its error will be
// returned by Wait.
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()
		var funcerr error
		defer func() {
			if funcerr == nil {
				return
			}
			g.errOnce.Do(func() {
				g.err = funcerr
				if g.cancel != nil {
					g.cancel()
				}
			})
		}()
		defer func() {
			if e := recover(); e != nil {
				if err, ok := e.(error); ok {
					funcerr = err
					return
				}
				funcerr = sentinel.Errorf("UNKNOWN: %v", e)
			}
		}()

		if err := f(); err != nil {
			funcerr = err
		}
	}()
}
