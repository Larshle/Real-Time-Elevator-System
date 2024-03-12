package distributor

import (
	"fmt"
	"root/config"
	"root/elevator"
	"root/elevio"
	"root/network/peers"
	"time"
)

type StatshType int

const (
	RemoveCall StatshType = iota
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

	var cs CommonState
	var StateStash elevator.State
	var NewOrderStash elevio.ButtonEvent
	var RemoveOrderStash elevio.ButtonEvent
	var StashType StatshType
	var peers peers.PeerUpdate

	disconnectTimer := time.NewTimer(config.DisconnectTime)
	heartbeatTimer := time.NewTicker(15 * time.Millisecond)

	stashed := false
	acking := false
	isolated := false

	cs.initCommonState()

	for {

		select {
		case <-disconnectTimer.C:
			isolated = true
		case P := <-receiverPeersC:
			peers = P
			cs.makeElevav(ElevatorID)
			fmt.Println("Peers", peers)

		default:
		}

		switch {
		case !acking: // Idle
			select {

			case newOrder := <-elevioOrdersC:
				NewOrderStash = newOrder
				StashType = AddCall
				cs.AddCall(newOrder, ElevatorID)
				cs.NullAckmap()
				cs.Ackmap[ElevatorID] = Acked
				stashed = true
				acking = true

			case removeOrder := <-deliveredAssignmentC:
				RemoveOrderStash = removeOrder
				StashType = RemoveCall
				cs.removeCall(removeOrder, ElevatorID)
				cs.NullAckmap()
				cs.Ackmap[ElevatorID] = Acked
				stashed = true
				acking = true

			case newElevState := <-newLocalElevStateC:
				StateStash = newElevState
				StashType = StateChange
				cs.updateLocalElevState(newElevState, ElevatorID)
				cs.NullAckmap()
				cs.Ackmap[ElevatorID] = Acked
				stashed = true
				acking = true

			case arrivedCommonState := <-receiverFromNetworkC:
				disconnectTimer = time.NewTimer(config.DisconnectTime)

				switch {
				case (arrivedCommonState.Origin > cs.Origin && arrivedCommonState.Seq == cs.Seq) || arrivedCommonState.Seq > cs.Seq:
					cs = arrivedCommonState
					cs.Ackmap[ElevatorID] = Acked
					acking = true
					cs.makeElevUnav(peers)
				}
			default:
			}

		case isolated:
			select {
			case <-receiverFromNetworkC:
				isolated = false

			case newOrder := <-elevioOrdersC:
				cs.makeElevUnavExceptOrigin(ElevatorID)
				cs.AddCabCall(newOrder, ElevatorID)
				fmt.Println("ISOLATED")
				toAssignerC <- cs

			case removeOrder := <-deliveredAssignmentC:
				cs.makeElevUnavExceptOrigin(ElevatorID)
				cs.removeCabCall(removeOrder, ElevatorID)
				fmt.Println("ISOLATED")
				toAssignerC <- cs

			case newElevState := <-newLocalElevStateC:
				cs.updateLocalElevState(newElevState, ElevatorID)
				cs.makeElevUnavExceptOrigin(ElevatorID)
				fmt.Println("ISOLATED")
				toAssignerC <- cs

			default:
			}

		default:
			select {
			case arrivedCommonState := <-receiverFromNetworkC:
				if arrivedCommonState.Seq < cs.Seq {
					break
				}
				disconnectTimer = time.NewTimer(config.DisconnectTime)

				switch {
				case (arrivedCommonState.Origin > cs.Origin && arrivedCommonState.Seq == cs.Seq) || arrivedCommonState.Seq > cs.Seq:
					cs = arrivedCommonState
					cs.Ackmap[ElevatorID] = Acked
					cs.makeElevUnav(peers)

				case arrivedCommonState.FullyAcked():
					cs = arrivedCommonState
					toAssignerC <- cs
					switch {
					case cs.Origin != ElevatorID && stashed:
						switch StashType {
						case AddCall:
							cs.AddCall(NewOrderStash, ElevatorID)
							cs.NullAckmap()
							cs.Ackmap[ElevatorID] = Acked

						case RemoveCall:
							cs.removeCall(RemoveOrderStash, ElevatorID)
							cs.NullAckmap()
							cs.Ackmap[ElevatorID] = Acked

						case StateChange:
							cs.updateLocalElevState(StateStash, ElevatorID)
							cs.NullAckmap()
							cs.Ackmap[ElevatorID] = Acked
						}
					case cs.Origin == ElevatorID && stashed:
						stashed = false
						acking = false
					default:
						acking = false
					}

				case commonStatesEqual(cs, arrivedCommonState):
					cs = arrivedCommonState
					cs.Ackmap[ElevatorID] = Acked
					cs.makeElevUnav(peers)

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
