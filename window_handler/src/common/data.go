package common

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// ------------------------------ QMQ

/**

RQP:

任务消息：
0x       | 000    | 0000    | 0
任务消息头 | 任务编码 | 任务序号(SN) | 任务状态位

任务初始化消息：
01 | 0000 | 0000 | 00...
消息头 1| 任务序号 5| 初始化配置位图9| 自定义消息

数据消息：（最大长度4096byte）
00       | 0000    | 00000000 | 00..00 | 00000000
数据消息头 | 任务序号(SN) | 消息序列 | 数据段 | 校验位

本地任务激活：
0x | 00000000
无实际内容

*/

const TaskOverFlag = "1"
const (
	optRQPHead    = "0x"
	initRQPHead   = "01"
	msgRQPHead    = "00"
	NULL_INIT_MAP = "0000"
)
const (
	//远程业务，0x011 ~ 0x100
	remoteSingleSyncPre = optRQPHead + "011"
	//非业务，0x101 ~ 0x200
	qnqAuthPre = optRQPHead + "101"
)

func GetRQPTaskPre(busType int) string {
	switch busType {
	case TYPE_REMOTE_SINGLE:
		return remoteSingleSyncPre
	case TYPE_REMOTE_QNQ_AUTH:
		return qnqAuthPre
	}
	return "0x000"
}

func GetRQPOptSignal(sn string, busType int, overFlag bool, remoteIp string) string {
	time.Sleep(100 * time.Millisecond)
	pre := GetRQPTaskPre(busType)
	var signal string
	if overFlag {
		signal = pre + sn + fmt.Sprint(1)
	} else {
		signal = pre + sn + fmt.Sprint(0)
	}
	log.Printf("send signal : %v, ip : %v", signal, remoteIp)
	return signal
}

func GetRQPInitSignal(sn string, initMap string, message string) string {
	return initRQPHead + sn + initMap + message
}

const (
	TASK_FREE = iota
	TASK_READY
	TASK_RUNNING
	TASK_OVER
)

// business type
const (
	placeholder = iota
	TYPE_LOCAL_BATCH
	TYPE_LOCAL_SING
	TYPE_PARTITION
	TYPE_REMOTE_SINGLE
	TYPE_REMOTE_BATCH
	TYPE_CREATE_TIMEPOINT
	TYPE_TEST_SPEED
	TYPE_REMOTE_QNQ_AUTH
)

const (
	SIGNAL_AUTH_PASS    = 1000
	SIGNAL_AUTH_NO_PASS = 1001
)

func GetForceDoneSignal(busType int) int {
	return busType * -1
}

func GetRunningSignal(busType int) int {
	return busType
}

var WorkerFactoryMap = map[string]func(SN string) *QWorker{}

type QWorker struct {
	SN              string
	Active          bool
	Status          int
	Sub             chan interface{}
	OverChan        chan int
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

func (w *QWorker) Execute(v ...interface{}) {
	defer w.Deconstruct()
	if w.Sub == nil {
		w.ExecuteFunc(nil, w)
		return
	}
	qmqWaitGroup.Add(1)
	defer qmqWaitGroup.Done()
	for {
		select {
		case m := <-w.Sub:
			if m != nil {
				w.ExecuteFunc(m, w)
			}
		case <-w.OverChan:
			log.Printf("Consumer %s over", w.SN)
			return
		}
	}
}

type QSender struct {
	SN                 string
	Active             bool
	Status             int
	RecCount           int
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

// ------------------------------ network

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

type QNetCell struct {
	QTarget        *net.Conn
	currentTWorker *QSender
	QServer        *net.Conn
	currentSWorker *QSender
	netCellLock    *sync.RWMutex
	//target status | server status
	status int
}

func (qn *QNetCell) setTargetStatus(status bool) {
	if status {
		qn.status = qn.status | 10
	} else {
		qn.status = qn.status & 01
	}
}

func (qn *QNetCell) setServerStatus(status bool) {
	if status {
		qn.status = qn.status | 01
	} else {
		qn.status = qn.status & 10
	}
}

func (qn *QNetCell) GetTargetStatus() bool {
	return qn.status&10 >= 10
}

func (qn *QNetCell) GetServerStatus() bool {
	return qn.status&01 == 1
}
