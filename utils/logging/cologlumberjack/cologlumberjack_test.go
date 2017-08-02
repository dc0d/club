package cologlumberjack

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/comail/colog"
	"github.com/dc0d/goroutines"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
	wg     = &sync.WaitGroup{}
)

func TestMain(m *testing.M) {
	defer func() {
		goroutines.New().
			Timeout(time.Second).
			Go(func() {
				wg.Wait()
			})
	}()
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	m.Run()
}

var _ colog.Hook = &Hook{}

func TestDummy(t *testing.T) {
	logger := &lumberjack.Logger{
		Filename:   "/tmp/cologlumberjack_test.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     1,
	}
	defer logger.Close()

	colog.Register()
	colog.SetFormatter(&colog.StdFormatter{})
	colog.AddHook(New(logger))
	colog.SetFlags(log.Lshortfile | log.Ltime)
	colog.ParseFields(true)

	log.Println("error: OK OK x=1 y=\"W O\"")
	log.Println("OK")
	log.Println(`OK1
OK2`)
}
