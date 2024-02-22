package distributor

import (
	"fmt"
	"root/assigner"
	"root/elevator"
	"root/network/network_modules/peers"
	"root/driver/elevio"
)

var N_floors = 4

var Commonstate = assigner.HRAInput{
	Origin: "string",
	ID: 1,
	Ackmap: map[string]string{},
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

var Unacked_Commonstate = assigner.HRAInput{
	Origin: "string",
	ID: 1,
	Ackmap: map[string]string{},
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
	fmt.Println("\nOrigin:", cs.Origin)
	fmt.Println("ID:", cs.ID)
	fmt.Println("Ackmap:", cs.Ackmap)
	fmt.Println("Hall Requests:", cs.HallRequests)

	for i, state := range cs.States {
		fmt.Printf("Elevator %s:\n", string(i))
		fmt.Printf("\tBehaviour: %s\n", state.Behaviour)
		fmt.Printf("\tFloor: %d\n", state.Floor)
		fmt.Printf("\tDirection: %s\n", state.Direction)
		fmt.Printf("\tCab Requests: %v\n\n", state.CabRequests)
	}
}

func Update_Commonstate(local_elevator_state elevator.elevator) {

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


func Commonstates_are_equal(new_commonstate, Commonstate assigner.HRAInput) bool {
    // Compare Ackmaps
    // for k, v := range a.Ackmap {
    //     if b.Ackmap[k] != v {
    //         return false
    //     }
    // }
	
	if new_commonstate.ID != Commonstate.ID {
		return false
	}
    // Compare HallRequests
    if len(new_commonstate.HallRequests) != len(Commonstate.HallRequests) {
        return false
    }
    for i, v := range new_commonstate.HallRequests {
        if Commonstate.HallRequests[i] != v {
            return false
        }
    }

    // Compare States
    for k, v := range new_commonstate.States {
		bv, ok := Commonstate.States[k]
		if !ok {
			return false
		}
		if bv.Behaviour != v.Behaviour || bv.Floor != v.Floor || bv.Direction != v.Direction {
			return false
		}
		if len(bv.CabRequests) != len(v.CabRequests) {
			return false
		}
		for i, cabRequest := range bv.CabRequests {
			if cabRequest != v.CabRequests[i] {
				return false
			}
		}
	}

	return true
}

func Recieve_commonstate(new_commonstate assigner.HRAInput) {
	if Commonstates_are_equal(new_commonstate, Unacked_Commonstate) {
		return
	}
	// if fullack (skriv dette senere)
	// Commonstate = new_commonstate 
	// broadcast
	// kjør til assigner

	// if new_commonstate har lavere prioritet
	// return
	if new_commonstate.ID < Unacked_Commonstate.ID {
		return
	}

	// else
	// ack, oppdater ack_commonstate og broadcast denne helt til den er acket eller det kommer en ny med høyere prioritet

}