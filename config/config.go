package config

import (
	"flag"
	"fmt"
	"os"
	"root/network/network_modules/localip"
)

func Generate_ID() string {
	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}
	id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	return id
}

var Elevator_id = Generate_ID()

const(
	N_floors = 4
	N_elevators = 2
)
