package distributor

import (
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
	AckingOtherWhileTryingToSend
	Isolated
	
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
	var lastesSelfState elevator.State
	var stash localAssignments
	var state State = Idle
	
	


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

		switch state {
			case Idle: 
				select {
					case assingmentUpdate := <-newAssingemntC: //bufferes lage stor kanal 64 feks lage tømmefunksjon 
						stash = assingmentUpdate
						commonState.Update_Assingments(assingmentUpdate)
						state = SendingSelf

					case newElevState := <-newElevStateC: //bufferes lage stor kanal 64 feks
						lastesSelfState = newElevState
						commonState.toHRAElevState(newElevState)
						state = SendingSelf
						

					case arrivedCommonState := <-receiveFromNetworkC: //bufferes lage stor kanal 64 feks
						arrivedCommonState.ensureElevatorState(arrivedCommonState.States[config.Elevator_id])
						if arrivedCommonState.Origin == config.Elevator_id {
							state = SendingSelf
						}
						if arrivedCommonState.Origin != config.Elevator_id {
							arrivedCommonState.Ack()
							commonState = arrivedCommonState
							state = Acking
						}
					case peers := <- peerUpdateC: //bufferes lage stor kanal 64 feks
						commonState.makeElevUnav(peers)
					default:
				}
			case SendingSelf:
				select {
				case arrivedCommonState := <-receiveFromNetworkC:
					if arrivedCommonState.Origin != config.Elevator_id && arrivedCommonState.seq >= commonState.seq && 
					
					
					arrivedCommonState.Origin {
						arrivedCommonState.Ack()
						commonState = arrivedCommonState
						state = AckingOtherWhileTryingToSend
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
					if arrivedCommonState.seq >= commonState.seq{ // && takePriortisedCommonState(commonState, arrivedCommonState) priority of higher  {
						arrivedCommonState.Ack()
						commonState = arrivedCommonState
					}
					if Fully_acked(commonState.Ackmap){
						state = Idle
						messageToAssinger <- commonState
					}

				case peers := <- peerUpdateC:
					commonState.makeElevUnav(peers)
					if Fully_acked(commonState.Ackmap){
						state = Idle
						messageToAssinger <- commonState
					}
				default:
			}
				
			case AckingOtherWhileTryingToSend:
				select {
				case arrivedCommonState := <-receiveFromNetworkC:
					if arrivedCommonState.seq < commonState.seq{
						break
					} 
					if arrivedCommonState.seq >= commonState.seq{ // && takePriortisedCommonState(commonState, arrivedCommonState) priority of higher  {
						arrivedCommonState.Ack()
						commonState = arrivedCommonState
					}
					if Fully_acked(commonState.Ackmap){
						state = SendingSelf
						commonState.Update_Assingments(stash)
						commonState.toHRAElevState(lastesSelfState)
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
				case peers := <- peerUpdateC:
					commonState.makeElevUnav(peers)
					state = Idle
				case arrivedCommonState := <-receiveFromNetworkC:
					commonState = arrivedCommonState
					arrivedCommonState.Ack()
					state = Acking

				case assingmentUpdate := <-newAssingemntC: //bufferes lage stor kanal 64 feks lage tømmefunksjon 
						stash = assingmentUpdate
						commonState.Update_Assingments(assingmentUpdate)
						state = SendingSelf

				case newElevState := <-newElevStateC: //bufferes lage stor kanal 64 feks
					lastesSelfState = newElevState
					commonState.toHRAElevState(newElevState)
					state = SendingSelf

				default:
				}


					

	
			select {
			case <-heartbeatTimer.C:
				giverToNetwork <- commonState
			default:
				}
				
		}
						// }

	} // to do: add case when for elevator lost network connection
			
}
	
