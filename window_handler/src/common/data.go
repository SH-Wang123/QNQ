package common

import (
	"os"
)

// ------------------------------ QMQ

/**
任务消息：
0x       | 000    | 0000    | 0
任务消息头 | 任务编码 | 任务序号(SN) | 任务状态位

任务初始化消息：
01 | 0000 | 0000 | 00
消息头 1| 任务序号 5| 初始化配置位图9| 冗余位11

数据消息：（最大长度4096byte）
00       | 0000    | 00..00 | 00000000
数据消息头 | 任务序号(SN) | 数据段 | 校验位

本地任务激活：
0x | 00000000
无实际内容

*/

const (
	TASK_FREE = iota
	TASK_READY
	TASK_RUNNING
	TASK_OVER
)

// business type
const (
	TYPE_LOCAL_BATCH = iota
	TYPE_LOCAL_SING
	TYPE_PARTITION
	TYPE_REMOTE_SINGLE
	TYPE_REMOTE_BATCH
	TYPE_CREATE_TIMEPOINT
	TYPE_TEST_SPEED
)

// channel signal
const (
	localBatchRunning = iota
	localBatchForceDone
	localSingleRunning
	localSingleForceDone
	partitionRunning
	partitionForceDone
	createTimePointRunning
	createTimePointForceDone
	testDiskSpeedRunning
	testDiskSpeedForceDone
)

func GetForceDoneSignal(busType int) int {
	switch busType {
	case TYPE_LOCAL_BATCH:
		return localBatchForceDone
	case TYPE_LOCAL_SING:
		return localSingleForceDone
	case TYPE_PARTITION:
		return partitionForceDone
	case TYPE_CREATE_TIMEPOINT:
		return createTimePointForceDone
	case TYPE_TEST_SPEED:
		return testDiskSpeedForceDone
	}
	return -1
}

func GetRunningSignal(busType int) int {
	switch busType {
	case TYPE_LOCAL_BATCH:
		return localBatchRunning
	case TYPE_LOCAL_SING:
		return localSingleRunning
	case TYPE_PARTITION:
		return partitionRunning
	case TYPE_CREATE_TIMEPOINT:
		return createTimePointRunning
	case TYPE_TEST_SPEED:
		return testDiskSpeedRunning
	}
	return -1
}

const TaskOverFlag = "1"
const RemoteSingleSyncType = "0x010"

var WorkerFactoryMap = map[string]func(SN string) *QWorker{}

type QWorker struct {
	SN              string
	Active          bool
	Status          int
	Sub             chan interface{}
	ExecuteFunc     func(msg interface{}, w *QWorker)
	DeconstructFunc func(w *QWorker)
	PrivateFile     *os.File //usually source file
	TargetFile      *os.File
	PrivateNet      os.File
	Md5CacheFlag    bool
}

func (w *QWorker) Deconstruct() {
	w.DeconstructFunc(w)
}

type QSender struct {
	SN                 string
	Active             bool
	Status             int
	ExecuteFunc        func(s *QSender)
	PrivateVariableMap map[string]interface{}
}

func (s *QSender) GetExecuteFunc() func(s *QSender) {
	return s.ExecuteFunc
}

// ------------------------------ Observer

type Observer interface {
	UpdateAd(interface{})
	GetName() string
	SetName(string)
}

type Subject interface {
	Register(Observer)
	Deregister(Observer)
	NotifyAll()
}

// ------------------------------ rest entry
const QNQ_TARGET_REST_PORT = ":9915"

type QResponse struct {
	Code int `json:"code"`
	Data any `json:"data"`
}

func NewQResponse(code int, data any) *QResponse {
	return &QResponse{
		Code: code,
		Data: data,
	}
}
