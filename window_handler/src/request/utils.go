package request

import (
	"io"
	"log"
)

func CloseRspBody(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		log.Printf("close body err : %v", err)
	}
}
