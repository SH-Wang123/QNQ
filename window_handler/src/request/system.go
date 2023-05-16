package request

import (
	"io"
	"net/http"
	"regexp"
)

func GetLastVersion() string {
	rsp, err := http.Get("http://114.132.74.122:6600/version")
	if rsp == nil || err != nil {
		return ""
	}
	defer CloseRspBody(rsp.Body)
	ret, err := io.ReadAll(rsp.Body)
	reg := regexp.MustCompile("^V\\d+.\\d+.\\d+")
	b := reg.Match(ret)
	if err != nil || !b {
		return ""
	}

	return string(ret)
}
