// Package singletongroup is most likely a bad practice.
package singletongroup

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dc0d/ctxgrp"
)

//-----------------------------------------------------------------------------

// Singleton returns the singleton instance of StopGroup
func Singleton() StopGroup { return grp }

var (
	grp = newGroup()
)

//-----------------------------------------------------------------------------

// creates a group that gets canceled on signal
func newGroup() StopGroup {
	ctx, cancel := context.WithCancel(context.Background())
	onSignal(func() { cancel() })
	tg := ctxgrp.New(ctx)

	return &sg{
		grp:    tg,
		cancel: cancel,
	}
}

//-----------------------------------------------------------------------------

// StopGroup a group that can be stopped
type StopGroup interface {
	Group() (context.Context, ctxgrp.WaitGroup)
	SignalStop()
}

//-----------------------------------------------------------------------------

type sg struct {
	grp    ctxgrp.Group
	cancel context.CancelFunc
}

func (x *sg) SignalStop() {
	if x.cancel == nil {
		return
	}
	x.cancel()
}

func (x *sg) Group() (context.Context, ctxgrp.WaitGroup) { return x.grp() }

//-----------------------------------------------------------------------------

func onSignal(f func(), sig ...os.Signal) {
	if f == nil {
		return
	}
	sigc := make(chan os.Signal, 1)
	if len(sig) > 0 {
		signal.Notify(sigc, sig...)
	} else {
		signal.Notify(sigc,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGSTOP,
			syscall.SIGTSTP,
			syscall.SIGKILL)
	}
	go func() {
		<-sigc
		f()
	}()
}

//-----------------------------------------------------------------------------
