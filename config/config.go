package config

import (
	"fmt"
	"os"
	"root/network/network_modules/localip"
)

func Generate_ID() string {
	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id := fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	return id
}

var Elevator_id = Generate_ID()

const(
	N_floors = 4
	N_elevators = 1
	RT_port_number = 58735
)
