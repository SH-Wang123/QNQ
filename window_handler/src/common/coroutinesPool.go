package common

import (
	"sync"
)

var globalCoroPool *CoroutinesPool

type QTask interface {
	Execute(v ...interface{})
}

type CoroutinesPool struct {
	TaskChannel chan func(v ...interface{})
	GoNum       int
	Wg          *sync.WaitGroup
}

func NewFixedPool(cap int) *CoroutinesPool {
	var n int
	if n == 0 {
		n = 1
	}

	p := &CoroutinesPool{
		TaskChannel: make(chan func(v ...interface{})),
		GoNum:       n,
		Wg:          &sync.WaitGroup{},
	}

	return p
}

func (p *CoroutinesPool) StartPool() {
	for i := 0; i < p.GoNum; i++ {
		go func() {
			for task := range p.TaskChannel {
				task()
			}
			//select {
			//case task := <-p.TaskChannel:
			//	task()
			//}
		}()
	}
}

func (p *CoroutinesPool) Submit(executeFunc func(v ...interface{})) {
	p.TaskChannel <- executeFunc
}

// TODO 并发问题
func GetCoroutinesPool() *CoroutinesPool {
	return globalCoroPool
}
