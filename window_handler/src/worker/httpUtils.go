package worker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"window_handler/common"
)

func getObjFromResponse(resp *http.Response, obj any) any {
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	qresponse := common.QResponse{}
	err := json.Unmarshal(b, &qresponse)
	if err != nil {
		return nil
	}
	retJson, err := json.Marshal(qresponse.Data)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(retJson, obj)
	if err != nil {
		return nil
	}
	return obj
}

func sendGet(url string, params ...map[string]string) (resp *http.Response, err error) {
	var paramsStr = ""
	for k, v := range params {
		if paramsStr == "" {
			paramsStr = fmt.Sprintf("?%s=%s", k, v)
		} else {
			paramsStr = paramsStr + fmt.Sprintf("&%s=%s", k, v)
		}
	}
	return http.Get(url + paramsStr)
}
