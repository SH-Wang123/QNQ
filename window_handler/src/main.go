/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	_ "window_handler/Qlog"
	"window_handler/cmd"
	"window_handler/common"
	"window_handler/gui"
)

func main() {
	//os.Setenv("FYNE_FONT", "msyh.ttc")
	//config.GetTargetSystemInfo()
	if common.CLI_FLAG {
		cmd.Execute()
	} else {
		gui.StartGUI()
	}
}
