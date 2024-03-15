package main

import (
	"flag"
	"fmt"
	"root/assigner"
	"root/config"
	"root/distributor"
	"root/elevator"
	"root/elevio"
	"root/lights"
	"root/network/bcast"
	"root/network/peers"
	"strconv"
)

var Port int
var ElevatorID int

func main() {

	port := flag.Int("port", 15301, "<-- Default verdi, men kan overskrives som en command line argument ved bruk av -port=xxxxx")
	id := flag.Int("id", 0, "<-- Default verdi, men kan overskrives som en command line argument ved bruk av -id=xxxxx")
	flag.Parse()

	Port = *port
	ElevatorID = *id

	fmt.Println()
	elevio.Init("localhost:"+strconv.Itoa(Port), config.NumFloors)

	fmt.Println("Elevator initialized with ID", ElevatorID, "on port", Port)
	fmt.Println("System has", config.NumFloors, "floors and", config.NumElevators, "elevators.")

	newAssignmentC := make(chan elevator.Assignments, 10000)
	deliveredAssignmentC := make(chan elevio.ButtonEvent, 10000)
	newLocalStateC := make(chan elevator.State, 10000)
	networkTx := make(chan distributor.CommonState, 10000) // Endre navn?
	networkRx := make(chan distributor.CommonState, 10000) // Endre navn? ja faktisk
	confirmedCommonstateC := make(chan distributor.CommonState, 10000)
	peersRx := make(chan peers.PeerUpdate, 10000) // Endre navn?
	peersTx := make(chan bool, 10000)

	go peers.Receiver(config.PeersPortNumber, peersRx)
	go peers.Transmitter(config.PeersPortNumber, ElevatorID, peersTx)

	go bcast.Receiver(config.BcastPortNumber, networkRx)
	go bcast.Transmitter(config.BcastPortNumber, networkTx)

	go distributor.Distributor(
		deliveredAssignmentC,
		newLocalStateC,
		networkTx,
		networkRx,
		confirmedCommonstateC,
		peersRx,
		ElevatorID)

	go elevator.Elevator(
		newAssignmentC,
		deliveredAssignmentC,
		newLocalStateC)

	for {
		select {
		case cs := <-confirmedCommonstateC:
			localAssingment := assigner.CalculateOptimalAssignments(cs, ElevatorID)
			newAssignmentC <- localAssingment
			lights.SetLights(cs, ElevatorID)

		default:
			continue
		}

	}
}
