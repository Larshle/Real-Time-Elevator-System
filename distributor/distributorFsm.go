package distributor

import (
	"root/config"
	"root/elevator"
	"root/elevio"
	"root/network/peers"
	"time"
)

type StashType int

const (
	None StashType = iota
	RemoveCall
	AddCall
	StateChange
)

func Distributor(
	deliveredAssignmentC <-chan elevio.ButtonEvent,
	newLocalElevStateC <-chan elevator.State,
	giverToNetworkC chan<- CommonState,
	receiverFromNetworkC <-chan CommonState,
	toAssignerC chan<- CommonState,
	receiverPeersC <-chan peers.PeerUpdate,
	id int) {

	elevioOrdersC := make(chan elevio.ButtonEvent, 10000)

	go elevio.PollButtons(elevioOrdersC)

	var stateStash elevator.State
	var NewOrderStash elevio.ButtonEvent
	var RemoveOrderStash elevio.ButtonEvent
	var stashType StashType
	var peers peers.PeerUpdate

	var cs CommonState

	disconnectTimer := time.NewTimer(config.DisconnectTime)
	heartbeatTimer := time.NewTicker(config.HeartbeatTime)

	acking := false
	aloneOnNetwork := false
	stashed := false

	for {
		select {
		case <-disconnectTimer.C:
			aloneOnNetwork = true
			cs.makeOthersUnavailable(id)

		case P := <-receiverPeersC:
			peers = P

		case <-heartbeatTimer.C:
			giverToNetworkC <- cs

		default:
		}

		switch {
		
		case acking:
			select {
			case arrivedCs := <-receiverFromNetworkC:
				disconnectTimer = time.NewTimer(config.DisconnectTime)
				if arrivedCs.Seq < cs.Seq {
					break
				}
				switch {
				case (arrivedCs.Origin > cs.Origin && arrivedCs.Seq == cs.Seq) || arrivedCs.Seq > cs.Seq:
					cs = arrivedCs
					cs.Ackmap[id] = Acked
					cs.makeLostPeersUnavailable(peers)

				case arrivedCs.fullyAcked(id):
					cs = arrivedCs
					toAssignerC <- cs
					acking = false

				case cs.equals(arrivedCs):
					cs = arrivedCs
					cs.Ackmap[id] = Acked
					cs.makeLostPeersUnavailable(peers)

				default:
				}
			default:
			}
		
		case stashed:
			select{
			case arrivedCs := <-receiverFromNetworkC:
				disconnectTimer = time.NewTimer(config.DisconnectTime)
				if arrivedCs.Seq < cs.Seq {
					break
				}
				switch {
				case (arrivedCs.Origin > cs.Origin && arrivedCs.Seq == cs.Seq) || arrivedCs.Seq > cs.Seq:
					cs = arrivedCs
					cs.Ackmap[id] = Acked
					cs.makeLostPeersUnavailable(peers)
					acking = true

				case arrivedCs.fullyAcked(id):
					cs = arrivedCs
					toAssignerC <- cs
					stashed = false
				}
		case aloneOnNetwork:
			select {
			case <-receiverFromNetworkC:
				if cs.States[id].CabRequests == [config.NumFloors]bool{} {
					aloneOnNetwork = false
				}

			case newOrder := <-elevioOrdersC:
				if cs.States[id].State.Motorstop {
					break
				}
				cs.Ackmap[id] = Acked
				cs.addCabCall(newOrder, id)
				toAssignerC <- cs

			case removeOrder := <-deliveredAssignmentC:
				cs.Ackmap[id] = Acked
				cs.removeAssignments(removeOrder, id)
				toAssignerC <- cs

			case newElevState := <-newLocalElevStateC:
				cs.Ackmap[id] = Acked
				cs.updateLocalElevState(newElevState, id)
				toAssignerC <- cs

			default:
			}

		default:
			select {
			case newOrder := <-elevioOrdersC:
				stashType = AddCall
				NewOrderStash = newOrder
				cs.prepNewCs(id)
				cs.addAssignments(newOrder, id)
				cs.Ackmap[id] = Acked
				acking = true
				stashed = true

			case removeOrder := <-deliveredAssignmentC:
				stashType = RemoveCall
				RemoveOrderStash = removeOrder
				cs.prepNewCs(id)
				cs.removeAssignments(removeOrder, id)
				cs.Ackmap[id] = Acked
				acking = true
				stashed = true

			case newElevState := <-newLocalElevStateC:
				stashType = StateChange
				stateStash = newElevState
				cs.prepNewCs(id)
				cs.updateLocalElevState(newElevState, id)
				cs.Ackmap[id] = Acked
				acking = true
				stashed = true

			case arrivedCs := <-receiverFromNetworkC:
				disconnectTimer = time.NewTimer(config.DisconnectTime)
				if (arrivedCs.Origin > cs.Origin && arrivedCs.Seq == cs.Seq) || arrivedCs.Seq > cs.Seq {
					cs = arrivedCs
					cs.makeLostPeersUnavailable(peers)
					cs.Ackmap[id] = Acked
					acking = true
				}

			default:
			}
		}
	}

}
}