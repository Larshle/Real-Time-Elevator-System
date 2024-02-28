package distributor

import (
	"root/config"
	"root/driver/elevio"
	"root/network/network_modules/peers"
	"root/network/network_modules/bcast"
	"root/elevator"
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
	
	go elevio.PollButtons(elevioOrdersC)
	go Update_Assingments(elevioOrdersC, deliveredOrderC, newAssingemntC)
	go peers.Receiver(15647, peerUpdateC)
	go bcast.Transmitter(15647, commonState) // MÅ ENDRES

	for{
		select{
			case localAssignments = <- newAssingemntC:
				commonState.Update_Assingments(localAssignments)
				giverToNetwork <- commonState

			case newElevState := <- newElevStateC:
				commonState.Update_local_state(newElevState)
				giverToNetwork <- commonState
			
			case peers := <- peerUpdateC:
				switch{
					case peers.New != "":
					giverToNetwork <- commonState
				}
			
			case receivedCommonState := <- receiveFromNetworkC:
				switch{
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