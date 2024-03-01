// midlertidig main fil for å teste coden vår

package main

import (
	"flag"
	"fmt"
	"root/assigner"
	"root/config"
	"root/distributor"
	"root/driver/elevio"
	"root/elevator"
	"root/lights"
	"root/network/network_modules/bcast"
	"root/network/network_modules/peers"
	"strconv"
)

func main() {

	fmt.Println("Hello, World!")
	fmt.Println("Elevator ID: ", config.Elevator_id)
	fmt.Println("N_floors: ", config.N_floors)
	fmt.Println("N_elevators: ", config.N_elevators)

	// var port string
	// flag.StringVar(&port, "port", "", "port of this peer")
	// flag.Parse()

	// fmt.Println("Port: ", port)

	// elevio.Init(fmt.Sprintf("127.0.0.1:%v", port), config.N_floors)

	// var port string
	// flag.StringVar(&port, "--port", "", "port of this peer")
	// flag.Parse()

	// fmt.Println("Port: ", port)

	// elevio.Init(fmt.Sprintf("127.0.0.1:%v", port), config.N_floors)

	var id string
	flag.StringVar(&id, "peerID", "", "id of this peer")
	flag.Parse()

	idInt, _ := strconv.Atoi(id)
	port := 15360 + idInt
	addr := "localhost:" + fmt.Sprint(port)
	elevio.Init(addr, config.N_floors)

	// // Storing for powerloss, hentet fra vetle sin kode kan sees på
	// store, err := skv.Open(fmt.Sprintf("elev%v.db", Elevator_id))
	// if err != nil {
	// 	panic(err)
	// }

	// var cs central.CentralState
	// if err = store.Get("cs", &cs); err != nil && err != skv.ErrNotFound {
	// 	panic(err)
	// }
	// cs.Origin = Elevator_id

	deliveredOrderC := make(chan elevio.ButtonEvent)
	newElevStateC := make(chan elevator.State)
	giverToNetwork := make(chan distributor.HRAInput)
	receiveFromNetworkC := make(chan distributor.HRAInput)
	messageToAssinger := make(chan distributor.HRAInput)
	eleveatorAssingmentC := make(chan elevator.Assingments)
	lightsAssingmentC := make(chan elevator.Assingments)
	chan_receiver_from_peers := make(chan peers.PeerUpdate)
	chan_giver_to_peers := make(chan bool)

	fmt.Println("1")

	go peers.Receiver(15360 + idInt, chan_receiver_from_peers)
	go peers.Transmitter(15360 + idInt, config.Elevator_id, chan_giver_to_peers)

	fmt.Println("2")

	go bcast.Receiver(16568, receiveFromNetworkC) // må endres
	go bcast.Transmitter(16568, giverToNetwork)

	fmt.Println("3")

	go distributor.Distributor(
		deliveredOrderC,
		newElevStateC,
		giverToNetwork,
		receiveFromNetworkC,
		messageToAssinger)

	fmt.Println("4")

	go assigner.Assigner(
		eleveatorAssingmentC,
		lightsAssingmentC,
		messageToAssinger)

	fmt.Println("5")

	go elevator.Elevator(
		eleveatorAssingmentC,
		newElevStateC,
		deliveredOrderC)

	fmt.Println("6")

	go lights.Lights(lightsAssingmentC)

	fmt.Println("7")

	select {} // for å kjøre alltid lol lars er gey
}
