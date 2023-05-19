package common

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	ServerNetworkType = "tcp"
	ServerPort        = ":9916"
	RecoverMessage    = "cx00000615"
)

var AuthLock = &sync.RWMutex{}
var AuthFlag = false

var qNetCells = make(map[string]*QNetCell, 8)

var NetChan = NewProducer()

func handleConnect(conn net.Conn, isTarget bool) {
	log.Printf("client %v connected\n", conn.RemoteAddr())
	remoteIp := GetIpFromAddr(conn.RemoteAddr().String())
	var buf [65542]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := conn.Read(buf[0:])
		result.Write(buf[0:n])
		if err != nil {
			time.Sleep(1 * time.Millisecond)
			if isTarget {
				continue
			}
			qNetCells[remoteIp].setTargetStatus(false)
			qNetCells[remoteIp] = nil
			conn.Close()
			return
		} else {
			scanner := bufio.NewScanner(result)
			scanner.Split(packetSlitFunc)
			for scanner.Scan() {
				handlerRecMessage(string(scanner.Bytes()[6:]), remoteIp)
				if len(string(scanner.Bytes()[6:])) > 0 {
					NetChan.Produce(string(scanner.Bytes()[6:]))
					break
				}
			}
		}
		result.Reset()
	}
}

func handlerRecMessage(msg string, remoteIp string) {
	if msg == RecoverMessage {
		qSender := GetCellQSender(remoteIp)
		if qSender != nil {
			qSender.RecCount++
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
		remoteIp := GetIpFromAddr(connect.RemoteAddr().String())
		//if !checkQTargetAuth(remoteIp) {
		//	err := connect.Close()
		//	if err != nil {
		//		continue
		//	}
		//}
		if err != nil {
			log.Println(err)
		}
		CurrentWaitAuthIp = remoteIp
		if strings.Contains(remoteIp, "127.0.0.1") {
			connect.Close()
			AuthFlag = false
			continue
		}
		SendSignal2WGChannel(GetRunningSignal(TYPE_REMOTE_QNQ_AUTH))
		AuthLock.Lock()
		AuthLock.Lock()
		if AuthFlag {
			var netCell *QNetCell
			if qNetCells[remoteIp] != nil {
				netCell = qNetCells[remoteIp]
			} else {
				netCell = &QNetCell{
					netCellLock: &sync.RWMutex{},
				}
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
	log.Printf("try to connect remote qnq : %v\n", ip)
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
			netCell = &QNetCell{
				netCellLock: &sync.RWMutex{},
			}
		}
		netCell.QTarget = &connect
		qNetCells[ip] = netCell
		_, err := WriteStrToQTarget([]byte("test"), ip, nil)
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

func DisconnectTarget(ip string) {
	log.Printf("try to disconnect remote qnq : %v\n", ip)
	if qNetCells[ip] == nil || qNetCells[ip].QTarget == nil {
		return
	}
	ptr := qNetCells[ip].QTarget
	target := *ptr
	err := target.Close()
	if err != nil {
		log.Printf("disconnect err, ip: %v, err : %v", ip, err)
		return
	}
}

// WriteStrToQTarget 如果不由QSender发送，参数传nil
func WriteStrToQTarget(message []byte, targetIp string, s *QSender) (string, error) {
	if s != nil && qNetCells[targetIp] != nil && qNetCells[targetIp].currentTWorker != s {
		log.Printf("%v want stael conn", s.SN)
		return "", nil
	}
	var err error
	var ret string
	err = writeToConn(qNetCells[targetIp].QTarget, message)
	if err != nil {
		log.Printf(err.Error())
		return "", err
	}
	return ret, err
}

func GetQNetCell(ip string) *QNetCell {
	cell := qNetCells[ip]
	if cell == nil {
		cell = &QNetCell{
			netCellLock: &sync.RWMutex{},
		}
		qNetCells[ip] = cell
	}
	return cell
}

func GetAllQNetCells() map[string]*QNetCell {
	return qNetCells
}
