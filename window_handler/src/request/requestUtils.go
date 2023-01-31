package request

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"window_handler/data"
)

var logger = log.Default()

func GetRequest(url string) *data.ApiResponse {
	resp, err := http.Get(url)
	if err != nil {
		logger.Printf("error send to  : ", url, err)
	}
	defer resp.Body.Close()
	return getApiResponse(resp, url)
}

func PostRequest(url string, params map[string]string) *data.ApiResponse {
	dataType, _ := json.Marshal(params)
	dataString := string(dataType)
	resp, err := http.Post(url,
		"text/json;charset=utf-8",
		strings.NewReader(dataString))
	if err != nil {
		logger.Printf("error ModifySysConfig", err)
	}
	defer resp.Body.Close()
	return getApiResponse(resp, url)
}

func getApiResponse(resp *http.Response, url string) *data.ApiResponse {
	body, err := ioutil.ReadAll(resp.Body)
	ret := new(data.ApiResponse)
	err = json.Unmarshal(body, &ret)
	if err != nil {
		logger.Printf("tran body err : ", url, err)
		return nil
	}
	return ret
}
