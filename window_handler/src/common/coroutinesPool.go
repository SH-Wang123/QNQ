package common

import (
	"log"
	"runtime"
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
		n = runtime.NumCPU() * 32
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
	if globalCoroPool == nil {
		InitCoroutinesPool()
		globalCoroPool.StartPool()
	}
	return globalCoroPool
}

func InitCoroutinesPool() {
	if globalCoroPool == nil {
		globalCoroPool = NewFixedPool(0)
		log.Printf("Create a new coroutines pool, size : %v", globalCoroPool.GoNum)
	}
}
