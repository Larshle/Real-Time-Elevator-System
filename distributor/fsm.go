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

	var commonState CommonState
	var StateStash elevator.State
	var NewOrderStash elevio.ButtonEvent
	var RemoveOrderStash elevio.ButtonEvent
	var StashType StatshType
	var P peers.PeerUpdate

	disconnectTimer := time.NewTimer(config.DisconnectTime)
	heartbeatTimer := time.NewTicker(15 * time.Millisecond)

	stashed := false
	acking := false
	isolated := false

	commonState = initCommonState()

	// commonState = CommonState{
	// 	Origin: 0,
	// 	Seq:    0,
	// 	Ackmap: []AckStatus{NotAcked,NotAcked,NotAck				toAssignerC <- commonStateed},
	// 	HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
	// 	States: []LocalElevState{
	// 		{
	// 			Behaviour:   "idle",
	// 			Floor:       2,
	// 			Direction:   "down",
	// 			CabRequests: []bool{false, false, false, false},
	// 		},
	// 		{
	// 			Behaviour:   "idle",
	// 			Floor:       2,
	// 			Direction:   "down",
	// 			CabRequests: []bool{false, false, false, false},
	// 		},
	// 		{
	// 			Behaviour:   "idle",
	// 			Floor:       2,
	// 			Direction:   "down",
	// 			CabRequests: []bool{false, false, false, false},
	// 		},
	// 	},
	// }

	for {

		select {
		case <-disconnectTimer.C:
			isolated = true
		case Penis := <-receiverPeersC:
			P = Penis
			commonState.makeElevav(ElevatorID)
			fmt.Println("Penis", P)

		default:
		}

		switch {
		case !acking: // Idle
			select {

			case newOrder := <-elevioOrdersC:
				NewOrderStash = newOrder
				StashType = AddCall
				commonState.AddCall(newOrder, ElevatorID)
				commonState.NullAckmap()
				commonState.Ack(ElevatorID)
				stashed = true
				acking = true

			case removeOrder := <-deliveredAssignmentC:
				RemoveOrderStash = removeOrder
				StashType = RemoveCall
				commonState.removeCall(removeOrder, ElevatorID)
				commonState.NullAckmap()
				commonState.Ack(ElevatorID)
				stashed = true
				acking = true

			case newElevState := <-newLocalElevStateC:
				StateStash = newElevState
				StashType = StateChange
				commonState.updateLocalElevState(newElevState, ElevatorID)
				commonState.NullAckmap()
				commonState.Ack(ElevatorID)
				stashed = true
				acking = true

			case arrivedCommonState := <-receiverFromNetworkC:
				disconnectTimer = time.NewTimer(config.DisconnectTime)

				switch {
				case (arrivedCommonState.Origin > commonState.Origin && arrivedCommonState.Seq == commonState.Seq) || arrivedCommonState.Seq > commonState.Seq:
					commonState = arrivedCommonState
					commonState.Ack(ElevatorID)
					acking = true
					commonState.makeElevUnav(P)
				}
			default:
			}

		case isolated:
			select {
			case <-receiverFromNetworkC:
				isolated = false

			case newOrder := <-elevioOrdersC:
				commonState.makeElevUnavExceptOrigin(ElevatorID)
				commonState.AddCabCall(newOrder, ElevatorID)
				fmt.Println("ISOLATED")
				toAssignerC <- commonState

			case removeOrder := <-deliveredAssignmentC:
				commonState.makeElevUnavExceptOrigin(ElevatorID)
				commonState.removeCabCall(removeOrder, ElevatorID)
				fmt.Println("ISOLATED")
				toAssignerC <- commonState

			case newElevState := <-newLocalElevStateC:
				commonState.updateLocalElevState(newElevState, ElevatorID)
				commonState.makeElevUnavExceptOrigin(ElevatorID)
				fmt.Println("ISOLATED")
				toAssignerC <- commonState

			default:
			}

		default:
			select {
			case arrivedCommonState := <-receiverFromNetworkC:
				if arrivedCommonState.Seq < commonState.Seq {
					break
				}
				disconnectTimer = time.NewTimer(config.DisconnectTime)

				switch {
				case (arrivedCommonState.Origin > commonState.Origin && arrivedCommonState.Seq == commonState.Seq) || arrivedCommonState.Seq > commonState.Seq:
					commonState = arrivedCommonState
					commonState.Ack(ElevatorID)
					commonState.makeElevUnav(P)

				case FullyAcked(arrivedCommonState.Ackmap):
					commonState = arrivedCommonState
					toAssignerC <- commonState
					switch {
					case commonState.Origin != ElevatorID && stashed:
						switch StashType {
						case AddCall:
							commonState.AddCall(NewOrderStash, ElevatorID)
							commonState.NullAckmap()
							commonState.Ack(ElevatorID)

						case RemoveCall:
							commonState.removeCall(RemoveOrderStash, ElevatorID)
							commonState.NullAckmap()
							commonState.Ack(ElevatorID)

						case StateChange:
							commonState.updateLocalElevState(StateStash, ElevatorID)
							commonState.NullAckmap()
							commonState.Ack(ElevatorID)
						}
					case commonState.Origin == ElevatorID && stashed:
						stashed = false
						acking = false
					default:
						acking = false
					}

				case commonStatesEqual(commonState, arrivedCommonState):
					commonState = arrivedCommonState
					commonState.Ack(ElevatorID)
					commonState.makeElevUnav(P)

				default:
				}
			default:
			}
		}
		select {
		case <-heartbeatTimer.C:
			giverToNetworkC <- commonState
		default:
		}
	}
}
