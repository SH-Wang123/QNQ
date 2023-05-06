package worker

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
		ExecuteFunc:     remoteSingleFileSyncReceiver,
		DeconstructFunc: receiverDeconstruct,
	}
}

func NewRemoteSyncSender() *common.QSender {
	return &common.QSender{
		SN:                 common.GetSNCount(),
		Active:             false,
		Status:             common.TASK_FREE,
		ExecuteFunc:        sendSingleFile,
		PrivateVariableMap: make(map[string]interface{}),
	}
}

func NewQNQAuthSender() *common.QSender {
	return &common.QSender{
		SN:                 common.GetSNCount(),
		Active:             false,
		Status:             common.TASK_FREE,
		ExecuteFunc:        requestQNQAuth,
		PrivateVariableMap: make(map[string]interface{}),
	}
}

func requestQNQAuth(s *common.QSender) {
	//激活认证接收者
	workerSignal := common.GetQMQTaskPre(common.TYPE_REMOTE_QNQ_AUTH) + s.SN + "0"
	_, _ = network.WriteStrToQTarget(workerSignal)
	//发送认证信息，携带IP和MAC
	var msgPrefix = dataMsgPreFix + s.SN

	network.WriteStrToQTarget(network.LoadContent(msgPrefix, config.SystemConfigCache.Value().QnqSTarget.RemotePath))
	//发送结束标志
}

func remoteSingleFileSyncReceiver(msg interface{}, w *common.QWorker) {
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

func sendSingleFile(s *common.QSender) {
	localFilePath := fmt.Sprintf("%v", s.PrivateVariableMap["local_file_path"])
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
	workerSignal := common.GetQMQTaskPre(common.TYPE_REMOTE_SINGLE) + s.SN + "0"
	_, _ = network.WriteStrToQTarget(workerSignal)
	network.WriteStrToQTarget(network.LoadContent(msgPrefix, config.SystemConfigCache.Value().QnqSTarget.RemotePath))
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
	workerSignal = common.GetQMQTaskPre(common.TYPE_REMOTE_SINGLE) + s.SN + "1"
	_, _ = network.WriteStrToQTarget(workerSignal)
}

func GetRemoteFileRootMap(ip string, abstractPath string, anchorPointPath string) map[string][]string {
	var params = make(map[string]string)
	params["abstractPath"] = abstractPath
	params["anchorPointPath"] = anchorPointPath
	resp, err := sendGet(URL_HRED+ip+common.QNQ_TARGET_REST_PORT+GET_FILE_ROOT_URI, params)
	if err != nil {
		return nil
	}
	var ret = make(map[string][]string)
	getObjFromResponse(resp, &ret)
	return ret
}

func TestQnqTarget(ip string) bool {
	resp, err := http.Get(URL_HRED + ip + common.QNQ_TARGET_REST_PORT + TEST_CONNECT)
	if err != nil {
		return false
	}
	var retStr string
	getObjFromResponse(resp, &retStr)
	return retStr == "ok"
}
