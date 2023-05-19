package Qlog

import (
	"io"
	"log"
	"os"
)

const logPath = "QNQ_W.log"

func init() {
	_, err := os.Stat(logPath)
	if err != nil {
		filePtr, _ := os.Create(logPath)
		defer func() {
			filePtr.Close()
		}()
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}

	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Welcome to QNQ.")

}
