package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

var taskCount = 0
var countLock sync.Mutex

func GetSNCount() string {
	SNCountStr := fmt.Sprintf("%v", GetTaskCount())
	for i := len(SNCountStr); i < 4; i++ {
		SNCountStr = "0" + SNCountStr
	}
	return SNCountStr
}

func GetTaskCount() string {
	countLock.Lock()
	if taskCount == 9999 {
		taskCount = 0
	}
	taskCount++
	countLock.Unlock()
	return fmt.Sprintf("%v", taskCount)
}

func GetNowTimeStr() string {
	now := time.Now()
	ret := fmt.Sprintf("%v", now.Format("2006/01/02 15:04:05"))
	return ret
}

func GetIpFromAddr(addr string) string {
	s := strings.Split(addr, ":")
	if len(s) == 2 {
		return s[0]
	} else {
		return ""
	}
}

//------------------------------------network util

// LoadQMQMsgNum 装载QMQ消息的序列
func LoadQMQMsgNum(num int) string {
	numStr := fmt.Sprint(num)
	l := 8 - len(numStr)
	for i := 0; i < l; i++ {
		numStr = fmt.Sprintf("0%s", numStr)
	}
	return numStr
}

func writeToConn(conn *net.Conn, content []byte) (err error) {
	l := len(content)
	magicNum := make([]byte, 4)
	binary.BigEndian.PutUint32(magicNum, 0x123456)
	lenNum := make([]byte, 2)
	binary.BigEndian.PutUint16(lenNum, uint16(l))
	packetBuf := bytes.NewBuffer(magicNum)
	packetBuf.Write(lenNum)
	packetBuf.Write(content)
	connR := *conn
	_, err = connR.Write(packetBuf.Bytes())
	if err != nil {
		fmt.Printf("write failed , err : %v\n", err)
		return
	}
	return
}

func loadContentWithCheck(prefix string, numStr string, msg string, checkBit string) string {
	return prefix + numStr + msg + checkBit
}

func LoadContent(prefix string, numStr string, msg string) string {
	return loadContentWithCheck(prefix, numStr, msg, "00000000")
}

// StealConn 窃取Target，标记持有方的QSender。只有发送端过快，接收端过慢才会出现粘包。
func StealConn(ip string, w *QSender) bool {
	if qNetCells[ip] == nil {
		return false
	}
	qNetCells[ip].netCellLock.Lock()
	defer qNetCells[ip].netCellLock.Unlock()
	if qNetCells[ip].currentTWorker != nil {
		return false
	}
	qNetCells[ip].currentTWorker = w
	return true
}

// ReleaseConn 只允许持有方释放
func ReleaseConn(ip string, w *QSender) bool {
	if qNetCells[ip] == nil {
		return false
	}
	qNetCells[ip].netCellLock.Lock()
	defer qNetCells[ip].netCellLock.Unlock()
	if qNetCells[ip].currentTWorker != w {
		return false
	}
	qNetCells[ip].currentTWorker = nil
	return true
}

func GetCellQSender(ip string) *QSender {
	if qNetCells[ip] != nil {
		return qNetCells[ip].currentTWorker
	} else {
		return nil
	}
}

func packetSlitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if !atEOF && len(data) > 6 && binary.BigEndian.Uint32(data[:4]) == 0x123456 {
		var l int16
		binary.Read(bytes.NewReader(data[4:6]), binary.BigEndian, &l)
		pl := int(l) + 6
		if pl <= len(data) {
			return pl, data[:pl], nil
		}
	}
	return
}

func GetInitRQPMsg(rqpStr string) string {
	return rqpStr[4:]
}
