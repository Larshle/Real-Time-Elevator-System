
// midlertidig main fil for å teste coden vår

package main

import (
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

	config.Init()
	elevio.Init("localhost:"+strconv.Itoa(config.Port), config.N_floors)

	deliveredOrderC := make(chan elevio.ButtonEvent, 10000)
	newElevStateC := make(chan elevator.State,10000)
	giverToNetwork := make(chan distributor.HRAInput,10000)
	receiveFromNetworkC := make(chan distributor.HRAInput,10000)
	messageToAssinger := make(chan distributor.HRAInput, 10000)
	eleveatorAssingmentC := make(chan elevator.Assingments,10000)
	lightsAssingmentC := make(chan elevator.Assingments,10000)
	chan_receiver_from_peers := make(chan peers.PeerUpdate,10000)
	chan_giver_to_peers := make(chan bool,10000)


	go peers.Receiver(config.RT_port_number, chan_receiver_from_peers)
	go peers.Transmitter(config.RT_port_number, config.Elevator_id, chan_giver_to_peers)

	go bcast.Receiver(config.RT_port_number+15, receiveFromNetworkC) // må endres
	go bcast.Transmitter(config.RT_port_number+15, giverToNetwork)

	go distributor.Distributor(
		deliveredOrderC,
		newElevStateC,
		giverToNetwork,
		receiveFromNetworkC,
		messageToAssinger, 
		chan_receiver_from_peers)

	go assigner.Assigner(
		eleveatorAssingmentC,
		lightsAssingmentC,
		messageToAssinger)

	go elevator.Elevator(
		eleveatorAssingmentC,
		newElevStateC,
		deliveredOrderC)

	go lights.Lights(lightsAssingmentC)

	select {} // for å kjøre alltid lol
}
