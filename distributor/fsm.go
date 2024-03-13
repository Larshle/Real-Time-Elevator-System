package distributor

import (
	"fmt"
	"root/config"
	"root/elevator"
	"root/elevio"
	"root/network/peers"
	"time"
)

type StashType int

const (
	RemoveCall StashType = iota
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

	var StateStash elevator.State
	var NewOrderStash elevio.ButtonEvent
	var RemoveOrderStash elevio.ButtonEvent
	var StashType StashType
	var peers peers.PeerUpdate
	var cs CommonState

	cs.initCommonState()

	disconnectTimer := time.NewTimer(config.DisconnectTime)
	heartbeatTimer := time.NewTicker(15 * time.Millisecond)

	stashed := false
	acking := false
	aloneOnNetwork := false

	for {

		select {
		case <-disconnectTimer.C:
			aloneOnNetwork = true
			cs.makeOthersUnavailable(ElevatorID)

		case P := <-receiverPeersC:
			peers = P

		default:
		}

		switch {
		case !acking: // Idle
			select {
			case newOrder := <-elevioOrdersC:
				NewOrderStash = newOrder
				StashType = AddCall
				cs.addAssignments(newOrder, ElevatorID)
				cs.nullAckmap()
				cs.Ackmap[ElevatorID] = Acked
				stashed = true
				acking = true

			case removeOrder := <-deliveredAssignmentC:
				RemoveOrderStash = removeOrder
				StashType = RemoveCall
				cs.removeAssignments(removeOrder, ElevatorID)
				cs.nullAckmap()
				cs.Ackmap[ElevatorID] = Acked
				stashed = true
				acking = true

			case newElevState := <-newLocalElevStateC:
				StateStash = newElevState
				StashType = StateChange
				cs.updateLocalElevState(newElevState, ElevatorID)
				cs.nullAckmap()
				cs.Ackmap[ElevatorID] = Acked
				stashed = true
				acking = true

			case arrivedCs := <-receiverFromNetworkC:
				disconnectTimer = time.NewTimer(config.DisconnectTime)

				switch {
				case (arrivedCs.Origin > cs.Origin && arrivedCs.Seq == cs.Seq) || arrivedCs.Seq > cs.Seq:
					cs = arrivedCs
					cs.Ackmap[ElevatorID] = Acked
					acking = true
					cs.makeLostPeersUnavailable(peers)
				}
			default:
			}

		case aloneOnNetwork:
			select {
			case <-receiverFromNetworkC:
				if cs.States[ElevatorID].CabRequests == [config.NumFloors]bool{} {
					aloneOnNetwork = false
					fmt.Println("Hello network!")
				}

			case newOrder := <-elevioOrdersC:
				cs.addCabCall(newOrder, ElevatorID)
				fmt.Println("Goodbye network :(")
				toAssignerC <- cs

			case removeOrder := <-deliveredAssignmentC:
				cs.removeAssignments(removeOrder, ElevatorID)
				fmt.Println("Goodbye network :(")
				toAssignerC <- cs

			case newElevState := <-newLocalElevStateC:
				cs.updateLocalElevState(newElevState, ElevatorID)
				fmt.Println("Goodbye network :(")
				toAssignerC <- cs

			default:
			}

		default:
			select {
			case arrivedCs := <-receiverFromNetworkC:
				if arrivedCs.Seq < cs.Seq {
					break
				}
				disconnectTimer = time.NewTimer(config.DisconnectTime)

				switch {
				case (arrivedCs.Origin > cs.Origin && arrivedCs.Seq == cs.Seq) || arrivedCs.Seq > cs.Seq: // Higher priority
					cs = arrivedCs
					cs.Ackmap[ElevatorID] = Acked
					cs.makeLostPeersUnavailable(peers)

				case arrivedCs.fullyAcked(ElevatorID):
					cs = arrivedCs
					toAssignerC <- cs
					switch {
					case cs.Origin != ElevatorID && stashed:
						switch StashType {
						case AddCall:
							cs.addAssignments(NewOrderStash, ElevatorID)
							cs.nullAckmap()
							cs.Ackmap[ElevatorID] = Acked

						case RemoveCall:
							cs.removeAssignments(RemoveOrderStash, ElevatorID)
							cs.nullAckmap()
							cs.Ackmap[ElevatorID] = Acked

						case StateChange:
							cs.updateLocalElevState(StateStash, ElevatorID)
							cs.nullAckmap()
							cs.Ackmap[ElevatorID] = Acked
						}
					case cs.Origin == ElevatorID && stashed:
						stashed = false
						acking = false
					default:
						acking = false
					}

				case commonStatesEqual(cs, arrivedCs):
					cs = arrivedCs
					cs.Ackmap[ElevatorID] = Acked
					cs.makeLostPeersUnavailable(peers)

				default:
				}
			default:
			}
		}
		select {
		case <-heartbeatTimer.C:
			giverToNetworkC <- cs
		default:
		}
	}
}
