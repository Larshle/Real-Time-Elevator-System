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
	giverToNetwork chan<- HRAInput2,
	receiveFromNetworkC <-chan HRAInput2,
	messageToAssinger chan<- HRAInput2) {

	elevioOrdersC := make(chan elevio.ButtonEvent)
	newAssingemntC := make(chan localAssignments)
	peerUpdateC := make(chan peers.PeerUpdate)
	var peers peers.PeerUpdate

	// commonState = HRAInput2{
	// 	Origin:       config.Elevator_id,
	// 	ID:           0,
	// 	Ackmap:       make(map[string]Ack_status),
	// 	HallRequests: make([][2]bool, 4), // Assuming you want 4 pairs of bools
	// 	States:       make(map[string]HRAElevState),
	// }

	commonState := HRAInput2{
		Origin:       config.Elevator_id,
		ID:           0,
		Ackmap:       make(map[string]Ack_status),
		HallRequests: [][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
		States: map[string]HRAElevState{
			config.Elevator_id: {
				Behaviour:   "idle",
				Floor:       0,
				Direction:   "up",
				CabRequests: []bool{false, false, false, true},
			},
		},
	}

	// localCommonState = HRAInput2{
	// 	Origin:       config.Elevator_id,
	// 	ID:           0,
	// 	Ackmap:       make(map[string]Ack_status),
	// 	HallRequests: [][2]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}},
	// 	States: map[string]HRAElevState{
	// 		config.Elevator_id: {
	// 			Behaviour:   "idle",
	// 			Floor:       0,
	// 			Direction:   "up",
	// 			CabRequests: []bool{false, false, false, true},
	// 		},
	// 	},
	// }

	// commonState = HRAInput2{
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

	ticker := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-ticker:
			giverToNetwork <- commonState
			// fmt.Println("Distributor: Sent commonstate")
		case assingmentUpdate := <-newAssingemntC:
			// localCs.Update_Assingments(assingmentUpdate)
			// commonState.MergeCommonState(localCs)
			commonState.Update_Assingments(assingmentUpdate)

			// fmt.Println("Distributor: Updated assingments")
			// PrintCommonState(commonState)
			// giverToNetwork	<- commonState

		case newElevState := <-newElevStateC:
			// localCs.Update_local_state(newElevState)
			// commonState.MergeCommonState(localCs)
			commonState.Update_local_state(newElevState)

			// giverToNetwork	<- commonState

		case peers = <-peerUpdateC:
			// localCs.makeElevUnav(peers)
			// commonState.MergeCommonState(localCs)
			commonState.makeElevUnav(peers)

			// giverToNetwork	<- commonState

		case arrivedCommonState := <-receiveFromNetworkC:
			switch {
			case 
			case Fully_acked(arrivedCommonState.Ackmap):
			}
				fmt.Println("Distributor: Fully acked")
				PrintCommonState(arrivedCommonState)

				messageToAssinger <- arrivedCommonState
				// commonState.MergeCommonState(arrivedCommonState)

				for key := range commonState.Ackmap {
					commonState.Ackmap[key] = NotAcked
				}
				commonState.Ack()
				for i := range commonState.HallRequests {
					for j := range commonState.HallRequests[i] {
						if commonState.HallRequests[i][j] == 2 {
							commonState.HallRequests[i][j] = 0
						}
					}
				}
				commonState.Origin = config.Elevator_id
				commonState.ID++

				// giverToNetwork	<- commonState

				// til assigner
				// øke id på commonstate
				// tømme ackmap
				// oppdater commonstate med dine lokale endringer
				// gjøre 2 i hallrequest til 0
				// ack
				// broadcast

			//case IsEqual(commonState, arrivedCommonState):
			//	continue

			default:
				commonState = takePriortisedCommonState(commonState, arrivedCommonState)
				commonState.MergeCommonState(arrivedCommonState)
				// giverToNetwork	<- commonState
			}
		}

	} // to do: add case when for elevator lost network connection
}
