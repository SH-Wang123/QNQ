package network

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
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

var qNetCells = make(map[string]qNetCell, 8)

var NetChan = common.NewProducer()

func handleConnect(conn net.Conn, isClient bool) {
	log.Printf("client %v connected\n", conn.RemoteAddr())
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(30)))
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
		if !checkQTargetAuth(remoteIp) {
			err := connect.Close()
			if err != nil {
				return
			}
		}
		if err != nil {
			log.Fatalln(err)
		}
		netCell := qNetCells[remoteIp]
		netCell.qServer = &connect
		go handleConnect(connect, false)
		log.Printf("qnq client connect, ip : %v ", remoteIp)
	}
}

func StartQTargets() {
	if config.SystemConfigCache.Value().QnqSTarget.Ip == "0.0.0.0" {
		log.Printf("not set QNQ Target")
		return
	}
	targetIps := strings.Split(config.SystemConfigCache.Value().QnqSTarget.Ip, ",")
	for _, ip := range targetIps {
		var err error
		log.Printf("client connect %v\n", config.SystemConfigCache.Value().QnqSTarget.Ip+ServerPort)
		connect, err := net.Dial(ServerNetworkType, ip+ServerPort)
		if err != nil {
			return
		}
		connect.SetReadDeadline(time.Now().Add(time.Second * time.Duration(30)))
		if err != nil {
			log.Printf(err.Error())
		} else {
			log.Printf("client %v connected \n", connect.RemoteAddr())
			ret, err := WriteStrToQTarget("test", ip)
			netCell := qNetCells[ip]
			if err == nil && ret != "" {
				netCell.setTargetStatus(true)
				go handleConnect(connect, true)
				log.Printf("client start ...")
				netCell.qTarget = &connect
			} else {
				netCell.setTargetStatus(false)
			}
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
	//if ret, err := readStr(ConnectClient); err != nil {
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
