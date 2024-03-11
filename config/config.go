package config

import (
	//"fmt"
	// "os"
	// "root/network/network_modules/localip"
	"strconv"
	"flag"
)

// func Generate_ID() string {
// 	localIP, err := localip.LocalIP()
// 	if err != nil {
// 		fmt.Println(err)
// 		localIP = "DISCONNECTED"
// 	}
// 	id := fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
// 	return id
// }

const (
	N_floors       = 4
	N_elevators    = 2
	RT_port_number = 58735
)

var Port int
var Elevator_id string

func Init() {
	//fmt.Println("Init")
	//fmt.Println("N_floors: ", N_floors)
	//fmt.Println("N_elevators: ", N_elevators)
	port := flag.Int("port", 15357, "<-- Default verdi, men kan overskrives som en command line argument ved bruk av -port=xxxxx")
	id := flag.Int("id", 10000, "id")
	flag.Parse()
	//fmt.Printf("Port: %d\n", *port)
	//fmt.Printf("id: %d\n", *id)

	Port = *port
	Elevator_id = "peer-10.22.229.227-" + strconv.Itoa(*id)
	//fmt.Printf("Elevator_id: %s\n", Elevator_id)
}