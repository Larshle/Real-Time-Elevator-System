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
				Floor:       2,
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

	giverToNetwork <- commonState

	queue := &CommonStateQueue{}

	for {
		select {
		case localAssignments = <-newAssingemntC:
			fmt.Println("LOCAL ASSIGNMENTS:")
			commonState.Update_Assingments(localAssignments)
			giverToNetwork <- commonState
			queue.Enqueue(commonState)

		case newElevState := <-newElevStateC:
			commonState.Update_local_state(newElevState)
			giverToNetwork <- commonState
			queue.Enqueue(commonState)

		case peers := <-peerUpdateC: 
			if len(peers.Lost) != 0{
				commonState.Update_ackmap(peers)
				commonState.ID++
				giverToNetwork <- commonState
			}
		


		case arrivedCommonState := <-receiveFromNetworkC:
			switch {
			case Fully_acked(arrivedCommonState.Ackmap):
				fmt.Println("Sjekke liit opp")
				commonState = arrivedCommonState
				messageToAssinger <- commonState

			case Higher_priority(commonState, commonState):
				commonState.Ack()
				giverToNetwork <- commonState

			case commonState.ID < arrivedCommonState.ID:
				commonState = HighestIDState(commonState, arrivedCommonState)
				commonState.Ack()
				giverToNetwork <- commonState

			default:
				commonState = arrivedCommonState
				commonState.Ack()
				giverToNetwork <- commonState
			}
		default:

			if newcommonState, ok := queue.Dequeue(); ok && Fully_acked(commonState.Ackmap) {
				newcommonState.ID = commonState.ID + 1
				giverToNetwork <- newcommonState
			}
		}

	} 

}
