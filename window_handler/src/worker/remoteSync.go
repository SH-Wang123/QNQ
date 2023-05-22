package worker

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
	"window_handler/common"
)

var outputFileLock = sync.Mutex{}

func NewRemoteSyncReceiver(SN string) *common.QWorker {
	return &common.QWorker{
		SN:              SN,
		Active:          true,
		Sub:             make(chan any),
		OverChan:        make(chan int),
		Status:          common.TASK_FREE,
		ExecuteFunc:     remoteSingleFileSyncReceiver,
		DeconstructFunc: receiverDeconstruct,
	}
}

func newRemoteSyncSender(sn string) *common.QSender {
	if sn == "" {
		sn = common.GetSNCount()
	}
	return &common.QSender{
		SN:                 sn,
		Active:             false,
		Status:             common.TASK_FREE,
		RecCount:           1,
		ExecuteFunc:        sendSingleFile,
		PrivateVariableMap: make(map[string]interface{}),
	}
}

func RemoteSingleSyncSingleTime(localPath string, remotePath string, ip string) {
	sn, lock, startTime := preSyncSingleTime(common.TYPE_REMOTE_SINGLE)
	defer afterSyncSingleTime(common.TYPE_REMOTE_SINGLE)
	GetTotalSize(&sn, localPath, false, lock)
	lock.Done()
	sender := newRemoteSyncSender(sn)
	sender.PrivateVariableMap["local_file_path"] = localPath
	sender.PrivateVariableMap["remoteIp"] = ip
	sender.GetExecuteFunc()(sender)
	recordOLog(
		common.TYPE_REMOTE_SINGLE,
		startTime,
		localPath,
		localPath)
}

func NewQNQAuthReceiver(SN string) *common.QWorker {
	return &common.QWorker{
		SN:              SN,
		Active:          true,
		Status:          common.TASK_FREE,
		ExecuteFunc:     qnqAuthReceiver,
		DeconstructFunc: receiverDeconstruct,
	}
}

func qnqAuthReceiver(msg interface{}, w *common.QWorker) {

}

func remoteSingleFileSyncReceiver(msg interface{}, w *common.QWorker) {
	outputFileLock.Lock()
	defer outputFileLock.Unlock()
	msgStr := fmt.Sprintf("%v", msg)

	//Create File
	if w.Status == common.TASK_FREE {
		filePath := common.GetInitRQPMsg(msgStr)
		var err error
		log.Printf("RemoteSingleFileSyncReceiver %v: get file path %v", w.SN, filePath)
		w.PrivateFile, err = common.OpenFile(filePath, true)
		if err != nil {
			return
		}
		w.Status = common.TASK_READY
		return
	}

	w.Status = common.TASK_RUNNING
	//Write File
	if len(msgStr) == 0 {
		return
	}
	rm := msgStr[8:]
	num := msgStr[0:8]
	log.Printf("%v handle msg num : %v", w.SN, num)
	w.PrivateFile.Write([]byte(rm))
}

func receiverDeconstruct(w *common.QWorker) {
	if w.PrivateFile != nil {
		w.PrivateFile.Close()
	}
}

func sendSingleFile(s *common.QSender) {
	localFilePath := fmt.Sprintf("%v", s.PrivateVariableMap["local_file_path"])
	common.SetCurrentSyncFile(s.SN, SYNC_RUNNING, localFilePath)
	if localFilePath == "" || localFilePath == "remoteFilePath" {
		log.Printf("SendSingleFile QSender SN : {%v}, get a null file path", s.SN)
		return
	}
	var msgPrefix = dataMsgPreFix + s.SN
	f, err := os.Open(localFilePath)
	if err != nil {
		log.Printf("Open %v err : %v", localFilePath, err.Error())
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("close file err : %v", err.Error())
		}
	}(f)
	num := 0
	remoteIp := fmt.Sprintf("%v", s.PrivateVariableMap["remoteIp"])
	workerSignal := common.GetRQPOptSignal(s.SN, common.TYPE_REMOTE_SINGLE, false, remoteIp)
	_, _ = common.WriteStrToQTarget([]byte(workerSignal), remoteIp, nil)
	time.Sleep(1000 * time.Millisecond)
	initSignal := common.GetRQPInitSignal(s.SN, common.NULL_INIT_MAP, localFilePath)
	log.Printf("send init signal : %s", initSignal)
	common.WriteStrToQTarget([]byte(initSignal), remoteIp, nil)
	time.Sleep(1000 * time.Millisecond)
	buf := make([]byte, 4082)
	log.Printf("send file start")
	common.StealConn(remoteIp, s)
	defer common.ReleaseConn(remoteIp, s)
	count := 0
	for {
		time.Sleep(1 * time.Millisecond)
		num++
		count++
		if num%200 == 0 {
			log.Printf("send msg num : %d, loss pck : %v, SN : %s", num, count != 200, s.SN)
			count = 0
		}
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("send file over")
			} else {
				log.Printf("f.Read error : %v", err.Error())
			}
			break
		}
		addSizeToDoneMap(s.SN, uint64(n))
		var fBytes = buf[:n]
		numStr := common.LoadQMQMsgNum(num)
		_, err = common.WriteStrToQTarget([]byte(common.LoadContent(msgPrefix, numStr, string(fBytes))), remoteIp, s)
		if err != nil {
			log.Printf("conn.Write error : %v", err.Error())
			break
		}
	}
	workerSignal = common.GetRQPOptSignal(s.SN, common.TYPE_REMOTE_SINGLE, true, remoteIp)
	_, _ = common.WriteStrToQTarget([]byte(workerSignal), remoteIp, s)
}
