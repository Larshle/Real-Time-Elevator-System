package distributor

import (
	"bytes"
	"fmt"
	"net"
	"root/config"
	"root/elevator"
	"root/network/network_modules/peers"
	"strconv"
	"strings"
	"reflect"
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
	Unavailable  int
	ID           int
	Origin       string
	Ackmap       map[string]Ack_status 
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func (es *HRAElevState) toHRAElevState(localElevState elevator.State) {
	es.Behaviour = localElevState.Behaviour.ToString()
	es.Floor = localElevState.Floor
	es.Direction = localElevState.Direction.ToString()
}

func PrintCommonState(cs HRAInput) {
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

func (holding localAssignments) Update_Assingments(local_elevator_assignments localAssignments) {

	for f := 0; f < config.N_floors; f++ {
		for b := 0; b < 3; b++ {
			if local_elevator_assignments.localHallAssignments[f][b] == add {
				holding.localHallAssignments[f][b] = add
			}
			if local_elevator_assignments.localHallAssignments[f][b] == remove {
				holding.localHallAssignments[f][b] = remove
				fmt.Println("Hall request removed")
			}
		}
	}
}

func (cs *HRAInput) Update_local_state(local_elevator_state elevator.State) {
	hraElevState := cs.States[config.Elevator_id]

	hraElevState.toHRAElevState(local_elevator_state)

	cs.States[config.Elevator_id] = hraElevState

}

func Fully_acked(ackmap map[string]Ack_status) bool {
	for _, value := range ackmap {
		if value == 0 {
			return false
		}
	}
	return true
}


func commonStatesNotEqual(oldCS, newCS HRAInput) bool {
	oldCS.Ackmap = nil
	newCS.Ackmap = nil
	return !reflect.DeepEqual(oldCS, newCS)
}

func (cs *HRAInput) makeElevUnav(p peers.PeerUpdate) {
	for _, id := range p.Lost {
		cs.Ackmap[id] = NotAvailable
		delete(cs.States, id)
	}
}

func (cs *HRAInput) Ack() {
	cs.Ackmap[config.Elevator_id] = Acked
}


func takePriortisedCommonState(oldCS, newCS HRAInput) HRAInput{
	if(oldCS.ID < newCS.ID){
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

func (localCS *HRAInput)MergeCommonState(globalCS HRAInput, lc localAssignments){
	globalCS.States[config.Elevator_id] = localCS.States[config.Elevator_id]
	for f := 0; f < config.N_floors; f++ {
		if lc.localCabAssignments[f] == add {
			localCS.States[config.Elevator_id].CabRequests[f] = true
		}
		if lc.localCabAssignments[f] == remove {
			localCS.States[config.Elevator_id].CabRequests[f] = false
		}
	}

	for f := 0; f < config.N_floors; f++ {
		for b := 0; b < 3; b++ {
			if lc.localHallAssignments[f][b] == add {
				globalCS.HallRequests[f][b] = true
			}
			if lc.localHallAssignments[f][b] == remove {
				globalCS.HallRequests[f][b] = false
			}
		}
	}

	localCS.States = globalCS.States
	localCS.HallRequests = globalCS.HallRequests

	localCS.Ack()
	localCS.Origin = config.Elevator_id
	localCS.ID = globalCS.ID + 1

}
