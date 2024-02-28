package distributor

import (
	"root/driver/elevio"
	"root/network/network_modules/peers"
	"root/assigner"
	"root/elevator"
	"root/network/network_modules/peers"
	
)

func Distributor_fsm(
	deliveredOrderC <-chan elevio.ButtonEvent, 
	newElevStateC <-chan elevator.State, 
	giverToNetwork chan<- Commonstate, 
	receiveFromNetworkC <-chan Commonstate,
	messageToAssinger chan<- assigner.HRAInput) {

	elevioOrdersC := make(chan elevio.ButtonEvent)
	newAssingemntC := make(chan localAssignments)
	peerUpdateC := make(chan peers.PeerUpdate)
	var localAssignments localAssignments
	var state elevator.State

	// Initialize the distributor
	var localAssignments = localAssignments
	var commonState = assigner.HRAInput
	var elevatorID = network.Generate_ID()
	
	go elevio.PollButtons(elevioOrdersC)
	go Update_Assingments(elevioOrdersC, deliveredOrderC, newAssingemntC)
	go peers.Receiver(15647, peerUpdateC)
	go network.

	for{
		select{
		case localAssignments := <- newAssingemntC:
			commonState = commonState.Update_Assingments(localAssignments, elevatorID)
			giverToNetwork <- commonState

		case newElevState := <- newElevStateC:
			commonState = commonState.Update_ElevState(newElevState, elevatorID)
			giverToNetwork <- commonState
	
		}
		case peerUpdate := <- peerUpdateC:
			switch{
				case peer.Update.Lost != 0:
				//Blalal	
					giverToNetwork <- commonState
				}



		
		case receivedCommonState := <- receiveFromNetworkC:
			switch{
				case fullyAcked(receivedCommonState, localAssignments):
					messageToAssinger <- receivedCommonState
				}
				case !fullyAcked(receivedCommonState, localAssignments):



	}

}
	// This is the main function for the distributor FSM. It listens to the channels for new orders, delivered orders, new elevator states, delivered elevator states, peer updates, and messages from the network. It then updates the state of the distributor, and sends messages to the network if necessary.

	func Recieve_commonstate(new_commonstate assigner.HRAInput, cToAssingerC chan <- assigner.HRAInput) {

		if Commonstates_are_equal(new_commonstate, Unacked_Commonstate) {
			return
		}
	
		if Fully_acked(new_commonstate.Ackmap) {
			Unacked_Commonstate = new_commonstate // vet ikke om dette er nødvendig
			Commonstate = new_commonstate
			// broadcast (gjøres hele tiden fra main)
			cToAssingerC <- Commonstate
		}
	
		// if new_commonstate har lavere prioritet
		// return
		if new_commonstate.ID < Unacked_Commonstate.ID || id_is_lower(new_commonstate.Origin, Unacked_Commonstate.Origin) {
			return
		}
	
		// else
		// ack, oppdater ack_commonstate og broadcast denne helt til den er acket eller det kommer en ny med høyere prioritet
		new_commonstate.Ackmap[Elevator_id] = assigner.Acked
		Unacked_Commonstate = new_commonstate
	
	}
	
	func Create_commonstate(peerUpdateCh <-chan peers.PeerUpdate) {
		// Create a new commonstate
		Unacked_Commonstate.ID++
		Unacked_Commonstate.Origin = Elevator_id
		Unacked_Commonstate.Ackmap = make(map[string]assigner.Ack_status)
		Unacked_Commonstate.Ackmap[Elevator_id] = 1
	}