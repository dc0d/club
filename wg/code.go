package wg

func init() {
	println("deprecated")
}

// import (
// 	"sync"
// 	"time"

// 	"github.com/dc0d/errgo/sentinel"
// )

// //-----------------------------------------------------------------------------

// // WaitGroup .
// type WaitGroup interface {
// 	Add(delta int)
// 	Done()
// 	Wait(...time.Duration) error
// }

// // New creates an instance of WaitGroup
// func New() WaitGroup { return &waitGroup{} }

// //-----------------------------------------------------------------------------

// type waitGroup struct {
// 	wg sync.WaitGroup
// }

// func (wg *waitGroup) Add(delta int) { wg.wg.Add(delta) }

// func (wg *waitGroup) Done() { wg.wg.Done() }

// func (wg *waitGroup) Wait(timeout ...time.Duration) error {
// 	if len(timeout) == 0 || timeout[0] <= 0 {
// 		wg.wg.Wait()
// 		return nil
// 	}
// 	done := make(chan struct{})
// 	go func() {
// 		defer close(done)
// 		wg.wg.Wait()
// 	}()
// 	select {
// 	case <-done:
// 	case <-time.After(timeout[0]):
// 		return ErrTimeout
// 	}
// 	return nil
// }

// //-----------------------------------------------------------------------------

// // errors
// var (
// 	ErrTimeout = sentinel.Errorf("TIMEOUT")
// )

// //-----------------------------------------------------------------------------
