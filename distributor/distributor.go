package distributor

import (
	"bytes"
	"fmt"
	"net"
	"reflect"
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
	Seq           int
	Origin       string
	Ackmap       map[string]Ack_status
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func (input *HRAInput)ensureElevatorState( state HRAElevState) {
    _, exists := input.States[config.Elevator_id]
    if !exists {
        input.States[config.Elevator_id] = state
    }
	input.Seq++
}

func (es *HRAInput) toHRAElevState(localElevState elevator.State) {
	HRA := es.States[config.Elevator_id]
	HRA.Behaviour = localElevState.Behaviour.ToString()
	HRA.Floor = localElevState.Floor
	HRA.Direction = localElevState.Direction.ToString()
	HRA.CabRequests = es.States[config.Elevator_id].CabRequests
	es.States[config.Elevator_id] = HRA
	es.Seq++	
	es.Origin = config.Elevator_id
}

func PrintCommonState(cs HRAInput) {
	fmt.Println("\nOrigin:", cs.Origin)
	fmt.Println("seq:", cs.Seq)
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

func (cs *HRAInput) Update_Assingments(local_elevator_assignments localAssignments) {

	for f := 0; f < config.N_floors; f++ {
		for b := 0; b < 2; b++ {
			if local_elevator_assignments.localHallAssignments[f][b] == add {
				cs.HallRequests[f][b] = true
				fmt.Println("Hall request added")
			}
			if local_elevator_assignments.localHallAssignments[f][b] == remove {
				cs.HallRequests[f][b] = false
				fmt.Println("Hall request removed")
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
	
	cs.Seq++
	cs.Origin = config.Elevator_id
	fmt.Println("Updated common state with Update_Assignment:")
	//PrintCommonState(*cs)

}

func (cs *HRAElevState) Update_local_state(local_elevator_state elevator.State) {
	cs.Behaviour = local_elevator_state.Behaviour.ToString()
	cs.Floor = local_elevator_state.Floor
	cs.Direction = local_elevator_state.Direction.ToString()
	
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
	}
}

func (cs *HRAInput) Ack() {
	cs.Ackmap[config.Elevator_id] = Acked
	fmt.Println("Acked")
}

func higherPriority(oldCS, newCS HRAInput) bool {
	return oldCS.Seq > newCS.Seq || oldCS.Origin > newCS.Origin && oldCS.Seq == newCS.Seq
}

func takePriortisedCommonState(oldCS, newCS HRAInput) HRAInput {
	if oldCS.Seq < newCS.Seq {
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
func (cs *HRAInput) makeElevUnavExceptOrigin() {
	for id := range cs.Ackmap {
		if id != config.Elevator_id {
			cs.Ackmap[id] = NotAvailable
		}
	}
}

func (cs *HRAInput) UpdateCabAssignments(local_elevator_assignments localAssignments) {
	for f := 0; f < config.N_floors; f++ {
		if local_elevator_assignments.localCabAssignments[f] == add {
			cs.States[config.Elevator_id].CabRequests[f] = true
		}
		if local_elevator_assignments.localCabAssignments[f] == remove {
			cs.States[config.Elevator_id].CabRequests[f] = false
		}
	}
	cs.Seq++
	fmt.Println("Updated common state with CabAssignment:")
	PrintCommonState(*cs)
}
func (cs *HRAInput) makeOriginElevUnav(){
	cs.Ackmap[config.Elevator_id] = NotAvailable
}

func (cs *HRAInput) NullAckmap() {
	for id := range cs.Ackmap {
		cs.Ackmap[id] = NotAcked
	}
}