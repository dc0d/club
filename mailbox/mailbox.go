package mailbox

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
	Send() chan<- interface{}
	Receive() <-chan interface{}
	Close() error
}

type mailbox struct {
	close   chan struct{}
	send    chan interface{}
	receive chan interface{}
	mails   Storage
}

func (mb *mailbox) Send() chan<- interface{}    { return mb.send }
func (mb *mailbox) Receive() <-chan interface{} { return mb.receive }
func (mb *mailbox) Close() error                { close(mb.close); return nil }

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
			if mb.mails.Len() > 0 {
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
