package tgroup

import "time"

// Wait waits for a WaitGroup with timeout (if provided)
func Wait(wg WaitGroup, timeout ...time.Duration) error {
	if len(timeout) == 0 || timeout[0] <= 0 {
		wg.Wait()
		return nil
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()
	select {
	case <-done:
	case <-time.After(timeout[0]):
		return ErrTimeout
	}
	return nil
}

// WaitFinish waits for a Group to finish, first locks
// (<-ctx.Done()) then waits for the WaitGroup.
func WaitFinish(
	g Group,
	timeout ...time.Duration) error {
	ctx, wg := g.Group()
	<-ctx.Done()
	return Wait(wg, timeout...)
}
