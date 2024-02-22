package distributor

import (
	"root/driver/elevio"
	"root/network/network_modules/peers"
	"root/assigner"
	"root/elevator"
)



func Distributor_fsm(newOrderC <-chan elevio.ButtonEvent, 
	deliveredOrderC <-chan elevio.ButtonEvent, 
	newElevStateC <-chan elevator.State,  
	peerUpdateC <-chan peers.PeerUpdate, 
	giverToNetwork chan<- Commonstate, 
	receiveFromNetworkC <-chan Coomonstate) {

}
	// This is the main function for the distributor FSM. It listens to the channels for new orders, delivered orders, new elevator states, delivered elevator states, peer updates, and messages from the network. It then updates the state of the distributor, and sends messages to the network if necessary.