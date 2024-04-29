package QNQ

import (
	"bytes"
	"context"
	"log/slog"
	"sync"
	"time"
)

var (
	taskChan = make(chan *Task, 256)
	overChan = make(chan struct{})
	taskWG   sync.WaitGroup
	taskMap  = make(map[uint32]*Task)
)

var (
	currentId = uint32(0)
	idLock    = sync.Mutex{}
)

func newId() uint32 {
	idLock.Lock()
	defer idLock.Unlock()
	currentId++
	return currentId
}

type TaskResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (tr *TaskResult) SetCode(key int) {
	tr.Code = key
	tr.Message = GetMsg(key)
}

type DataBuffer struct {
	Id int
	bytes.Buffer
}

type Task struct {
	Id          uint32
	Buffer      chan DataBuffer
	Context     context.Context
	Status      uint32
	Result      *TaskResult
	Probe       *ProgressProbe
	executeFunc func(t *Task)
	cancelFunc  func(t *Task)
}

func (t *Task) Execute() {
	slog.Info("task start : ", "id", t.Id, "params", t.Context.Value("params"))
	defer slog.Info("task done : ", "id", t.Id, "res", t.Result)
	t.executeFunc(t)
}

func (t *Task) Cancel() {
	t.cancelFunc(t)
}

func NewTask(executeFunc func(t *Task), cancelFunc func(t *Task), key any, val any) *Task {
	task := &Task{
		Id:          newId(),
		Buffer:      make(chan DataBuffer),
		Context:     context.WithValue(context.Background(), key, val),
		Result:      &TaskResult{},
		executeFunc: executeFunc,
		cancelFunc:  cancelFunc,
		Probe:       NewProgressProbe(),
	}
	taskMap[task.Id] = task
	return task
}

type TaskCache struct {
	preTask  *Task
	nextTask *Task
}

func (tc *TaskCache) GetTask(id int) {

}

func (tc *TaskCache) AddTask() {

}

type TaskGroup struct {
	Id int
}

const goNum = 10

func initTaskSchedule() {
	for i := 0; i < goNum; i++ {
		go func() {
			for {
				select {
				case t, ok := <-taskChan:
					if t != nil && ok {
						t.Execute()
						taskWG.Done()
					}
				case <-overChan:
					return
				}
			}
		}()
	}
}

func CommitTask(t *Task) {
	taskWG.Add(1)
	t.Probe.watch()
	taskChan <- t
}

func WaitTaskSchedulerStop() {
	taskWG.Wait()
	close(overChan)
	time.Sleep(1 * time.Second)
	close(taskChan)
}

func getTask(id uint32) *Task {
	if v, ok := taskMap[id]; ok {
		return v
	}
	return nil
}

type TaskCenter struct {
	UnimplementedTaskCenterServer
}

func (t *TaskCenter) GetTaskInfo(ctx context.Context, req *GetTaskRequest) (*TaskInfoResult, error) {
	var res *TaskInfoResult
	var err error
	task := getTask(req.GetTaskId())
	if task != nil {
		res = &TaskInfoResult{
			Status:    task.Status,
			Progress:  float64(task.Probe.GetProgress()),
			TotalSize: task.Probe.totalSize,
		}
	}
	return res, err
}
