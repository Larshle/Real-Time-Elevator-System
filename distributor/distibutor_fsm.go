package distributor

import (
	"root/driver/elevio"
	"root/network/network_modules/peers"
	"root/assigner"
	"root/elevator"
)



func Distributor_fsm(
	deliveredOrderC <-chan elevio.ButtonEvent, 
	newElevStateC <-chan elevator.State,  
	peerUpdateC <-chan peers.PeerUpdate, 
	giverToNetwork chan<- Commonstate, 
	receiveFromNetworkC <-chan Commonstate,
	messageToAssinger chan<- assigner.HRAInput,) {

	elevioOrdersC := make(chan elevio.ButtonEvent)
	newAssingemntC := make(chan localAssignments)
	var localAssignments localAssignments
	var state elevator.State

	// Initialize the distributor
	var localAssignments = localAssignments
	var commonState = Commonstate
	
	go elevio.PollButtons(elevioOrdersC)
	go Update_Assingments(elevioOrdersC, deliveredOrderC, newAssingemntC)

	for{
		select{
		case localAssignments := <- newAssingemntC:
			switch{
			case localAssignments.localCabAssignments != commonState.HallRequests:
				commonState.HallRequests = localAssignments.localHallAssignments
				giverToNetwork <- commonState
			}

		case newElevState := <- newElevStateC:
			switch{

			}
		}
		case peerUpdate := <- peerUpdateC:
			switch{
			case peerUpdate.New != "":
				giverToNetwork <- commonState
			}
		
		case receivedCommonState := <- receiveFromNetworkC:
			switch{
				case fullyAcked(receivedCommonState, localAssignments):
					messageToAssinger <- receivedCommonState
				}	


	}

}
	// This is the main function for the distributor FSM. It listens to the channels for new orders, delivered orders, new elevator states, delivered elevator states, peer updates, and messages from the network. It then updates the state of the distributor, and sends messages to the network if necessary.