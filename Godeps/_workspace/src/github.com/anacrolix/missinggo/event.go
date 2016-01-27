package missinggo

import "sync"

type Event struct {
	mu     sync.Mutex
	ch     chan struct{}
	closed bool
}

func (me *Event) C() <-chan struct{} {
	me.mu.Lock()
	defer me.mu.Unlock()
	if me.ch == nil {
		me.ch = make(chan struct{})
	}
	return me.ch
}

func (me *Event) Set() {
	me.mu.Lock()
	defer me.mu.Unlock()
	if me.closed {
		return
	}
	close(me.ch)
	me.closed = true
}
