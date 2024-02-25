package distributor

import (
	"fmt"
	"root/assigner"
	"root/elevator"
	"root/network"
	"root/network/network_modules/peers"
	"root/driver/elevio"
	"root/main"
	"strconv"
	"strings"
	"bytes"
	"net"
)

var Elevator_id = network.Generate_ID()
var N_floors = 4


// func (a localAssignments) Add_Assingment(newAssignments elevio.ButtonEvent) localAssignments{
// 	if newAssignments.Button == elevio.BT_Cab {
// 		a.localCabAssignments[newAssignments.Floor] = true
// 	} else {
// 		a.localHallAssignments[newAssignments.Floor][newAssignments.Button] = true
// 	}
// 	return a
// }

// func (a localAssignments) Remove_Assingment( deliveredAssingement elevio.ButtonEvent) localAssignments{
// 	if deliveredAssingement.Button == elevio.BT_Cab {
// 		a.localCabAssignments[deliveredAssingement.Floor] = false
// 	} else {
// 		a.localHallAssignments[deliveredAssingement.Floor][deliveredAssingement.Button] = false
// 	}
// 	return a
// }

var Commonstate = assigner.HRAInput{
	Origin: "string",
	ID: 1,
	Ackmap: make(map[string]assigner.Ack_status),
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
	Ackmap: make(map[string]assigner.Ack_status),
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
var DirectionMap = map[elevator.Direction]string{
	elevator.Down:   "down",
	elevator.Up: "up",
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

func Update_assignments(local_elevator_assignments elevator.Assingments) {

	var new_HallRequests [][2]bool
	var new_CabRequests []bool

	new_HallRequests = make([][2]bool, N_floors)
	new_CabRequests = make([]bool, N_floors)

	// Endrer format fra Assignments til Commonstate
	for i, row := range local_elevator_assignments {
		new_HallRequests[len(local_elevator_assignments)-1-i][0] = row[0]
		new_HallRequests[len(local_elevator_assignments)-1-i][1] = row[1]
		new_CabRequests[len(local_elevator_assignments)-1-i] = row[2]
	}
	// Oppdaterer hall requests
	Unacked_Commonstate.HallRequests = new_HallRequests

	// Oppdaterer cab requests
	temp_state := Unacked_Commonstate.States[Elevator_id]
	temp_state.CabRequests = new_CabRequests
	Unacked_Commonstate.States[Elevator_id] = temp_state
}

func Update_local_state(local_elevator_state elevator.State) {

	// Create a temporary variable to hold the updated state
	tempState := Unacked_Commonstate.States[Elevator_id]
	tempState.Behaviour = string(local_elevator_state.Behaviour)

	// Assign the updated state back to the map
	Unacked_Commonstate.States[Elevator_id] = tempState
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

func Fully_acked(ackmap map[string]assigner.Ack_status) bool {
	for id, value := range ackmap {
		if value == 0 && id != Elevator_id {
			return false
		}
	}
	return true
}	

func id_is_lower(id1, id2 string) bool {
	// Parse IP addresses and process IDs
	parts1 := strings.Split(id1, "-")
	parts2 := strings.Split(id2, "-")
	ip1 := net.ParseIP(parts1[1])
	ip2 := net.ParseIP(parts2[1])
	pid1, _ := strconv.Atoi(parts1[2])
	pid2, _ := strconv.Atoi(parts2[2])

	// Compare IP addresses
	cmp := bytes.Compare(ip1, ip2)
	if cmp < 0 {
		return true
	} else if cmp > 0 {
		return false
	}

	// If IP addresses are equal, compare process IDs
	return pid1 < pid2
}


func Recieve_commonstate(new_commonstate assigner.HRAInput) {
	
	if Commonstates_are_equal(new_commonstate, Unacked_Commonstate) {
		return
	}

	if Fully_acked(new_commonstate.Ackmap) {
		Unacked_Commonstate = new_commonstate // vet ikke om dette er nødvendig
		Commonstate = new_commonstate
		// broadcast
		// kjør til assigner
	}

	// if new_commonstate har lavere prioritet
	// return
	if new_commonstate.ID < Unacked_Commonstate.ID || id_is_lower(new_commonstate.Origin, Unacked_Commonstate.Origin) {
		return
	}

	// else
	// ack, oppdater ack_commonstate og broadcast denne helt til den er acket eller det kommer en ny med høyere prioritet

}


