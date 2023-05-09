package network

import "net"

type qNetCell struct {
	qTarget *net.Conn
	qServer *net.Conn
	//target status | server status
	status int
}

func (qn *qNetCell) setTargetStatus(status bool) {
	if status {
		qn.status = qn.status | 10
	} else {
		qn.status = qn.status & 01
	}
}

func (qn *qNetCell) setServerStatus(status bool) {
	if status {
		qn.status = qn.status | 01
	} else {
		qn.status = qn.status & 10
	}
}

func (qn *qNetCell) getTargetStatus() bool {
	return qn.status&10 >= 10
}

func (qn *qNetCell) getServerStatus() bool {
	return qn.status&01 == 1
}

var idCount = 0

var qnqAuthList = make([]string, 4)

func checkQTargetAuth(ip string) bool {
	for _, v := range qnqAuthList {
		if v == ip {
			return true
		}
	}
	return false
}

func addAuth(ip string) {
	qnqAuthList = append(qnqAuthList, ip)
}
