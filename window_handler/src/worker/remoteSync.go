package worker

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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

func newRemoteSyncSender() *common.QSender {
	return &common.QSender{
		SN:                 common.GetSNCount(),
		Active:             false,
		Status:             common.TASK_FREE,
		RecCount:           1,
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
	num := 1
	remoteIp := fmt.Sprintf("%v", s.PrivateVariableMap["remoteIp"])
	workerSignal := common.GetRQPOptSignal(s.SN, common.TYPE_REMOTE_SINGLE, false, remoteIp)
	_, _ = common.WriteStrToQTarget([]byte(workerSignal), remoteIp, nil)
	initSignal := common.GetRQPInitSignal(s.SN, common.NULL_INIT_MAP, localFilePath)
	common.WriteStrToQTarget([]byte(initSignal), remoteIp, nil)
	time.Sleep(1000 * time.Millisecond)
	buf := make([]byte, 4082)
	log.Printf("send file start")
	common.StealConn(remoteIp, s)
	defer common.ReleaseConn(remoteIp, s)
	for {
		num++
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
		_, err = common.WriteStrToQTarget([]byte(common.LoadContent(msgPrefix, common.LoadQMQMsgNum(num), string(fBytes))), remoteIp, s)
		if err != nil {
			log.Printf("conn.Write error : %v", err.Error())
			break
		}
	}
	workerSignal = common.GetRQPOptSignal(s.SN, common.TYPE_REMOTE_SINGLE, true, remoteIp)
	time.Sleep(30 * time.Second)
	_, _ = common.WriteStrToQTarget([]byte(workerSignal), remoteIp, s)
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
	common.ConnectTarget(ip)
}

func GetQNetCell(ip string) *common.QNetCell {
	return common.GetQNetCell(ip)
}

func GetAllQNetCells() map[string]*common.QNetCell {
	return common.GetAllQNetCells()
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
