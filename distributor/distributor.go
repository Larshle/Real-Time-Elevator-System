package distributor

import (
	"fmt"
	"reflect"
	"root/elevio"
	"root/elevator"
	"root/network/peers"
)

type AckStatus int
const (
	NotAcked AckStatus = iota
	Acked
	NotAvailable
)

type LocalElevState struct {
	Behaviour   string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type CommonState struct {
	Seq          int
	Origin       int
	Ackmap       map[int]AckStatus
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[int]LocalElevState `json:"states"`
}

func (es *CommonState) AddCall(newCall elevio.ButtonEvent, ElevatorID int) {
	if newCall.Button == elevio.BT_Cab {
		es.States[ElevatorID].CabRequests[newCall.Floor] = true
	} else {
		es.HallRequests[newCall.Floor][newCall.Button] = true
	}
	es.Seq++
	es.Origin = ElevatorID
}

func (es *CommonState) removeCall(deliveredAssingement elevio.ButtonEvent, ElevatorID int) {
	if deliveredAssingement.Button == elevio.BT_Cab {
		es.States[ElevatorID].CabRequests[deliveredAssingement.Floor] = false
	} else {
		es.HallRequests[deliveredAssingement.Floor][deliveredAssingement.Button] = false
	}
	es.Seq++
	es.Origin = ElevatorID
}

func (es *CommonState) updateLocalElevState(localElevState elevator.State, ElevatorID int) {
	HRA := es.States[ElevatorID]
	HRA.Behaviour = localElevState.Behaviour.ToString()
	HRA.Floor = localElevState.Floor
	HRA.Direction = localElevState.Direction.ToString()
	HRA.CabRequests = es.States[ElevatorID].CabRequests
	es.States[ElevatorID] = HRA
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

func FullyAcked(ackmap map[int]AckStatus) bool {
	for _, value := range ackmap {
		if value == 0 {
			return false
		}
	}
	return true
}

func commonStatesEqual(oldCS, newCS CommonState) bool {
	oldCS.Ackmap = nil
	newCS.Ackmap = nil
	return reflect.DeepEqual(oldCS, newCS)
}

func (cs *CommonState) makeElevUnav(p peers.PeerUpdate) {
	for _, id := range p.Lost {
		cs.Ackmap[id] = NotAvailable
	}
}

func (cs *CommonState) makeElevav(ElevatorID int) {
	if cs.Ackmap[ElevatorID] == NotAvailable {
		cs.Ackmap[ElevatorID] = NotAcked
	}
}

func (cs *CommonState) Ack(ElevatorID int) {
	cs.Ackmap[ElevatorID] = Acked
}

func (cs *CommonState) makeElevUnavExceptOrigin(ElevatorID int) {
	for id := range cs.Ackmap {
		if id != ElevatorID {
			cs.Ackmap[id] = NotAvailable
		}
	}
}

func (cs *CommonState) NullAckmap() {
	for id := range cs.Ackmap {
		if cs.Ackmap[id] == Acked {
			cs.Ackmap[id] = NotAcked
		}
	}
}
