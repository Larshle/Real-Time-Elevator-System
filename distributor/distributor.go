package distributor

import (
	"fmt"
	"root/elevator"
	"strconv"
	"strings"
	"bytes"
	"net"
)

var N_floors = 4

type Ack_status int
const (
	NotAcked Ack_status = iota
	Acked
	NotAvailable
)

type HRAElevState struct {
	Behaviour   string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	ID int
	Origin string
	Ackmap map[string]Ack_status
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

var Commonstate = HRAInput{
	Origin: "string",
	ID: 1,
	Ackmap: make(map[string]Ack_status),
	HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
	States: map[string]HRAElevState{
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

var Unacked_Commonstate = HRAInput{
	Origin: "string",
	ID: 1,
	Ackmap: make(map[string]Ack_status),
	HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
	States: map[string]HRAElevState{
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

func (es HRAElevState) toHRAElevState(localElevState elevator.State) {
	es.Behaviour = string(localElevState.Behaviour)
	es.Floor = localElevState.Floor
	es.Direction = string(localElevState.Direction)
}


func printCommonState(cs HRAInput) {
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

func (cs HRAInput) Update_Assingments(local_elevator_assignments localAssignments) HRAInput {

	for f := 0; f < N_floors; f++ {
		for b := 0; b < 2; b++ {
			if local_elevator_assignments.localHallAssignments[f][b] == add {
				cs.HallRequests[f][b] = true
			}
			if local_elevator_assignments.localHallAssignments[f][b] == remove {
				cs.HallRequests[f][b] = false
			}
		}
	}

	for f := 0; f < N_floors; f++ {
		if local_elevator_assignments.localCabAssignments[f] == add {
			cs.States[Elevator_id].CabRequests[f] = true
		}
		if local_elevator_assignments.localCabAssignments[f] == remove {
			cs.States[Elevator_id].CabRequests[f] = false
		}
	}
	cs.ID++
	cs.Origin = Elevator_id
	return cs
}

func (cs HRAInput) Update_local_state(local_elevator_state elevator.State) {

	// Create a temporary variable to hold the updated state
	tempState := Unacked_Commonstate.States[Elevator_id]
	tempState.Behaviour = string(local_elevator_state.Behaviour)

	// Assign the updated state back to the map
	Unacked_Commonstate.States[Elevator_id] = tempState
}


func Commonstates_are_equal(new_commonstate, Commonstate HRAInput) bool {	

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

func Fully_acked(ackmap map[string]Ack_status) bool {
	for id, value := range ackmap {
		if value == 0 && id != Elevator_id {
			return false
		}
	}
	return true
}	

func Higher_priority(cs1, cs2 HRAInput) bool {

	if cs1.ID > cs2.ID {
		return true
	}

	id1 := cs1.Origin
	id2 := cs2.Origin
    parts1 := strings.Split(id1, "-")
    parts2 := strings.Split(id2, "-")
    ip1 := net.ParseIP(parts1[1])
    ip2 := net.ParseIP(parts2[1])
    pid1, _ := strconv.Atoi(parts1[2])
    pid2, _ := strconv.Atoi(parts2[2])

    // Compare IP addresses
    cmp := bytes.Compare(ip1, ip2)
    if cmp > 0 {
        return true
    } else if cmp < 0 {
        return false
    }

    // If IP addresses are equal, compare process IDs
    return pid1 > pid2
}