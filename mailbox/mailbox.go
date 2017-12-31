package mailbox

import (
	"time"
)

// Storage like a []interface{} slice
type Storage interface {
	Len() int
	Peek() interface{}
	Drop()
	Append(interface{})
}

// SliceStorage .
type SliceStorage []interface{}

// Len .
func (store SliceStorage) Len() int { return len(store) }

// Append .
func (store *SliceStorage) Append(v interface{}) { *store = append(*store, v) }

// Peek .
func (store SliceStorage) Peek() interface{} { return store[0] }

// Drop .
func (store *SliceStorage) Drop() { *store = (*store)[1:] }

// Mailbox .
type Mailbox interface {
	Send(interface{}, ...time.Duration) bool
	Receive(...time.Duration) (interface{}, bool)
	Close() error
}

type mailbox struct {
	close   chan struct{}
	send    chan interface{}
	receive chan interface{}
	mails   Storage
}

func (mb *mailbox) Send(v interface{}, timeout ...time.Duration) bool {
	var toch <-chan time.Time
	if len(timeout) > 0 && timeout[0] > 0 {
		toch = time.After(timeout[0])
	}
	select {
	case <-toch:
		return false
	case mb.send <- v:
	}
	return true
}

func (mb *mailbox) Close() error { close(mb.close); return nil }

func (mb *mailbox) Receive(timeout ...time.Duration) (interface{}, bool) {
	var toch <-chan time.Time
	if len(timeout) > 0 && timeout[0] > 0 {
		toch = time.After(timeout[0])
	}
	select {
	case <-toch:
		return nil, false
	case v, ok := <-mb.receive:
		return v, ok
	}
}

func (mb *mailbox) loop() {
	defer close(mb.receive)
	var actualReceive chan interface{}
	first := func() interface{} {
		if mb.mails.Len() == 0 {
			return nil
		}
		return mb.mails.Peek()
	}
	for {
		select {
		case <-mb.close:
			if mb.mails.Len() > 0 { // (?) this may cause to not close ever
				continue
			}
			return
		case v := <-mb.send:
			mb.mails.Append(v)
			actualReceive = mb.receive
		case actualReceive <- first():
			mb.mails.Drop()
			if mb.mails.Len() == 0 {
				actualReceive = nil
			}
		}
	}
}

// New .
func New(store Storage) Mailbox {
	res := &mailbox{
		close:   make(chan struct{}),
		send:    make(chan interface{}),
		receive: make(chan interface{}),
		mails:   store,
	}
	go res.loop()
	return res
}
