package main

import (
	"QNQ"
	_ "QNQ"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

const logPath = "QNQ_W.log"

func initLog() {
	_, err := os.Stat(logPath)
	if err != nil {
		filePtr, _ := os.Create(logPath)
		defer func() {
			err := filePtr.Close()
			if err != nil {
				return
			}
		}()
	}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}

	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	slog.Info("Welcome to QNQ.", "version", QNQ.Version)
}

func main() {
	initLog()
	QNQ.Start()
	runtime.Gosched()
	terminal()
}

func clearResource() {

}

func terminal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	select {
	case <-c:
		slog.Info("SIGTERM signal caught")
	}
	clearResource()
}
