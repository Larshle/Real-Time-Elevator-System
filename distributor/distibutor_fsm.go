package distributor

import (
	"fmt"
	"root/config"
	"root/driver/elevio"
	"root/elevator"
	"root/network/network_modules/peers"
	"time"
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
	checkNettworkTimer := time.NewTimer(time.Hour)

	var commonState HRAInput
	var localCommonState HRAInput
	var localAssignments localAssignments
	var P peers.PeerUpdate


	// commonState = HRAInput{
	// 	Origin:       config.Elevator_id,
	// 	ID:           0,
	// 	Ackmap:       make(map[string]Ack_status),
	// 	HallRequests: make([][2]bool, 4), // Assuming you want 4 pairs of bools
	// 	States:       make(map[string]HRAElevState),
	// }

	commonState = HRAInput{
		Origin:       config.Elevator_id,
		ID:           0,
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
	
	localCommonState = HRAInput{
		Origin:       config.Elevator_id,
		ID:           0,
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

	queue := &CommonStateQueue{}

	go elevio.PollButtons(elevioOrdersC)
	go Update_Assingments(elevioOrdersC, deliveredOrderC, newAssingemntC)

	heartbeatTimer := time.NewTicker(100 * time.Millisecond)

	for {
		select {
			case assingmentUpdate := <-newAssingemntC:
				localAssignments.Update_Assingments(assingmentUpdate)
				queue.Enqueue(commonState)

			case newElevState := <-newElevStateC:
				localCommonState.Update_local_state(newElevState)
				

			case peers := <-peerUpdateC:
				P = peers

			case arrivedCommonState := <-receiveFromNetworkC:
				fmt.Println("receiveFromNetworkC")
				checkNettworkTimer = time.NewTimer(500*time.Millisecond)
				switch {
					case Fully_acked(arrivedCommonState.Ackmap):
						commonState = arrivedCommonState
						messageToAssinger <- commonState
						

					case commonStatesNotEqual(commonState, arrivedCommonState):
						commonState = takePriortisedCommonState(commonState, arrivedCommonState)
						commonState.Ack()


					default:
						commonState = arrivedCommonState
						commonState.makeElevUnav(P)
						commonState.Ack()
				}
		
			case <-heartbeatTimer.C:
				switch{
				case Fully_acked(commonState.Ackmap):
					localCommonState.MergeCommonState(commonState, localAssignments)
					giverToNetwork <- localCommonState
					
				default:
					giverToNetwork <- commonState
				}
			
			// case <-checkNettworkTimer.C:
			// 	localCommonState.MergeCommonState(localCommonState, localAssignments)
			// 	giverToNetwork <- localCommonState
			// }

		} // to do: add case when for elevator lost network connection
	}
}