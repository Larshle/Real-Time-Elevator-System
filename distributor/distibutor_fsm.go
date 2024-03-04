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

	fmt.Print("Distributor started\n")

	elevioOrdersC := make(chan elevio.ButtonEvent)
	newAssingemntC := make(chan localAssignments)
	peerUpdateC := make(chan peers.PeerUpdate)
	var localAssignments localAssignments
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


	fmt.Println("Første commonstate")
	PrintCommonState(commonState)

	go elevio.PollButtons(elevioOrdersC)
	go Update_Assingments(elevioOrdersC, deliveredOrderC, newAssingemntC)



	for {
		select {
		case assingmentUpdate := <-newAssingemntC:
			localAssignments.Update_Assingments(assingmentUpdate)

		case newElevState := <-newElevStateC:
			localCommonState.Update_local_state(newElevState)

		case peers := <-peerUpdateC:
			switch {
			case len(peers.Lost) != 0:
				P = peers


		case arrivedCommonState := <-receiveFromNetworkC:
			switch {
			case Fully_acked(arrivedCommonState.Ackmap):
				fmt.Println("Sjekke liit opp")
				commonState = arrivedCommonState
				messageToAssinger <- commonState
				localCommonState.MergeCommonState(commonState, localAssignments)
				giverToNetwork <- localCommonState

			case commonStatesNotEqual(commonState, arrivedCommonState):
				commonState = takePriortisedCommonState(commonState, arrivedCommonState)
				commonState.Ack()
				giverToNetwork <- commonState

			default:
				commonState = arrivedCommonState
				commonState.makeElevUnav(P)
				commonState.Ack()
				giverToNetwork <- commonState
			}
		


		default:
			if Fully_acked(commonState.Ackmap) {
				localCommonState.MergeCommonState(commonState, localAssignments)
				giverToNetwork <- localCommonState
				}
			}

		} // to do: add case when for elevator lost network connection

	}
}