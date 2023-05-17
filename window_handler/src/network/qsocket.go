package network

import (
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
	"window_handler/common"
)

const (
	ServerNetworkType = "tcp"
	ServerPort        = ":9916"
	MessageDelimiter  = '\f'
	RecoverMessage    = "cx00000615"
)

var AuthLock = &sync.RWMutex{}
var AuthFlag = false

var qNetCells = make(map[string]*QNetCell, 8)

var NetChan = common.NewProducer()

func handleConnect(conn net.Conn, isClient bool) {
	log.Printf("client %v connected\n", conn.RemoteAddr())
	for {
		err := conn.SetReadDeadline(time.Now().Add(time.Minute * time.Duration(5)))
		if err != nil {
			continue
		}
		if msg, err := read(conn); err != nil {
			time.Sleep(100 * time.Millisecond)
			if err == io.EOF {
				//log.Printf("client %v closed\n", conn.RemoteAddr())
			} else {
				//log.Printf("read error : %v\n", err.Error())
			}
		} else {
			if !isClient {
				writeStr(conn, RecoverMessage)
			}
			if len(msg) > 0 {
				go NetChan.Produce(msg)
			}
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
		remoteIp := common.GetIpFromAddr(connect.RemoteAddr().String())
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
			var netCell *QNetCell
			if qNetCells[remoteIp] != nil {
				netCell = qNetCells[remoteIp]
			} else {
				netCell = &QNetCell{}
			}
			netCell.QServer = &connect
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

// ConnectTarget 不允许非worker包调用network包
func ConnectTarget(ip string) bool {
	if strings.Contains(ip, "0.0.0.0") || strings.Contains(ip, "127.0.0.1") {
		return false
	}
	var err error
	log.Printf("try to connect remote qnq : %v\n", ip+ServerPort)
	connect, err := net.Dial(ServerNetworkType, ip+ServerPort)
	if err != nil {
		return false
	}
	connect.SetReadDeadline(time.Now().Add(time.Second * time.Duration(30)))
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("client %v connected \n", connect.RemoteAddr())
		netCell := qNetCells[ip]
		if netCell == nil {
			netCell = &QNetCell{}
		}
		netCell.QTarget = &connect
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
	return true
}

func WriteStrToQTarget(message string, targetIp string) (string, error) {
	var err error
	var ret string
	_, err = writeStr(*qNetCells[targetIp].QTarget, message)
	if err != nil {
		log.Printf(err.Error())
		return "", err
	}
	return ret, err
}

func GetQNetCell(ip string) *QNetCell {
	cell := qNetCells[ip]
	if cell == nil {
		cell = &QNetCell{}
		qNetCells[ip] = cell
	}
	return cell
}

func GetAllQNetCells() map[string]*QNetCell {
	return qNetCells
}
