package network

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"time"
)

var icmp ICMP

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func readStr(conn net.Conn) (string, error) {
	var str string
	var bytebuf bytes.Buffer
	bytearr := make([]byte, 1)
	for {
		if _, err := conn.Read(bytearr); err != nil {
			return str, err
		}
		item := bytearr[0]
		if item == MessageDelimiter {
			break
		}
		bytebuf.WriteByte(item)
	}
	str = bytebuf.String()
	return str, nil
}

func readBytes(conn net.Conn) (bytes.Buffer, error) {
	var bytebuf bytes.Buffer
	var err error
	bytearr := make([]byte, 1)
	for {
		if _, err := conn.Read(bytearr); err != nil {
			return bytebuf, err
		}
		item := bytearr[0]
		if item == MessageDelimiter {
			break
		}
		if item == EndMessage {
			return bytebuf, err
		}
		bytebuf.WriteByte(item)
	}
	return bytebuf, err
}

func write(conn net.Conn, content string) (int, error) {
	log.Printf("send %v : %v\n", conn.RemoteAddr(), content)
	var bytebuf bytes.Buffer
	bytebuf.WriteString(content)
	bytebuf.WriteByte(MessageDelimiter)
	bytearr := bytebuf.Bytes()
	return conn.Write(bytearr)
}

func loadContentWithCheck(prefix string, msg string, checkBit string) string {
	return prefix + msg + checkBit
}

func LoadContent(prefix string, msg string) string {
	return loadContentWithCheck(prefix, msg, "00000000")
}

func TestPing(ip string) bool {
	icmp.Type = 8 // 8->echo message  0->reply message
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequenceNum = 0

	recvBuf := make([]byte, 32)
	var buffer bytes.Buffer
	//先在buffer中写入icmp数据报求去校验和
	binary.Write(&buffer, binary.BigEndian, icmp)
	icmp.Checksum = CheckSum(buffer.Bytes())
	//然后清空buffer并把求完校验和的icmp数据报写入其中准备发送
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)

	Time, _ := time.ParseDuration("2s")
	conn, err := net.DialTimeout("ip4:icmp", ip, Time)
	if err != nil {
		return false
	}
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		log.Println("conn.Write error:", err)
		return false
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	num, err := conn.Read(recvBuf)
	if err != nil {
		log.Println("conn.Read error:", err)
		return false
	}

	conn.SetReadDeadline(time.Time{})

	if string(recvBuf[0:num]) != "" {
		return true
	}
	return false

}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}
