package club

import (
	"testing"
	"time"

	"github.com/dc0d/goroutines"
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
	err := goroutines.New().
		WaitStart().
		WaitGo(time.Second * 3).
		Go(func() {
			defer Finit(-1, time.Millisecond)
			AppPool.Go(func() { <-AppCtx.Done() })
			AppCancel()
			AppPool.Wait()
		})
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
