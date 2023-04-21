/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"net"
	"window_handler/cmd"
	"window_handler/common"
	"window_handler/gui"
	"window_handler/network"
)

func main() {
	//os.Setenv("FYNE_FONT", "msyh.ttc")
	//config.GetTargetSystemInfo()
	if common.CLI_FLAG {
		cmd.Execute()
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
