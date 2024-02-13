package distributor

import (
	"fmt"
	"root/assigner"
	"root/elevator/localElevator"
	"root/network/network_modules/peers"
	"root/driver/elevio"
)

var N_floors = 4

var Commonstate = assigner.HRAInput{
	HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
	States: map[string]assigner.HRAElevState{
		"one":{
			Behaviour:   "moving",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{false, false, false, true},
		},
		"two":{
			Behaviour:   "idle",
			Floor:       3,
			Direction:   "stop",
			CabRequests: []bool{true, false, false, false},
		},
	},
}

// Map for å endre fra type til string
var motorDirectionMap = map[elevio.MotorDirection]string{
	elevio.MD_Up:   "up",
	elevio.MD_Down: "down",
	elevio.MD_Stop: "stop",
}

func printCommonState(cs assigner.HRAInput) {
	fmt.Println("\nHall Requests:")
	fmt.Println(cs.HallRequests)

	for i, state := range cs.States {
		fmt.Printf("\nElevator %s:\n", string(i))
		fmt.Printf("\tBehaviour: %s\n", state.Behaviour)
		fmt.Printf("\tFloor: %d\n", state.Floor)
		fmt.Printf("\tDirection: %s\n", state.Direction)
		fmt.Printf("\tCab Requests: %v\n\n", state.CabRequests)
	}
}

func Update_Commonstate(local_elevator_state Elevator) {

	// skal bytte dette ut med unik id
	elevator_id := "one"

	var new_HallRequests [][2]bool
	var new_CabRequests []bool

	new_HallRequests = make([][2]bool, N_floors)
	new_CabRequests = make([]bool, N_floors)

	// Endrer format fra Elevator til Commonstate
	for i, row := range local_elevator_state.Assignments {
		new_HallRequests[len(local_elevator_state.Assignments)-1-i][0] = row[0]
		new_HallRequests[len(local_elevator_state.Assignments)-1-i][1] = row[1]
		new_CabRequests[len(local_elevator_state.Assignments)-1-i] = row[2]
	}

	// Oppdaterer Commonstate
	Commonstate.HallRequests = new_HallRequests
	Commonstate.States[elevator_id] = assigner.HRAElevState{
		Behaviour:   string(local_elevator_state.Behaviour),
		Floor:       local_elevator_state.CurrentFloor,
		Direction:   motorDirectionMap[local_elevator_state.Direction],
		CabRequests: new_CabRequests,
	}
}