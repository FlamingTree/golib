package waitquit

import (
	"sync"
)

type WaitQuit struct {
	*sync.WaitGroup
	exitChan chan struct{}
}

func NewWaitQuit() WaitQuit {
	return WaitQuit{
		WaitGroup: &sync.WaitGroup{},
		exitChan:  make(chan struct{}),
	}
}

func (wq WaitQuit) Exit() {
	close(wq.exitChan)
	wq.Wait()
}

func (wq WaitQuit) Wrap(cb func()) {
	wq.Add(1)
	go func() {
		cb()
		wq.WaitGroup.Done()
	}()
}

func (wq WaitQuit) Done() <-chan struct{} {
	return wq.exitChan
}
