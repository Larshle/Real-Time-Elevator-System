package distributor

import (
	"fmt"
	"root/assigner"
	"root/elevator/localElevator"
	"root/network/network_modules/peers"
)

var N_floors = 4

var Commonstate = assigner.HRAInput{
	HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
	States: map[string]assigner.HRAElevState{
		"one":{
			Behavior:    "moving",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{false, false, false, true},
		},
		"two":{
			Behavior:    "idle",
			Floor:       3,
			Direction:   "stop",
			CabRequests: []bool{true, false, false, false},
		},
	},
}

func update_Commonstate(local_elevator_state Elevator) {

	elevator_id := "one"

	var new_HallRequests [][]bool
	var new_CabRequests  []bool

	new_HallRequests = make([][]bool, N_floors)
    for i := range new_HallRequests {
        new_HallRequests[i] = make([]bool, 2) // Initialize the inner slices
    }
    new_CabRequests = make([]bool, N_floors)



	// Iterate through each "row" of the 2D slice
	for i, row := range local_elevator_state.Assignments {
		new_HallRequests[len(local_elevator_state.Assignments)-1-i][0] = row[0]
		new_HallRequests[len(local_elevator_state.Assignments)-1-i][1] = row[1]
		new_CabRequests[len(local_elevator_state.Assignments)-1-i]     = row[2]
	}

	Commonstate.States[elevator_id] = assigner.HRAElevState{
		Behavior:    local_elevator_state.Behavior,
		Floor:       local_elevator_state.CurrentFloor,
		Direction:   local_elevator_state.Direction,
		CabRequests: []bool{true, true, false, false},
	}

}