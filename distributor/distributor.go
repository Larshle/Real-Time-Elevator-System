package distributor

import (
	"bytes"
	"fmt"
	"net"
	"root/config"
	"root/elevator"
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

func (cs *HRAInput) Update_Assingments(local_elevator_assignments localAssignments) {

	for f := 0; f < config.N_floors; f++ {
		for b := 0; b < 2; b++ {
			if local_elevator_assignments.localHallAssignments[f][b] == add {
				cs.HallRequests[f][b] = true
			}
			if local_elevator_assignments.localHallAssignments[f][b] == remove {
				cs.HallRequests[f][b] = false
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
	cs.Origin = config.Elevator_id
}

func (cs *HRAInput) Update_local_state(local_elevator_state elevator.State) {
	hraElevState := cs.States[config.Elevator_id]

	hraElevState.toHRAElevState(local_elevator_state)

	cs.States[config.Elevator_id] = hraElevState

	cs.ID++
	cs.Origin = config.Elevator_id
}

func Fully_acked(ackmap map[string]Ack_status) bool {
	for id, value := range ackmap {
		if value == 0 && id != config.Elevator_id {
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
