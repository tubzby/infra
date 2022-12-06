package broadcast

import "sync"

type Broadcast[T any] interface {
	Start()
	Stop()
	Register(chan T)
	UnRegister(chan T)
	Submit(T)
}

func New[T any]() Broadcast[T] {
	return &broadcast[T]{
		ch: make(chan T),
	}
}

type broadcast[T any] struct {
	l       sync.RWMutex
	clients []chan T
	ch      chan T
}

func (b *broadcast[T]) Start() {
	go func() {
		for {
			v, ok := <-b.ch
			if !ok {
				return
			}
			b.l.RLock()
			for _, ch := range b.clients {
				ch <- v
			}
			b.l.RUnlock()
		}
	}()
}

func (b *broadcast[T]) Register(ch chan T) {
	b.l.Lock()
	defer b.l.Unlock()

	b.clients = append(b.clients, ch)
}

func (b *broadcast[T]) UnRegister(ch chan T) {
	b.l.Lock()
	defer b.l.Unlock()

	l := len(b.clients)
	for i, v := range b.clients {
		if v == ch {
			b.clients[i], b.clients[l-1] = b.clients[l-1], b.clients[i]
			b.clients = b.clients[:l-1]
			return
		}
	}
}

func (b *broadcast[T]) Submit(value T) {
	b.l.RLock()
	defer b.l.RUnlock()
	b.ch <- value
}

func (b *broadcast[T]) Stop() {
	b.l.Lock()
	defer b.l.Unlock()
	close(b.ch)
}
