package worker

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"window_handler/common"
	"window_handler/network"
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

func newRemoteSyncSender() *common.QSender {
	return &common.QSender{
		SN:                 common.GetSNCount(),
		Active:             false,
		Status:             common.TASK_FREE,
		ExecuteFunc:        sendSingleFile,
		PrivateVariableMap: make(map[string]interface{}),
	}
}

func RemoteSingleSyncSingleTime(localPath string, remotePath string, ip string) {
	sender := newRemoteSyncSender()
	sender.PrivateVariableMap["local_file_path"] = localPath
	sender.PrivateVariableMap["remoteIp"] = ip
	sender.GetExecuteFunc()(sender)
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
	if w.Status == common.TASK_READY {
		var err error
		log.Printf("RemoteSingleFileSyncReceiver : get file path %v", msgStr)
		w.PrivateFile, err = common.OpenFile(msgStr, true)
		if err != nil {
			return
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
	remoteIp := fmt.Sprintf("%v", s.PrivateVariableMap["remoteIp"])
	workerSignal := common.GetQMQTaskPre(common.TYPE_REMOTE_SINGLE) + s.SN + "0"
	_, _ = network.WriteStrToQTarget(workerSignal, remoteIp)
	network.WriteStrToQTarget(network.LoadContent(msgPrefix, localFilePath), remoteIp)
	buf := make([]byte, 1010)
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
		var fBytes = buf[:n]
		_, err = network.WriteStrToQTarget(network.LoadContent(msgPrefix, string(fBytes)), remoteIp)
		if err != nil {
			log.Printf("conn.Write error : %v", err.Error())
			break
		}
	}
	workerSignal = common.GetQMQTaskPre(common.TYPE_REMOTE_SINGLE) + s.SN + "1"
	_, _ = network.WriteStrToQTarget(workerSignal, remoteIp)
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

func ConnectTarget(ip string) {
	network.ConnectTarget(ip)
}

func GetQNetCell(ip string) *network.QNetCell {
	return network.GetQNetCell(ip)
}

func GetAllQNetCells() map[string]*network.QNetCell {
	return network.GetAllQNetCells()
}

// GetAllQSorT 获取当前network缓存中的所有server或target
func GetAllQSorT(serverFlag bool) []*net.Conn {
	cells := GetAllQNetCells()
	rets := make([]*net.Conn, 0)
	for _, cell := range cells {
		if serverFlag && cell.QServer != nil {
			rets = append(rets, cell.QServer)
		}
		if !serverFlag && cell.QTarget != nil {
			rets = append(rets, cell.QTarget)
		}
	}
	return rets
}
