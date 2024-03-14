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
	ElevatorID int) {

	elevioOrdersC := make(chan elevio.ButtonEvent, 10000)

	go elevio.PollButtons(elevioOrdersC)

	var stateStash elevator.State
	var NewOrderStash elevio.ButtonEvent
	var RemoveOrderStash elevio.ButtonEvent
	var stashType StashType
	var peers peers.PeerUpdate

	cs := initCommonState(ElevatorID)

	disconnectTimer := time.NewTimer(config.DisconnectTime)
	heartbeatTimer := time.NewTicker(config.HeartbeatTime)

	acking := false
	aloneOnNetwork := false

	for {
		select {
		case <-disconnectTimer.C:
			aloneOnNetwork = true
			cs.makeOthersUnavailable(ElevatorID)

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
				cs.prepNewCs(ElevatorID)
				cs.addAssignments(newOrder, ElevatorID)
				cs.Ackmap[ElevatorID] = Acked
				acking = true

			case removeOrder := <-deliveredAssignmentC:
				stashType = RemoveCall
				RemoveOrderStash = removeOrder
				cs.prepNewCs(ElevatorID)
				cs.removeAssignments(removeOrder, ElevatorID)
				cs.Ackmap[ElevatorID] = Acked
				acking = true

			case newElevState := <-newLocalElevStateC:
				stashType = StateChange
				stateStash = newElevState
				cs.prepNewCs(ElevatorID)
				cs.updateLocalElevState(newElevState, ElevatorID)
				cs.Ackmap[ElevatorID] = Acked
				acking = true

			case arrivedCs := <-receiverFromNetworkC:
				disconnectTimer = time.NewTimer(config.DisconnectTime)
				if (arrivedCs.Origin > cs.Origin && arrivedCs.Seq == cs.Seq) || arrivedCs.Seq > cs.Seq {
					cs = arrivedCs
					cs.makeLostPeersUnavailable(peers)
					cs.Ackmap[ElevatorID] = Acked
					acking = true
				}

			default:
			}

		case aloneOnNetwork:
			select {
			case <-receiverFromNetworkC:
				if cs.States[ElevatorID].CabRequests == [config.NumFloors]bool{} {
					aloneOnNetwork = false
				}

			case newOrder := <-elevioOrdersC:
				if cs.States[ElevatorID].Stuck {
					break
				}
				cs.Ackmap[ElevatorID] = Acked
				cs.addCabCall(newOrder, ElevatorID)
				toAssignerC <- cs

			case removeOrder := <-deliveredAssignmentC:
				cs.Ackmap[ElevatorID] = Acked
				cs.removeAssignments(removeOrder, ElevatorID)
				toAssignerC <- cs

			case newElevState := <-newLocalElevStateC:
				cs.Ackmap[ElevatorID] = Acked
				cs.updateLocalElevState(newElevState, ElevatorID)
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
					cs.Ackmap[ElevatorID] = Acked
					cs.makeLostPeersUnavailable(peers)

				case arrivedCs.fullyAcked(ElevatorID):
					cs = arrivedCs
					toAssignerC <- cs
					switch {
					case cs.Origin != ElevatorID && stashType != None:
						cs.prepNewCs(ElevatorID)

						switch stashType {
						case AddCall:
							cs.addAssignments(NewOrderStash, ElevatorID)
							cs.Ackmap[ElevatorID] = Acked

						case RemoveCall:
							cs.removeAssignments(RemoveOrderStash, ElevatorID)
							cs.Ackmap[ElevatorID] = Acked

						case StateChange:
							cs.updateLocalElevState(stateStash, ElevatorID)
							cs.Ackmap[ElevatorID] = Acked
						}

					case cs.Origin == ElevatorID && stashType != None:
						stashType = None
						acking = false

					default:
						acking = false
					}

				case cs.equals(arrivedCs):
					cs = arrivedCs
					cs.Ackmap[ElevatorID] = Acked
					cs.makeLostPeersUnavailable(peers)

				default:
				}
			default:
			}
		}
	}
}
