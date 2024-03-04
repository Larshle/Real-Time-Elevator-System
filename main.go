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
)

func main() {

	fmt.Println("Hello, World!")
	fmt.Println("Elevator ID: ", config.Elevator_id)
	fmt.Println("N_floors: ", config.N_floors)
	fmt.Println("N_elevators: ", config.N_elevators)

	port := flag.String("port", "15360", "Default verdi er 15360, men kan overskrives som et command line argument")
	flag.Parse()
	fmt.Printf("Port: %s\n", *port)
	elevio.Init("localhost:" + *port, config.N_floors)

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

	go peers.Receiver(15357, chan_receiver_from_peers)
	go peers.Transmitter(15357, config.Elevator_id, chan_giver_to_peers)

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
