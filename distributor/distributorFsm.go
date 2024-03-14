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

	cs := initCommonState(id)

	disconnectTimer := time.NewTimer(config.DisconnectTime)
	heartbeatTimer := time.NewTicker(config.HeartbeatTime)

	acking := false
	aloneOnNetwork := false

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
		case !acking:
			select {
			case newOrder := <-elevioOrdersC:
				stashType = AddCall
				NewOrderStash = newOrder
				cs.prepNewCs(id)
				cs.addAssignments(newOrder, id)
				cs.Ackmap[id] = Acked
				acking = true

			case removeOrder := <-deliveredAssignmentC:
				stashType = RemoveCall
				RemoveOrderStash = removeOrder
				cs.prepNewCs(id)
				cs.removeAssignments(removeOrder, id)
				cs.Ackmap[id] = Acked
				acking = true

			case newElevState := <-newLocalElevStateC:
				stashType = StateChange
				stateStash = newElevState
				cs.prepNewCs(id)
				cs.updateLocalElevState(newElevState, id)
				cs.Ackmap[id] = Acked
				acking = true

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

		case aloneOnNetwork:
			select {
			case <-receiverFromNetworkC:
				if cs.States[id].CabRequests == [config.NumFloors]bool{} {
					aloneOnNetwork = false
				}

			case newOrder := <-elevioOrdersC:
				if cs.States[id].Stuck {
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

		case acking:
			select {
			case arrivedCs := <-receiverFromNetworkC:
				if arrivedCs.Seq < cs.Seq {
					break
				}
				disconnectTimer = time.NewTimer(config.DisconnectTime)

				switch {
				case (arrivedCs.Origin > cs.Origin && arrivedCs.Seq == cs.Seq) || arrivedCs.Seq > cs.Seq:
					cs = arrivedCs
					cs.Ackmap[id] = Acked
					cs.makeLostPeersUnavailable(peers)

				case arrivedCs.fullyAcked(id):
					cs = arrivedCs
					toAssignerC <- cs
					switch {
					case cs.Origin != id && stashType != None:
						cs.prepNewCs(id)

						switch stashType {
						case AddCall:
							cs.addAssignments(NewOrderStash, id)
							cs.Ackmap[id] = Acked

						case RemoveCall:
							cs.removeAssignments(RemoveOrderStash, id)
							cs.Ackmap[id] = Acked

						case StateChange:
							cs.updateLocalElevState(stateStash, id)
							cs.Ackmap[id] = Acked
						}

					case cs.Origin == id && stashType != None:
						stashType = None
						acking = false

					default:
						acking = false
					}

				case cs.equals(arrivedCs):
					cs = arrivedCs
					cs.Ackmap[id] = Acked
					cs.makeLostPeersUnavailable(peers)

				default:
				}
			default:
			}
		}
	}
}
