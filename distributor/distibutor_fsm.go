package distributor

import (
	"root/driver/elevio"
	"root/network/network_modules/peers"
	"root/network/network_modules/bcast"
	"root/elevator"
	"root/network"
)

var Elevator_id string

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
	Elevator_id = network.Generate_ID()
	
	go elevio.PollButtons(elevioOrdersC)
	go Update_Assingments(elevioOrdersC, deliveredOrderC, newAssingemntC)
	go peers.Receiver(15647, peerUpdateC)
	go bcast.Transmitter(15647, Elevator_id, commonState) // MÅ ENDRES

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
						receivedCommonState.Ackmap[Elevator_id] = Acked
						commonState = receivedCommonState
					}
		}
	}

}