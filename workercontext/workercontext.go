package workercontext

import (
	"context"
	"sync"

	"github.com/dc0d/club/errors"
)

// WaitGroup interface for built-in WaitGroup
type WaitGroup interface {
	Add(delta int)
	Done()
	Wait()
}

type workerContext struct {
	ctx context.Context
	wg  *sync.WaitGroup
}

func (wcx *workerContext) Context() context.Context { return wcx.ctx }
func (wcx *workerContext) WaitGroup() WaitGroup     { return wcx.wg }

// WorkerContext combination of context.Context & WaitGroup, an execution context
type WorkerContext interface {
	Context() context.Context
	WaitGroup() WaitGroup
}

// ErrNilContext means we got a nil context while we shouldn't
var ErrNilContext = errors.Errorf("ERR_NIL_CONTEXT")

// New .
func New(ctx context.Context) (WorkerContext, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}
	return &workerContext{
		ctx: ctx,
		wg:  &sync.WaitGroup{},
	}, nil
}
