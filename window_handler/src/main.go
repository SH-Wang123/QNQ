package main

import (
	"net"
	"window_handler/cli"
	"window_handler/common"
	"window_handler/gui"
	"window_handler/network"
	_ "window_handler/rest"
	"window_handler/worker"
)

func main() {
	worker.GetPartitionsInfo()
	//os.Setenv("FYNE_FONT", "msyh.ttc")
	//config.GetTargetSystemInfo()
	if common.CLI_FLAG {
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
