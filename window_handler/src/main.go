package main

import (
	"net"
	"window_handler/Qlog"
	"window_handler/cli"
	"window_handler/common"
	"window_handler/config"
	"window_handler/gui"
	"window_handler/network"
	_ "window_handler/rest"
	"window_handler/worker"
)

func main() {
	worker.GetPartitionsInfo()
	//os.Setenv("FYNE_FONT", "msyh.ttc")
	Qlog.MakeLogger()
	common.InitCoroutinesPool()
	worker.LoadWorkerFactory()
	go network.StartQServer()
	go network.StartQClient()
	network.NetChan.StartPump()
	//config.GetTargetSystemInfo()
	worker.InitFileNode(false, true)
	common.GetCoroutinesPool().StartPool()
	if config.CLI_FALG {
		cli.StartCli()
	} else {
		gui.StartGUI()
	}
	if network.ConnectStauts {
		defer func(ConnectClient net.Conn) {
			err := ConnectClient.Close()
			if err != nil {

			}
		}(network.ConnectClient)
	}
}
