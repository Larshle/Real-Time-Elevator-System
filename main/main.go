// midlertidig main fil for å teste coden vår 

package main

import (
	"fmt"
	"root/network"
)


func main() {

	N_floors := 4
	N_elevators := 3
	var id = network.Generate_ID()


	fmt.Println(id)
	fmt.Println("Hello, World!")
	fmt.Println("N_floors: ", N_floors)
	fmt.Println("N_elevators: ", N_elevators)

	// Hva må skje her? Prøver å lage en liste
	// Init - alle heisene må sende staten sin til distributor som oppdaterer commonstate

}