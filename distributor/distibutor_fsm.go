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
	var peers peers.PeerUpdate

	commonState := HRAInput{
		Origin:       config.Elevator_id,
		ID:           1,
		Ackmap:       make(map[string]Ack_status),
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		States: map[string]HRAElevState{
			config.Elevator_id: {
				Behaviour:   "idle",
				Floor:       2,
				Direction:   "stop",
				CabRequests: []bool{false, false, false, true},
			},
		},
	}

	fmt.Println("Første commonstate")
	PrintCommonState(commonState)

	go elevio.PollButtons(elevioOrdersC)
	go Update_Assingments(elevioOrdersC, deliveredOrderC, newAssingemntC)

	for {
		select {
		case localAssignments = <-newAssingemntC:
			temp := commonState
			localstate := commonState.States[config.Elevator_id]
			commonState.Update_Assingments(localAssignments)
			commonState = temp
			commonState.States[config.Elevator_id] = localstate
			giverToNetwork <- commonState

		case newElevState := <-newElevStateC:
			temp := commonState
			localstate := commonState.States[config.Elevator_id]
			commonState.Update_local_state(newElevState)
			commonState = temp
			commonState.States[config.Elevator_id] = localstate
			giverToNetwork <- commonState

		case peers = <-peerUpdateC:
			commonState.makeElevUnav(peers)
			commonState.Origin = config.Elevator_id
			giverToNetwork <- commonState

		case receivedCommonState := <-receiveFromNetworkC:
			switch {
			case Fully_acked(receivedCommonState.Ackmap):
				commonState = receivedCommonState
				commonState.Origin = config.Elevator_id
				messageToAssinger <- receivedCommonState

			default:
				temp := takePriortisedCommonState(commonState, receivedCommonState)
				localstate := commonState.States[config.Elevator_id]
				commonState = temp
				commonState.States[config.Elevator_id] = localstate
				commonState.Ack()
				commonState.ID++
				commonState.Origin = config.Elevator_id
				giverToNetwork <- commonState

			}
		}
	}

}
