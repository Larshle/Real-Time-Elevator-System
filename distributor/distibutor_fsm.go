package distributor

import (
	// "fmt"
	"root/config"
	"root/driver/elevio"
	"root/elevator"
	"root/network/network_modules/peers"
	"time"
)

type StatshType int

const (
	AssingmetChange StatshType = iota
	StateChange
)

type State int

const (
	Idle State = iota
	Acking
	SendingSelf
	AckingOtherWhileTryingToSendSelf
	Isolated
	UnableToMove
)

func Distributor(
	deliveredOrderC <-chan elevio.ButtonEvent,
	newElevStateC <-chan elevator.State,
	giverToNetwork chan<- HRAInput,
	receiveFromNetworkC <-chan HRAInput,
	messageToAssinger chan<- HRAInput) {

	elevioOrdersC := make(chan elevio.ButtonEvent)
	newAssingemntC := make(chan localAssignments, 64)
	peerUpdateC := make(chan peers.PeerUpdate)

	var commonState HRAInput
	var StateStash elevator.State
	var AssignmentStash localAssignments
	var state State = Idle
	var StashType StatshType
	selfLostNetworkDuratio := 1 * time.Second

	commonState = HRAInput{
		Origin:       config.Elevator_id,
		Seq:          0,
		Ackmap:       make(map[string]Ack_status),
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		States: map[string]HRAElevState{
			config.Elevator_id: {
				Behaviour:   "idle",
				Floor:       0,
				Direction:   "up",
				CabRequests: []bool{false, false, false, true},
			},
		},
	}

	go elevio.PollButtons(elevioOrdersC)
	go Update_Assingments(elevioOrdersC, deliveredOrderC, newAssingemntC)

	heartbeatTimer := time.NewTicker(15 * time.Millisecond)
	timeCounter := time.NewTimer(500 * time.Millisecond)

	for {

		select {
		case <-timeCounter.C:
			state = Isolated
		default:
		}

		switch state {
		case Idle:
			select {
			case assingmentUpdate := <-newAssingemntC:
				AssignmentStash = assingmentUpdate
				StashType = AssingmetChange
				commonState.Update_Assingments(assingmentUpdate)
				commonState.Ack()
				state = SendingSelf

			case newElevState := <-newElevStateC:
				StateStash = newElevState
				StashType = StateChange
				commonState.toHRAElevState(newElevState)
				commonState.Ack()
				state = SendingSelf

			case arrivedCommonState := <-receiveFromNetworkC:
				timeCounter = time.NewTimer(selfLostNetworkDuratio)

				switch {
				case higherPriority(commonState, arrivedCommonState):
					commonState = arrivedCommonState
					commonState.Ack()
					state = Acking
				default:
					break //doing jack
				}
			case peers := <-peerUpdateC:
				commonState.makeElevUnav(peers)
			default:
			}
		case SendingSelf:
			select {
			case arrivedCommonState := <-receiveFromNetworkC:
				timeCounter = time.NewTimer(selfLostNetworkDuratio)
				switch {
				case arrivedCommonState.Origin != config.Elevator_id && higherPriority(commonState, arrivedCommonState):
					arrivedCommonState.Ack()
					commonState = arrivedCommonState
					state = AckingOtherWhileTryingToSendSelf

				case Fully_acked(arrivedCommonState.Ackmap):
					state = Idle
					commonState = arrivedCommonState
					messageToAssinger <- commonState
				default:
					break //doing jack
				}

			case peers := <-peerUpdateC:
				commonState.makeElevUnav(peers)
				if Fully_acked(commonState.Ackmap) {
					state = Idle
					messageToAssinger <- commonState
				}
			default:
			}

		case Acking:
			select {
			case arrivedCommonState := <-receiveFromNetworkC:
				timeCounter = time.NewTimer(selfLostNetworkDuratio)
				switch {

				case Fully_acked(arrivedCommonState.Ackmap):
					state = Idle
					commonState = arrivedCommonState
					messageToAssinger <- commonState

				case higherPriority(commonState, arrivedCommonState):
					arrivedCommonState.Ack()
					commonState = arrivedCommonState

				case !higherPriority(commonState, arrivedCommonState):
					break //doing jack
				}

			case peers := <-peerUpdateC:
				commonState.makeElevUnav(peers)
				if Fully_acked(commonState.Ackmap) {
					state = Idle
					messageToAssinger <- commonState
				}

			default:
			}

		case AckingOtherWhileTryingToSendSelf:
			select {
			case arrivedCommonState := <-receiveFromNetworkC:
				timeCounter = time.NewTimer(selfLostNetworkDuratio)
				switch {
				case !higherPriority(commonState, arrivedCommonState):
					break

				case higherPriority(commonState, arrivedCommonState):
					arrivedCommonState.Ack()
					commonState = arrivedCommonState

				case Fully_acked(arrivedCommonState.Ackmap):
					state = SendingSelf
					switch StashType {
					case AssingmetChange:
						arrivedCommonState.Update_Assingments(AssignmentStash)

					case StateChange:
						arrivedCommonState.toHRAElevState(StateStash)
					}
					commonState = arrivedCommonState
					commonState.emptyAckmap()
					commonState.Ack()
					messageToAssinger <- commonState
				}
			case peers := <-peerUpdateC:
				commonState.makeElevUnav(peers)
				if Fully_acked(commonState.Ackmap) {
					state = SendingSelf
					messageToAssinger <- commonState
				}
			default:
			}
		case Isolated:
			select {
			case <-receiveFromNetworkC:
				state = Idle

			case assingmentUpdate := <-newAssingemntC:
				commonState.makeElevUnavExceptOrigin()
				commonState.UpdateCabAssignments(assingmentUpdate)
				messageToAssinger <- commonState

			case newElevState := <-newElevStateC:
				commonState.toHRAElevState(newElevState)
				commonState.makeElevUnavExceptOrigin()
				messageToAssinger <- commonState

			default:
			}
		}

		select {
		case <-heartbeatTimer.C:
			giverToNetwork <- commonState
		default:
		}

	}
}
