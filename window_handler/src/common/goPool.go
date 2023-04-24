package common

import (
	"reflect"
	"sync"
)

var newGoPool *goPool

type Task struct {
	status   int
	function reflect.Value
	param    []reflect.Value
}

type goPool struct {
	taskChannel chan *Task
	goNum       int
	wg          *sync.WaitGroup
}

func newFixedGoPool(cap int) *goPool {
	newGoPool = &goPool{
		taskChannel: make(chan *Task),
		goNum:       10,
		wg:          &sync.WaitGroup{},
	}
	newGoPool.start()
	return newGoPool
}

func SubmitFunc2Pool(f reflect.Value, params ...any) {
	var refParams []*reflect.Value
	for _, v := range params {
		refP := reflect.ValueOf(v)
		refParams = append(refParams, &refP)
	}
	newGoPool.submit(f, refParams)
}

func SubmitTask2Pool() {

}

func (w *Task) execute() {
	w.function.Call(w.param)
}

func (g *goPool) start() {
	for i := 0; i < g.goNum; i++ {
		go func() {
			for task := range g.taskChannel {
				t := *task
				t.execute()
			}
		}()
	}
}

func (g *goPool) submit(f reflect.Value, params []*reflect.Value) {
	task := &Task{}
	task.function = f
	if params != nil {
		for _, v := range params {
			task.param = append(task.param, *v)
		}
	}
	g.taskChannel <- task
}
