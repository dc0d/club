package club

import (
	"testing"
	"time"
)

func TestGroup(t *testing.T) {
	g := NewGroup()

	items := make(chan int, 2)
	g.Go(func() { items <- 1 })
	g.Go(func() { items <- 1 })
	g.Wait()
	close(items)
	acc := 0
	for v := range items {
		acc += v
	}
	if acc != 2 {
		t.Fail()
	}
}

func TestApp(t *testing.T) {
	err := WaitGo(func() {
		defer Finit(-1, time.Millisecond)
		AppPool.Go(func() { <-AppCtx.Done() })
		AppCancel()
		AppPool.Wait()
	}, time.Second*3)
	if err != nil {
		t.Fail()
	}
}

func TestErrors(t *testing.T) {
	var es Errors
	es = append(es, Error(`ONE`))
	es = append(es, Error(`TWO`))
	if es.String() != `[ONE] [TWO]` {
		t.Fail()
	}
}
