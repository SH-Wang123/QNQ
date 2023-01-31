package Qlog

import (
	"io"
	"log"
	"os"
)

const LOG_PATH = "./QNQ_W.log"

func MakeLogger() {
	_, err := os.Stat(LOG_PATH)
	if err != nil {
		filePtr, _ := os.Create(LOG_PATH)
		defer func() {
			filePtr.Close()
		}()
	}

	f, err := os.OpenFile(LOG_PATH, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}

	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("test")

}
