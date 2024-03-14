package distributor

import (
	//"fmt"
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
	State       elevator.State
	CabRequests [config.NumFloors]bool
}

type CommonState struct {
	Seq          int
	Origin       int
	Ackmap       [config.NumElevators]AckStatus
	HallRequests [config.NumFloors][2]bool
	States       [config.NumElevators]LocalElevState
}



func (cs *CommonState) addAssignments(newCall elevio.ButtonEvent, id int) {
	if newCall.Button == elevio.BT_Cab {
		cs.States[id].CabRequests[newCall.Floor] = true
	} else {
		cs.HallRequests[newCall.Floor][newCall.Button] = true
	}
}

func (cs *CommonState) removeAssignments(deliveredAssingement elevio.ButtonEvent, id int) {
	if deliveredAssingement.Button == elevio.BT_Cab {
		cs.States[id].CabRequests[deliveredAssingement.Floor] = false
	} else {
		cs.HallRequests[deliveredAssingement.Floor][deliveredAssingement.Button] = false
	}
}

func (cs *CommonState) addCabCall(newCall elevio.ButtonEvent, id int) {
	if newCall.Button == elevio.BT_Cab {
		cs.States[id].CabRequests[newCall.Floor] = true
	}
}

func (cs *CommonState) updateLocalElevState(localElevState elevator.State, id int) {
	//localEs := cs.States[id]

	cs.States[id] = LocalElevState{
		State:       localElevState,
		CabRequests: cs.States[id].CabRequests,
	}
}
/*
func (cs *CommonState) Print() {
	fmt.Println("\nOrigin:", cs.Origin)
	fmt.Println("seq:", cs.Seq)
	fmt.Println("Ackmap:", cs.Ackmap)
	fmt.Println("Hall Requests:", cs.HallRequests)

	for i, state := range cs.States {
		fmt.Printf("Elevator %d:\n", int(i))
		fmt.Printf("\tStuck: %t\n", state.Stuck)
		fmt.Printf("\tBehaviour: %s\n", state.Behaviour)
		fmt.Printf("\tFloor: %d\n", state.Floor)
		fmt.Printf("\tDirection: %s\n", state.Direction)
		fmt.Printf("\tCab Requests: %v\n\n", state.CabRequests)
	}
}
*/

func (cs *CommonState) fullyAcked(id int) bool {
	if cs.Ackmap[id] == NotAvailable {
		return false
	}
	for index := range cs.Ackmap {
		if cs.Ackmap[index] == NotAcked {
			return false
		}
	}
	return true
}

func (oldCs CommonState) equals(newCs CommonState) bool {
	oldCs.Ackmap = [config.NumElevators]AckStatus{}
	newCs.Ackmap = [config.NumElevators]AckStatus{}
	return reflect.DeepEqual(oldCs, newCs)
}

func (cs *CommonState) makeLostPeersUnavailable(p peers.PeerUpdate) {
	for _, id := range p.Lost {
		cs.Ackmap[id] = NotAvailable
	}
}

func (cs *CommonState) makeOthersUnavailable(id int) {
	for elev := range cs.Ackmap {
		if elev != id {
			cs.Ackmap[elev] = NotAvailable
		}
	}
}

func (cs *CommonState) prepNewCs(id int) {
	cs.Seq++
	cs.Origin = id
	for id := range cs.Ackmap {
		if cs.Ackmap[id] == Acked {
			cs.Ackmap[id] = NotAcked
		}
	}
}
