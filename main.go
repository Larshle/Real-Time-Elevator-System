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

	port := flag.Int("port", 15357, "<-- Default verdi, men kan overskrives som en command line argument ved bruk av -port=xxxxx")
	flag.Parse()
	fmt.Printf("Port: %d\n", *port)

	elevio.Init("localhost:" + strconv.Itoa(*port), config.N_floors)

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

	go peers.Receiver(config.RT_port_number, chan_receiver_from_peers)
	go peers.Transmitter(config.RT_port_number, config.Elevator_id, chan_giver_to_peers)

	fmt.Println("2")

	go bcast.Receiver(config.RT_port_number, receiveFromNetworkC) // må endres
	go bcast.Transmitter(config.RT_port_number, giverToNetwork)

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
