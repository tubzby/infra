package broadcast

import (
	"sync"
	"testing"
)

var sg sync.WaitGroup

func client(ch chan int) {
	for {
		v, ok := <-ch
		if !ok {
			return
		}
		sg.Add(v)
	}
}

func TestBCast(t *testing.T) {
	b := New[int]()
	b.Start()
	defer b.Stop()

	var chs [10]chan int
	for i := 0; i < 10; i++ {
		chs[i] = make(chan int)
		b.Register(chs[i])
		go client(chs[i])
	}

	// 10 channels, 10 times
	sg.Add(100)

	for i := 0; i < 10; i++ {
		b.Submit(-1)
	}

	sg.Wait()
	for i := 0; i < 10; i++ {
		close(chs[i])
		b.UnRegister(chs[i])
	}
}
