package distributor

import (
	"bytes"
	"fmt"
	"net"
	// "reflect"
	"root/config"
	"root/elevator"
	"root/network/network_modules/peers"
	"strconv"
	"strings"
)

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
	ID           int
	Origin       string
	Ackmap       map[string]Ack_status
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

type HRAInput2 struct {
	ID           int
	Origin       string
	Ackmap       map[string]Ack_status
	HallRequests [][2]int               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func (es *HRAElevState) toHRAElevState(localElevState elevator.State) {
	es.Behaviour = localElevState.Behaviour.ToString()
	es.Floor = localElevState.Floor
	es.Direction = localElevState.Direction.ToString()
}

func PrintCommonState(cs HRAInput2) {
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

func (cs *HRAInput2) Update_Assingments(local_elevator_assignments localAssignments) {

	for f := 0; f < config.N_floors; f++ {
		for b := 0; b < 2; b++ {
			if local_elevator_assignments.localHallAssignments[f][b] == add {
				cs.HallRequests[f][b] = 1
			}
			if local_elevator_assignments.localHallAssignments[f][b] == remove {
				cs.HallRequests[f][b] = 2
			}
		}
	}

	for f := 0; f < config.N_floors; f++ {
		if local_elevator_assignments.localCabAssignments[f] == add {
			cs.States[config.Elevator_id].CabRequests[f] = true
		}
		if local_elevator_assignments.localCabAssignments[f] == remove {
			cs.States[config.Elevator_id].CabRequests[f] = false
		}
	}
	cs.ID++
}

func (cs *HRAInput2) Update_local_state(local_elevator_state elevator.State) {
    hraElevState := cs.States[config.Elevator_id]
    hraElevState.toHRAElevState(local_elevator_state)
    cs.States[config.Elevator_id] = hraElevState
	cs.ID++
}

func Fully_acked(ackmap map[string]Ack_status) bool {
	// if len(ackmap) > 1 {
	for _, value := range ackmap {
		if value == 0 {
			return false
		}
	}
	return true
	// }
	// return false
}

// func commonStatesNotEqual(oldCS, newCS HRAInput2) bool {
// 	oldCS.Ackmap = nil
// 	newCS.Ackmap = nil
// 	return !reflect.DeepEqual(oldCS, newCS)
// }

func (cs *HRAInput2) makeElevUnav(p peers.PeerUpdate) {
	for _, id := range p.Lost {
		cs.Ackmap[id] = NotAvailable
		delete(cs.States, id)
	}
	cs.ID++
}

func (cs *HRAInput2) Ack() {
	cs.Ackmap[config.Elevator_id] = Acked
}

func takePriortisedCommonState(oldCS, newCS HRAInput2) HRAInput2 {
	if oldCS.ID < newCS.ID {
		return newCS
	}
	id1 := oldCS.Origin
	id2 := newCS.Origin
	parts1 := strings.Split(id1, "-")
	parts2 := strings.Split(id2, "-")
	ip1 := net.ParseIP(parts1[1])
	ip2 := net.ParseIP(parts2[1])
	pid1, _ := strconv.Atoi(parts1[2])
	pid2, _ := strconv.Atoi(parts2[2])

	// Compare IP addresses
	cmp := bytes.Compare(ip1, ip2)
	if cmp > 0 {
		return oldCS
	} else if cmp < 0 {
		return newCS
	}

	// If IP addresses are equal, compare process IDs
	if pid1 > pid2 {
		return oldCS
	}
	return newCS
}

// func (localCS *HRAInput2) MergeCommonState(globalCS HRAInput2, lc localAssignments) {
// 	globalCS.States[config.Elevator_id] = localCS.States[config.Elevator_id]
// 	for f := 0; f < config.N_floors; f++ {
// 		if lc.localCabAssignments[f] == add {
// 			localCS.States[config.Elevator_id].CabRequests[f] = true
// 		}
// 		if lc.localCabAssignments[f] == remove {
// 			localCS.States[config.Elevator_id].CabRequests[f] = false
// 		}
// 	}

// 	for f := 0; f < config.N_floors; f++ {
// 		for b := 0; b < 2; b++ {
// 			if lc.localHallAssignments[f][b] == add {
// 				globalCS.HallRequests[f][b] = true
// 			}
// 			if lc.localHallAssignments[f][b] == remove {
// 				globalCS.HallRequests[f][b] = false
// 			}
// 		}
// 	}

// 	localCS.States = globalCS.States
// 	localCS.HallRequests = globalCS.HallRequests

// 	fmt.Println("3")
// 	localCS.Ack()
// 	fmt.Println("4")
// 	localCS.Origin = config.Elevator_id
// 	localCS.ID = globalCS.ID + 1
// }

func (cs *HRAInput2) MergeCommonState(newCS HRAInput2) {
	temp := cs.States[config.Elevator_id] 
	cs.States = newCS.States
	cs.States[config.Elevator_id] = temp
	cs.Ackmap = newCS.Ackmap
	cs.Ack()
	cs.Origin = config.Elevator_id
	
	for f := 0; f < config.N_floors; f++ {
		for b := 0; b < 2; b++ {
			if newCS.HallRequests[f][b] == 2{
				cs.HallRequests[f][b] = 2
			}
			if newCS.HallRequests[f][b] == 1 {
				if cs.HallRequests[f][b] == 0 {
					cs.HallRequests[f][b] = 1
				}
				if cs.HallRequests[f][b] == 1 {
					cs.HallRequests[f][b] = 1
				}
				if cs.HallRequests[f][b] == 2 {
					cs.HallRequests[f][b] = 2
				}
			}
			if newCS.HallRequests[f][b] == 0{
				if cs.HallRequests[f][b] == 0 {
					cs.HallRequests[f][b] = 0
				}
				if cs.HallRequests[f][b] == 1 {
					cs.HallRequests[f][b] = 1
				}
				if cs.HallRequests[f][b] == 2 {
					cs.HallRequests[f][b] = 2
				}
			}
		}
	}
}
