package distributor

import (
	"root/config"
	"root/driver/elevio"
	"root/elevator"
	"root/network/network_modules/peers"
	"time"
	"fmt"

	
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
	newAssingemntC := make(chan localAssignments)
	peerUpdateC := make(chan peers.PeerUpdate)

	var commonState HRAInput
	var StateStash elevator.State
	var AssignmentStash localAssignments
	var state State = Idle
	var StashType StatshType
	timeCounter := time.NewTimer(time.Hour)
	selfLostNetworkDuratio := 1 * time.Second
	


	// commonState = HRAInput{
	// 	Origin:       config.Elevator_id,
	// 	ID:           0,
	// 	Ackmap:       make(map[string]Ack_status),
	// 	HallRequests: make([][2]bool, 4), // Assuming you want 4 pairs of bools
	// 	States:       make(map[string]HRAElevState),
	// }

	commonState = HRAInput{
		Origin:       config.Elevator_id,
		seq:           0,
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
	go Update_Assingments(elevioOrdersC, deliveredOrderC, newAssingemntC)

	heartbeatTimer := time.NewTicker(15 * time.Millisecond)

	for {

		select{
			case <- timeCounter.C:
				state = Isolated
			default:
		}

		switch state {
			case Idle: 
				select {
					case assingmentUpdate := <-newAssingemntC: //bufferes lage stor kanal 64 feks lage tømmefunksjon 
						AssignmentStash = assingmentUpdate
						StashType = AssingmetChange
						commonState.Update_Assingments(assingmentUpdate)
						state = SendingSelf

					case newElevState := <-newElevStateC: //bufferes lage stor kanal 64 feks
						StateStash = newElevState
						StashType = StateChange
						commonState.toHRAElevState(newElevState)
						state = SendingSelf
						

					case arrivedCommonState := <-receiveFromNetworkC://bufferes lage stor kanal 64 feks
						timeCounter = time.NewTimer(selfLostNetworkDuratio) 
						arrivedCommonState.ensureElevatorState(arrivedCommonState.States[config.Elevator_id])

						switch {
							case higherPriority(commonState, arrivedCommonState):
								fmt.Println("something fishy")
								//if arrivedCommonState.Origin == config.Elevator_id {
								//state = SendingSelf
							//}
							if arrivedCommonState.Origin != config.Elevator_id {
								fmt.Println("arrived new commonstate")
								arrivedCommonState.Ack()
								commonState = arrivedCommonState
								state = Acking
							}
							default:
								break //doing jack
						}
					case peers := <- peerUpdateC: //bufferes lage stor kanal 64 feks
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

				case peers := <- peerUpdateC: //bufferes lage stor kanal 64 feks
					commonState.makeElevUnav(peers)
					if Fully_acked(commonState.Ackmap){
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

					case higherPriority(commonState, arrivedCommonState):// && takePriortisedCommonState(commonState, arrivedCommonState) priority of higher  {
						arrivedCommonState.Ack()
						commonState = arrivedCommonState
				
					case !higherPriority(commonState, arrivedCommonState):
						break //doing jack
					}

				case peers := <- peerUpdateC:
					commonState.makeElevUnav(peers)
					if Fully_acked(commonState.Ackmap){
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
						break //doing jack

					case higherPriority(commonState, arrivedCommonState):// && takePriortisedCommonState(commonState, arrivedCommonState) priority of higher  {
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
						messageToAssinger <- commonState
					}
				case peers := <- peerUpdateC:
					commonState.makeElevUnav(peers)
					if Fully_acked(commonState.Ackmap){
						state = SendingSelf
						messageToAssinger <- commonState
					}
				
				default:
			}
			case Isolated:
				select{
				//case <- peerUpdateC:
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

	
