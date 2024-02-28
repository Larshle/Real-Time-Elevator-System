// midlertidig main fil for å teste coden vår 

package main

import (
	"fmt"
	"root/network"
	"root/distributor"
	"root/driver/elevio"
	"root/elevator"
	"root/assigner"
	"root/lights"
)


func main() {

	N_floors := 4
	N_elevators := 3
	var Elevator_id = network.Generate_ID()

	fmt.Println(Elevator_id)
	fmt.Println("Hello, World!")
	fmt.Println("N_floors: ", N_floors)
	fmt.Println("N_elevators: ", N_elevators)

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

	go distributor.Distributor(
		deliveredOrderC,
		newElevStateC,
		giverToNetwork,
		receiveFromNetworkC,
		messageToAssinger,
		Elevator_id)
	
	go assigner.Assigner(
		eleveatorAssingmentC,
		lightsAssingmentC,
		messageToAssinger,
		Elevator_id)
	
	go elevator.Elevator(
		eleveatorAssingmentC,
		newElevStateC,
		deliveredOrderC)

	go lights.Lights(lightsAssingmentC)

	select {} // for å kjøre alltid lol lars er gey
}