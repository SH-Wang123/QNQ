package worker

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"window_handler/common"
	"window_handler/config"
	"window_handler/network"
)

var outputFileLock = sync.Mutex{}

func NewRemoteSyncReceiver(SN string) *common.QWorker {
	return &common.QWorker{
		SN:              SN,
		Active:          true,
		Status:          common.TASK_FREE,
		ExecuteFunc:     RemoteSingleFileSyncReceiver,
		DeconstructFunc: receiverDeconstruct,
	}
}

func NewRemoteSyncSender() *common.QSender {
	return &common.QSender{
		SN:                 common.GetSNCount(),
		Active:             false,
		Status:             common.TASK_FREE,
		ExecuteFunc:        SendSingleFile,
		PrivateVariableMap: make(map[string]interface{}),
	}
}

func RemoteSingleFileSyncReceiver(msg interface{}, w *common.QWorker) {
	outputFileLock.Lock()
	defer outputFileLock.Unlock()
	msgStr := fmt.Sprintf("%v", msg)

	//Create File
	if w.Status == common.TASK_READY {
		var err error
		log.Printf("RemoteSingleFileSyncReceiver : get file path %v", msgStr)
		w.PrivateFile, err = os.Open(msgStr)
		if err != nil {
			w.PrivateFile, _ = os.Create(msgStr)
		}
		w.Status = common.TASK_RUNNING
		return
	} else if w.Status != common.TASK_RUNNING {
		log.Printf("%v : error opt", w.SN)
		return
	}

	//Write File
	if len(msgStr) == 0 {
		return
	}
	w.PrivateFile.Write([]byte(msgStr))
}

func receiverDeconstruct(w *common.QWorker) {
	if w.PrivateFile != nil {
		w.PrivateFile.Close()
	}
}

func SendSingleFile(s *common.QSender) {
	filePath := fmt.Sprintf("%v", s.PrivateVariableMap["file_path"])
	if filePath == "" {
		log.Printf("SendSingleFile QSender SN : {%v}, get a null file path", s.SN)
		return
	}
	var msgPrefix = dataMsgPreFix + s.SN
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Open %v err : %v", filePath, err.Error())
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Printf("close file err : %v", err.Error())
		}
	}(f)
	workerSignal := common.RemoteSingleSyncType + s.SN + "0"
	_, _ = network.WriteStrToQTarget(workerSignal)
	network.WriteStrToQTarget(network.LoadContent(msgPrefix, config.SystemConfigCache.Value().QnqSTarget.LocalPath))
	buf := make([]byte, 4094)
	var msgfixBytes = []byte{'1', '1'}
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("send file over")
			} else {
				log.Printf("f.Read error : %v", err.Error())
			}
			break
		}
		var fBytes = append(msgfixBytes, buf[:n]...)
		_, err = network.WriteStrToQTarget(network.LoadContent(msgPrefix, string(fBytes)))
		if err != nil {
			log.Printf("conn.Write error : %v", err.Error())
			break
		}
	}
	workerSignal = common.RemoteSingleSyncType + s.SN + "1"
	_, _ = network.WriteStrToQTarget(workerSignal)
}
