// midlertidig main fil for å teste coden vår 

package main

import (
	"root/config"
	"fmt"
	"root/distributor"
	"root/driver/elevio"
	"root/elevator"
	"root/assigner"
	"root/lights"
	"root/network/network_modules/peers"
	"root/network/network_modules/bcast"
)


func main() {

	fmt.Println("Hello, World!")
	fmt.Println("Elevator ID: ", config.Elevator_id)
	fmt.Println("N_floors: ", config.N_floors)
	fmt.Println("N_elevators: ", config.N_elevators)

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

	go peers.Receiver(15647, chan_receiver_from_peers)
	go peers.Transmitter(15647, config.Elevator_id, chan_giver_to_peers)

	
	go bcast.Receiver(15647, receiveFromNetworkC) // må endres
	go bcast.Transmitter(15647, giverToNetwork)


	go distributor.Distributor(
		deliveredOrderC,
		newElevStateC,
		giverToNetwork,
		receiveFromNetworkC,
		messageToAssinger)
	
	go assigner.Assigner(
		eleveatorAssingmentC,
		lightsAssingmentC,
		messageToAssinger)
	
	go elevator.Elevator(
		eleveatorAssingmentC,
		newElevStateC,
		deliveredOrderC)

	go lights.Lights(lightsAssingmentC)

	select {} // for å kjøre alltid lol lars er gey
}