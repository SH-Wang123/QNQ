package network

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
	"window_handler/common"
	"window_handler/config"
)

const (
	ServerNetworkType = "tcp"
	ServerPort        = ":9916"
	MessageDelimiter  = '\t'
	RecoverMessage    = "cx00000615"
	EndMessage        = '\n'
)

var AuthLock = &sync.RWMutex{}
var AuthFlag = false

var qNetCells = make(map[string]QNetCell, 8)

var NetChan = common.NewProducer()

func handleConnect(conn net.Conn, isClient bool) {
	log.Printf("client %v connected\n", conn.RemoteAddr())
	for {
		conn.SetReadDeadline(time.Now().Add(time.Minute * time.Duration(5)))
		if msg, err := readStr(conn); err != nil {
			if err == io.EOF {
				//log.Printf("client %v closed\n", conn.RemoteAddr())
			} else {
				//log.Printf("read error : %v\n", err.Error())
			}
		} else {
			if !isClient {
				write(conn, RecoverMessage)
			}
			go NetChan.Produce(msg)
		}
	}
}

func StartQServers() {
	listener, err := net.Listen(ServerNetworkType, ServerPort)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	for {
		connect, err := listener.Accept()
		remoteIp := connect.RemoteAddr().String()
		//if !checkQTargetAuth(remoteIp) {
		//	err := connect.Close()
		//	if err != nil {
		//		continue
		//	}
		//}
		if err != nil {
			log.Fatalln(err)
		}
		common.CurrentWaitAuthIp = remoteIp
		if strings.Contains(remoteIp, "127.0.0.1") {
			connect.Close()
			AuthFlag = false
			continue
		}
		common.SendSignal2WGChannel(common.GetRunningSignal(common.TYPE_REMOTE_QNQ_AUTH))
		AuthLock.Lock()
		AuthLock.Lock()
		if AuthFlag {
			netCell := qNetCells[remoteIp]
			netCell.qServer = &connect
			qNetCells[remoteIp] = netCell
			go handleConnect(connect, false)
			log.Printf("qnq client connect, ip : %v ", remoteIp)
		} else {
			connect.Close()
		}
		AuthLock.Unlock()
		AuthFlag = false
	}
}

func StartQTargets() {
	for _, qnqConfig := range config.SystemConfigCache.Value().QNQNetCells {
		if qnqConfig.Ip == "0.0.0.0" && strings.Contains(qnqConfig.Ip, "127.0.0.1") {
			continue
		}
		if qNetCells[qnqConfig.Ip].qTarget != nil {
			continue
		}
		ConnectTarget(qnqConfig.Ip)
	}
}

// ConnectTarget 不允许非worker包调用network包
func ConnectTarget(ip string) {
	var err error
	log.Printf("try to connect remote qnq : %v\n", ip+ServerPort)
	connect, err := net.Dial(ServerNetworkType, ip+ServerPort)
	if err != nil {
		return
	}
	connect.SetReadDeadline(time.Now().Add(time.Second * time.Duration(30)))
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("client %v connected \n", connect.RemoteAddr())
		netCell := qNetCells[ip]
		netCell.qTarget = &connect
		qNetCells[ip] = netCell
		_, err := WriteStrToQTarget("test", ip)
		if err == nil {
			netCell.setTargetStatus(true)
			go handleConnect(connect, true)
			log.Printf("client start ...")
		} else {
			netCell.setTargetStatus(false)
		}
	}
}

func WriteStrToQTarget(message string, targetIp string) (string, error) {
	var err error
	var ret string
	_, err = write(*qNetCells[targetIp].qTarget, message)
	if err != nil {
		log.Printf(err.Error())
		return "", err
	}
	//ret, err = readStr(*qNetCells[targetIp].qTarget)
	//if err != nil {
	//	log.Printf(err.Error())
	//} else {
	//	log.Printf(ret)
	//}
	return ret, err
}

func ReadStrFromQTarget(targetIp string) string {
	ret, err := readStr(*qNetCells[targetIp].qTarget)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}
	return ret
}

func ReadBytesFromQTarget(targetIp string) bytes.Buffer {
	ret, err := readBytes(*qNetCells[targetIp].qTarget)
	if err != nil {
		log.Printf(err.Error())
		return bytes.Buffer{}
	}
	return ret
}

func GetQNetCell(ip string) QNetCell {
	return qNetCells[ip]
}
