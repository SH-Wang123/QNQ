package Qlog

import (
	"log"
)

func outputStatus(name string, id string, status string) {
	log.Printf("%v begin to %v, id(SN) is {%v}", status, name, id)
}

func OutputStart(name string, id string) {
	outputStatus(name, id, "start")
}

func OutputOver(name string, id string) {
	outputStatus(name, id, "over")
}

func OutputFailedStatus(name string, id string) {
	outputStatus(name, id, "over")
}
