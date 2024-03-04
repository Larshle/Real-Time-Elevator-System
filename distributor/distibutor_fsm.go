package distributor

import (
	"fmt"
	"root/config"
	"root/driver/elevio"
	"root/elevator"
	"root/network/network_modules/peers"
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
	var localAssignments localAssignments
	var commonState HRAInput

	// commonState = HRAInput{
	// 	Origin:       config.Elevator_id,
	// 	ID:           0,
	// 	Ackmap:       make(map[string]Ack_status),
	// 	HallRequests: make([][2]bool, 4), // Assuming you want 4 pairs of bools
	// 	States:       make(map[string]HRAElevState),
	// }

	commonState = HRAInput{
		Origin:       config.Elevator_id,
		ID:           1,
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

	for {
		select {
		case localAssignments = <-newAssingemntC:
			commonState.Update_Assingments(localAssignments)
			giverToNetwork <- commonState

		case newElevState := <-newElevStateC:
			commonState.Update_local_state(newElevState)
			giverToNetwork <- commonState

		case peers := <-peerUpdateC:
			switch {
			case peers.New != "":
				giverToNetwork <- commonState
			}

		case receivedCommonState := <-receiveFromNetworkC:
			switch {
			case Fully_acked(receivedCommonState.Ackmap):
				messageToAssinger <- receivedCommonState

			case Higher_priority(receivedCommonState, commonState):
				commonState = receivedCommonState

			default:
				receivedCommonState.Ackmap[config.Elevator_id] = Acked
				commonState = receivedCommonState
			}
		}
	}

}
