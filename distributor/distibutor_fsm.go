package distributor

import (
	"fmt"
	"root/config"
	"root/driver/elevio"
	"root/elevator"
	"root/network/network_modules/peers"
	"time"
)

type StatshType int

const (
	RemoveCall StatshType = iota
	AddCall
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
	messageToAssinger chan<- HRAInput, 
	recieveFromPeerC <- chan peers.PeerUpdate) {

	elevioOrdersC := make(chan elevio.ButtonEvent, 69)
	newAssingemntC := make(chan localAssignments, 69)
	

	var commonState HRAInput
	var StateStash elevator.State
	var NewOrderStash elevio.ButtonEvent
	var RemoveOrderStash elevio.ButtonEvent
	var state State = Idle
	var StashType StatshType
	timeCounter := time.NewTimer(time.Hour)
	selfLostNetworkDuratio := 10 * time.Second

	// commonState = HRAInput{
	// 	Origin:       config.Elevator_id,
	// 	ID:           0,
	// 	Ackmap:       make(map[string]Ack_status),
	// 	HallRequests: make([][2]bool, 4), // Assuming you want 4 pairs of bools
	// 	States:       make(map[string]HRAElevState),
	// }

	commonState = HRAInput{
		Origin: "peer-10.22.229.227-22222",
		Seq:    0,
		Ackmap: map[string]Ack_status{
			"peer-10.22.229.227-22222": NotAcked,
			"peer-10.22.229.227-11111": NotAcked,
			"peer-10.22.229.227-33333": NotAcked,

		},
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		States: map[string]HRAElevState{
			"peer-10.22.229.227-22222": {
				Behaviour:   "idle",
				Floor:       0,
				Direction:   "up",
				CabRequests: []bool{false, false, false, true},
			},
			"peer-10.22.229.227-11111": {
				Behaviour:   "idle",
				Floor:       0,
				Direction:   "up",
				CabRequests: []bool{false, false, false, true},
			},
			"peer-10.22.229.227-33333": {
				Behaviour:   "idle",
				Floor:       0,
				Direction:   "up",
				CabRequests: []bool{false, false, false, true},
			},
		},
	}

	// commonState = HRAInput{
	// 	Origin:       config.Elevator_id,
	// 	ID:           1,
	// 	Ackmap:       make(map[string]Ack_status),
	// 	HallRequests: make([][2]bool, 4), // Assuming you want 4 pairs of bools
	// 	States: map[string]HRAElevState{
	// 		config.Elevator_id: {}, // Replace "initialKey" with your key
	// 	},
	// }

	go elevio.PollButtons(elevioOrdersC)
	//go Update_Assingments(elevioOrdersC, deliveredOrderC, newAssingemntC)

	heartbeatTimer := time.NewTicker(15 * time.Millisecond)

	for {

		select {
		case <-timeCounter.C:
			state = Isolated
		default:
		}

		switch state {
		case Idle:
			select {
			//case assingmentUpdate := <-newAssingemntC: //bufferes lage stor kanal 64 feks lage tømmefunksjon
			//	AssignmentStash = assingmentUpdate
			//	StashType = AssingmetChange
			//	commonState.Update_Assingments(assingmentUpdate)
			//	commonState.NullAckmap()
			//	commonState.Ack()
			//	//PrintCommonState(commonState)
			//	state = SendingSelf
			case newOrder := <-elevioOrdersC:
				fmt.Println("New order")
				NewOrderStash = newOrder
				StashType = AddCall
				//PrintCommonState(commonState)
				commonState.AddCall(newOrder)
				commonState.NullAckmap()
				commonState.Ack()
				PrintCommonState(commonState)
				state = SendingSelf


			case removeOrder := <-deliveredOrderC:
				fmt.Println("Delivered PULL OUT ")
				RemoveOrderStash = removeOrder
				StashType = RemoveCall
				commonState.removeCall(removeOrder)
				commonState.NullAckmap()
				commonState.Ack()
				state = SendingSelf
				
			case newElevState := <-newElevStateC: //bufferes lage stor kanal 64 feks
				fmt.Println("newElevState")
				StateStash = newElevState
				StashType = StateChange
				commonState.toHRAElevState(newElevState)
				commonState.NullAckmap()
				commonState.Ack()
				state = SendingSelf
				PrintCommonState(commonState)


			case arrivedCommonState := <-receiveFromNetworkC: //bufferes lage stor kanal 64 feks
				timeCounter = time.NewTimer(selfLostNetworkDuratio)
				//fmt.Println("Vebjøn liker tss")
				//PrintCommonState(commonState)
				//fmt.Println("Arrived")
				//PrintCommonState(arrivedCommonState)
				
				//arrivedCommonState.ensureElevatorState(arrivedCommonState.States[config.Elevator_id])

				switch {
				case higherPriority(commonState, arrivedCommonState):
					fmt.Println("something fishy")
					//if arrivedCommonState.Origin == config.Elevator_id {
					//state = SendingSelf
				
					arrivedCommonState.Ack()
					commonState = arrivedCommonState
					state = Acking
					//PrintCommonState(commonState)
					//}
					//if arrivedCommonState.Origin != config.Elevator_id {
					//	fmt.Println("arrived new commonstate")
					//	arrivedCommonState.Ack()
					//	commonState = arrivedCommonState
					//	state = Acking
					//}
				default:
					break //doing jack
				}
			case peers := <-recieveFromPeerC:
				fmt.Println(peers) //bufferes lage stor kanal 64 feks
				fmt.Println("    ")
				fmt.Println("peers number 1 fucked")
				fmt.Println("    ")
				commonState.makeElevUnav(peers)
			default:
			}
		case SendingSelf:
			//fmt.Println("-")
			select {
			case arrivedCommonState := <-receiveFromNetworkC:
				//fmt.Println("Im in SendingSelf mode")
				timeCounter = time.NewTimer(selfLostNetworkDuratio)
				switch {
				case arrivedCommonState.Origin != config.Elevator_id && higherPriority(commonState, arrivedCommonState):
					fmt.Println("I am not priority:(")
					arrivedCommonState.Ack()
					commonState = arrivedCommonState
					state = AckingOtherWhileTryingToSendSelf

				case Fully_acked(arrivedCommonState.Ackmap):
					//fmt.Println("get in there")
					state = Idle
					//fmt.Println("Fucking get in there")
					commonState = arrivedCommonState
					messageToAssinger <- commonState
				default:
					//fmt.Println("doing jack")
					//break //doing jack
					//fmt.Println("Priority mofo")
				}

			case peers := <-recieveFromPeerC: //bufferes lage stor kanal 64 feks
				commonState.makeElevUnav(peers)
				fmt.Println(peers)
				fmt.Println("    ")
				fmt.Println("peers number 2 fucked")
				fmt.Println("    ")
				if Fully_acked(commonState.Ackmap) {
					state = Idle
					messageToAssinger <- commonState
				}
			default:
			}

		case Acking:
			select {
			case arrivedCommonState := <-receiveFromNetworkC:
				//fmt.Println("Im in acking mode")
				//PrintCommonState(arrivedCommonState)

				timeCounter = time.NewTimer(selfLostNetworkDuratio)
				switch {

				case Fully_acked(arrivedCommonState.Ackmap):
					state = Idle
					commonState = arrivedCommonState
					messageToAssinger <- commonState

				case higherPriority(commonState, arrivedCommonState): // && takePriortisedCommonState(commonState, arrivedCommonState) priority of higher  {
					arrivedCommonState.Ack()
					commonState = arrivedCommonState
					
				default: 
					arrivedCommonState.Ack()
					commonState = arrivedCommonState
				//
				}

			case peers := <-recieveFromPeerC:
				commonState.makeElevUnav(peers)
				fmt.Println(peers)
				fmt.Println("    ")
				fmt.Println("peers number 3 fucked")
				fmt.Println("    ")

				if Fully_acked(commonState.Ackmap) {
					state = Idle
					messageToAssinger <- commonState
				}

			default:
			}

		case AckingOtherWhileTryingToSendSelf:
			select {
			case arrivedCommonState := <-receiveFromNetworkC:
				//fmt.Println("AckingOtherWhileTryingToSendSelf")
				timeCounter = time.NewTimer(selfLostNetworkDuratio)
				switch {
				//case !higherPriority(commonState, arrivedCommonState):
				//	break //doing jack

				case higherPriority(commonState, arrivedCommonState): // && takePriortisedCommonState(commonState, arrivedCommonState) priority of higher  {
					arrivedCommonState.Ack()
					commonState = arrivedCommonState
					fmt.Println("Esssskeetit")
					PrintCommonState(commonState)

				case Fully_acked(arrivedCommonState.Ackmap):
					state = SendingSelf
					//fmt.Println("BOOOOOOOB")
					switch StashType {

					case AddCall:
						arrivedCommonState.AddCall(NewOrderStash)

					case RemoveCall:
						arrivedCommonState.removeCall(RemoveOrderStash)

					case StateChange:
						arrivedCommonState.toHRAElevState(StateStash)
						//fmt.Println("statechange")
					}

					//fmt.Println("TTTTTTTTTTTT")
					commonState = arrivedCommonState
					//fmt.Println("AAAAAAAAAAAAAA")
					messageToAssinger <- commonState
					//fmt.Println("BBBBBBBBBBBBBB")
					commonState.NullAckmap()
					commonState.Ack()
					//PrintCommonState(commonState)
				default:
					arrivedCommonState.Ack()
					//fmt.Println("suck a big ooooooone")
					commonState = arrivedCommonState

				}
			case peers := <-recieveFromPeerC:

				commonState.makeElevUnav(peers)
				fmt.Println(peers)
				fmt.Println("    ")
				fmt.Println("peers number 4 fucked")
				fmt.Println("    ")
				if Fully_acked(commonState.Ackmap) {
					state = SendingSelf
					messageToAssinger <- commonState
				}

			default:
			}
		case Isolated:
			fmt.Println("I should not be here???")
			select {
			//case <- recieveFromPeerC:
			//	state = Idle

			case <-receiveFromNetworkC:
				state = Idle

			case assingmentUpdate := <-newAssingemntC: //bufferes lage stor kanal 64 feks lage tømmefunksjon
				commonState.makeElevUnavExceptOrigin()
				commonState.UpdateCabAssignments(assingmentUpdate)
				messageToAssinger <- commonState

			case newElevState := <-newElevStateC: //bufferes lage stor kanal 64 feks
				commonState.toHRAElevState(newElevState)
				commonState.makeElevUnavExceptOrigin()
				messageToAssinger <- commonState

			default:
			}

			//case UnableToMove: // TODO: make channel for unav elevator
			//	select{
			//		case AbleToMove := <-newElevStateC:
			//			state = Idle
			//	default:
			//		commonState.makeOriginElevUnav()
			//
		}

		select {
		case <-heartbeatTimer.C:
			giverToNetwork <- commonState
		default:
		}

	} // to do: add case when for elevator lost network connection
}
