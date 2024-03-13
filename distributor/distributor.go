package distributor

import (
	"fmt"
	"reflect"
	"root/config"
	"root/elevator"
	"root/elevio"
	"root/network/peers"
)

type AckStatus int

const (
	NotAcked AckStatus = iota
	Acked
	NotAvailable
)

type LocalElevState struct {
	Behaviour   string                 `json:"behaviour"`
	Floor       int                    `json:"floor"`
	Direction   string                 `json:"direction"`
	CabRequests [config.NumFloors]bool `json:"cabRequests"`
}

type CommonState struct {
	Seq          int
	Origin       int
	Ackmap       [config.NumElevators]AckStatus
	HallRequests [config.NumFloors][2]bool           `json:"hallRequests"`
	States       [config.NumElevators]LocalElevState `json:"states"`
}

func (cs *CommonState) initCommonState() {
	var hallRequests [config.NumFloors][2]bool
	var cabRequests [config.NumFloors]bool
	for f := range hallRequests {
		hallRequests[f] = [2]bool{false, false}
		cabRequests[f] = false
	}

	var ackSlice [config.NumElevators]AckStatus
	var states [config.NumElevators]LocalElevState
	for i := 0; i < config.NumElevators; i++ {
		states[i] = LocalElevState{
			Behaviour:   "idle",
			Floor:       2,
			Direction:   "down",
			CabRequests: cabRequests,
		}
		ackSlice[i] = NotAcked
	}

	cs.Seq = 0
	cs.Origin = 0
	cs.Ackmap = ackSlice
	cs.HallRequests = hallRequests
	cs.States = states
}

func (es *CommonState) addAssignments(newCall elevio.ButtonEvent, ElevatorID int) {
	if newCall.Button == elevio.BT_Cab {
		es.States[ElevatorID].CabRequests[newCall.Floor] = true
	} else {
		es.HallRequests[newCall.Floor][newCall.Button] = true
	}
	es.Seq++
	es.Origin = ElevatorID
}

func (es *CommonState) removeAssignments(deliveredAssingement elevio.ButtonEvent, ElevatorID int) {
	if deliveredAssingement.Button == elevio.BT_Cab {
		es.States[ElevatorID].CabRequests[deliveredAssingement.Floor] = false
	} else {
		es.HallRequests[deliveredAssingement.Floor][deliveredAssingement.Button] = false
	}
	es.Seq++
	es.Origin = ElevatorID
}

func (es *CommonState) addCabCall(newCall elevio.ButtonEvent, ElevatorID int) {
	if newCall.Button == elevio.BT_Cab {
		es.States[ElevatorID].CabRequests[newCall.Floor] = true
	}
}

func (es *CommonState) updateLocalElevState(localElevState elevator.State, ElevatorID int) {
	localEs := es.States[ElevatorID]
	localEs.Behaviour = localElevState.Behaviour.ToString()
	localEs.Floor = localElevState.Floor
	localEs.Direction = localElevState.Direction.ToString()
	localEs.CabRequests = es.States[ElevatorID].CabRequests
	es.States[ElevatorID] = localEs
	es.Seq++
	es.Origin = ElevatorID
}

func (cs *CommonState) Print() {
	fmt.Println("\nOrigin:", cs.Origin)
	fmt.Println("seq:", cs.Seq)
	fmt.Println("Ackmap:", cs.Ackmap)
	fmt.Println("Hall Requests:", cs.HallRequests)

	for i, state := range cs.States {
		fmt.Printf("Elevator %d:\n", int(i))
		fmt.Printf("\tBehaviour: %s\n", state.Behaviour)
		fmt.Printf("\tFloor: %d\n", state.Floor)
		fmt.Printf("\tDirection: %s\n", state.Direction)
		fmt.Printf("\tCab Requests: %v\n\n", state.CabRequests)
	}
}

func (cs *CommonState) fullyAcked(ElevatorID int) bool {
	if cs.Ackmap[ElevatorID] == NotAvailable {
		return false
	}
	for index := range cs.Ackmap {
		if cs.Ackmap[index] == NotAcked {
			return false
		}
	}
	return true
}

func commonStatesEqual(oldCS, newCS CommonState) bool {
	oldCS.Ackmap = [config.NumElevators]AckStatus{}
	newCS.Ackmap = [config.NumElevators]AckStatus{}
	return reflect.DeepEqual(oldCS, newCS)
}

func (cs *CommonState) makeLostPeersUnavailable(p peers.PeerUpdate) {
	for _, id := range p.Lost {
		cs.Ackmap[id] = NotAvailable
	}
}

func (cs *CommonState) makeOthersUnavailable(ElevatorID int) {
	for id := range cs.Ackmap {
		if id != ElevatorID {
			cs.Ackmap[id] = NotAvailable
		}
	}
}

func (cs *CommonState) nullAckmap() {
	for id := range cs.Ackmap {
		if cs.Ackmap[id] == Acked {
			cs.Ackmap[id] = NotAcked
		}
	}
}