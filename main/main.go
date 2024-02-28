// midlertidig main fil for å teste coden vår 

package main

import (
	"fmt"
	"root/network"
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

	// Hva må skje her? Prøver å lage en liste
	// Init - alle heisene må sende staten sin til distributor som oppdaterer commonstate

	

}