package distributor

import (
	"root/elevio"
	"root/elevator"
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

	selfLostNetworkDuratio := 5 * time.Second
	timeCounter := time.NewTimer(selfLostNetworkDuratio)
	heartbeatTimer := time.NewTicker(15 * time.Millisecond)

	stashed := false
	acking := false
	isolated := false


	commonState = CommonState{
		Origin: 0,
		Seq:    0,
		Ackmap: map[int]AckStatus{
			0: NotAcked,
			1: NotAcked,
			2: NotAcked,
		},
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		States: map[int]LocalElevState{
			0: {
				Behaviour:   "idle",
				Floor:       2,
				Direction:   "down",
				CabRequests: []bool{false, false, false, false},
			},
			1: {
				Behaviour:   "idle",
				Floor:       2,
				Direction:   "down",
				CabRequests: []bool{false, false, false, false},
			},
			2: {
				Behaviour:   "idle",
				Floor:       2,
				Direction:   "down",
				CabRequests: []bool{false, false, false, false},
			},
		},
	}

	for {

		select {
		case <-timeCounter.C:
			isolated = true
		default:
		}

		switch {
		case !stashed && !acking: // Idle
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
				timeCounter = time.NewTimer(selfLostNetworkDuratio)

				switch {
				case (arrivedCommonState.Origin > commonState.Origin && arrivedCommonState.Seq == commonState.Seq) || arrivedCommonState.Seq > commonState.Seq:
					commonState = arrivedCommonState
					commonState.Ack(ElevatorID)
					acking = true
				}
			case peers := <-receiverPeersC:
				commonState.makeElevUnav(peers)
				commonState.makeElevav(ElevatorID)
			default:
			}

		case isolated:
			select {
			case <-receiverFromNetworkC:
				isolated = false

			// må legge til 2 caser med add/remove call

			case newElevState := <-newLocalElevStateC:
				commonState.updateLocalElevState(newElevState, ElevatorID)
				commonState.makeElevUnavExceptOrigin(ElevatorID)
				toAssignerC <- commonState

			default:
			}

		default:
			select {
			case arrivedCommonState := <-receiverFromNetworkC:
				if arrivedCommonState.Seq < commonState.Seq {
					break
				}
				timeCounter = time.NewTimer(selfLostNetworkDuratio)

				switch {
				case (arrivedCommonState.Origin > commonState.Origin && arrivedCommonState.Seq == commonState.Seq) || arrivedCommonState.Seq > commonState.Seq:
					commonState = arrivedCommonState
					commonState.Ack(ElevatorID)

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

				default:
				}
			case peers := <-receiverPeersC:
				commonState.makeElevUnav(peers)
				commonState.makeElevav(ElevatorID)
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
