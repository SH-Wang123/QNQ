package network

import (
	"bytes"
	"io"
	"log"
	"net"
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

var ConnectClient net.Conn
var ConnectServer net.Conn
var ConnectStauts = false
var currentTask = 0x00000001

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

func StartQServer() {
	log.Printf("start new server")
	listener, err := net.Listen(ServerNetworkType, ServerPort)
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	log.Printf("waiting client connect... ")
	for {
		ConnectServer, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		go handleConnect(ConnectServer, false)
	}
}

func StartQClient() {
	if config.SystemConfigCache.Value().QnqTarget.LocalPath == "0.0.0.0" {
		log.Printf("not set QNQ Target")
		return
	}
	var err error
	log.Printf("client connect %v\n", config.SystemConfigCache.Value().QnqTarget.Ip+ServerPort)
	ConnectClient, err = net.Dial(ServerNetworkType, config.SystemConfigCache.Value().QnqTarget.Ip+ServerPort)
	if err != nil {
		return
	}
	ConnectClient.SetReadDeadline(time.Now().Add(time.Second * time.Duration(30)))
	if err != nil {
		log.Printf(err.Error())
	} else {
		log.Printf("client %v connected \n", ConnectClient.RemoteAddr())
		ret, err := WriteStrToQTarget("test")
		if err == nil && ret != "" {
			ConnectStauts = true
			go handleConnect(ConnectClient, true)
			log.Printf("client start ...")
		} else {
			ConnectStauts = false
		}
	}
}

func WriteStrToQTarget(message string) (string, error) {
	var err error
	var ret string
	_, err = write(ConnectClient, message)
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

func ReadStrFromQTarget() string {
	ret, err := readStr(ConnectClient)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}
	return ret
}

func ReadBytesFromQTarget() bytes.Buffer {
	ret, err := readBytes(ConnectClient)
	if err != nil {
		log.Printf(err.Error())
		return bytes.Buffer{}
	}
	return ret
}
