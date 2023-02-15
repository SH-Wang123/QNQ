package request

import (
	"io"
	"net/http"
)

func GetLastVersion() string {
	rsp, _ := http.Get("http://114.132.74.122:6600/version")
	if rsp == nil {
		return ""
	}
	defer CloseRspBody(rsp.Body)
	ret, err := io.ReadAll(rsp.Body)
	if err != nil {
		return ""
	}
	return string(ret)
}
