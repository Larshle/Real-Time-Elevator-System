package distributor

import (
	"fmt"
	"root/config"
	"root/driver/elevio"
	"root/elevator"
	"root/network/network_modules/peers"
	"time"
)

type StatshType int

const (
	RemoveCall StatshType = iota
	AddCall
	StateChange
)


func Distributor(
	deliveredOrderC <-chan elevio.ButtonEvent,
	newElevStateC <-chan elevator.State,
	giverToNetwork chan<- HRAInput,
	receiveFromNetworkC <-chan HRAInput,
	messageToAssinger chan<- HRAInput, 
	recieveFromPeerC <- chan peers.PeerUpdate) {

	elevioOrdersC := make(chan elevio.ButtonEvent, 10000)
	newAssingemntC := make(chan localAssignments, 10000)

	var commonState HRAInput
	var StateStash elevator.State
	var NewOrderStash elevio.ButtonEvent
	var RemoveOrderStash elevio.ButtonEvent
	var StashType StatshType
	timeCounter := time.NewTimer(time.Second * 5)
	selfLostNetworkDuratio := 10 * time.Second
	stashed := false
	acking := false
	isolated:= false

	commonState = HRAInput{
		Origin: "peer-10.22.229.227-22222",
		Seq:    0,
		Ackmap: map[string]Ack_status{
			"peer-10.22.229.227-22222": NotAcked,
			"peer-10.22.229.227-11111": NotAcked,
			"peer-10.22.229.227-33333": NotAcked,

		},
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		States: map[string]HRAElevState{
			"peer-10.22.229.227-22222": {
				Behaviour:   "idle",
				Floor:       0,
				Direction:   "up",
				CabRequests: []bool{false, false, false, true},
			},
			"peer-10.22.229.227-11111": {
				Behaviour:   "idle",
				Floor:       0,
				Direction:   "up",
				CabRequests: []bool{false, false, false, true},
			},
			"peer-10.22.229.227-33333": {
				Behaviour:   "idle",
				Floor:       0,
				Direction:   "up",
				CabRequests: []bool{false, false, false, true},
			},
		},
	}

	go elevio.PollButtons(elevioOrdersC)

	heartbeatTimer := time.NewTicker(15 * time.Millisecond)

	for {

		select {
		case <-timeCounter.C:
			isolated = true
		default:
		}

		switch {
		case !stashed && !acking:
			select {

			case newOrder := <-elevioOrdersC:
				fmt.Println("New order")
				NewOrderStash = newOrder
				StashType = AddCall
				//PrintCommonState(commonState)
				commonState.AddCall(newOrder)
				commonState.NullAckmap()
				commonState.Ack()
				PrintCommonState(commonState)
				stashed = true
				acking = true

			case removeOrder := <-deliveredOrderC:
				fmt.Println("Delivered PULL OUT ")
				RemoveOrderStash = removeOrder
				StashType = RemoveCall
				commonState.removeCall(removeOrder)
				commonState.NullAckmap()
				commonState.Ack()
				stashed = true
				acking = true
				
			case newElevState := <-newElevStateC: //bufferes lage stor kanal 64 feks
				fmt.Println("newElevState")
				StateStash = newElevState
				StashType = StateChange
				commonState.toHRAElevState(newElevState)
				commonState.NullAckmap()
				commonState.Ack()
				stashed = true
				acking = true
				//PrintCommonState(commonState)

			case arrivedCommonState := <-receiveFromNetworkC: //bufferes lage stor kanal 64 feks
				timeCounter = time.NewTimer(selfLostNetworkDuratio)

				switch {
				case higherPriority(commonState, arrivedCommonState):
					fmt.Println("something fishy")
					arrivedCommonState.Ack()
					commonState = arrivedCommonState
					acking = true
				}
			case peers := <-recieveFromPeerC:
				commonState.makeElevUnav(peers)
				commonState.makeElevav()
			default:
			}

		case isolated:
			fmt.Println("I should not be here???")
			select {
			case <-receiveFromNetworkC:
				isolated = false

			case assingmentUpdate := <-newAssingemntC: //bufferes lage stor kanal 64 feks lage tømmefunksjon
				commonState.makeElevUnavExceptOrigin()
				commonState.UpdateCabAssignments(assingmentUpdate)
				messageToAssinger <- commonState

			case newElevState := <-newElevStateC: //bufferes lage stor kanal 64 feks
				commonState.toHRAElevState(newElevState)
				commonState.makeElevUnavExceptOrigin()
				messageToAssinger <- commonState

			default:
			}
		
		default:
			select {
			case arrivedCommonState := <-receiveFromNetworkC:
				if arrivedCommonState.Seq < commonState.Seq{
					break
				}
				timeCounter = time.NewTimer(selfLostNetworkDuratio)
				

				switch {
				case (arrivedCommonState.Origin > commonState.Origin && arrivedCommonState.Seq == commonState.Seq) || arrivedCommonState.Seq > commonState.Seq:
					arrivedCommonState.Ack()
					commonState = arrivedCommonState
				
				case Fully_acked(arrivedCommonState.Ackmap):	
					commonState = arrivedCommonState
					messageToAssinger <- commonState
					PrintCommonState(commonState)
					switch{
					case commonState.Origin != config.Elevator_id  && stashed:
						switch StashType {
						case AddCall:
							commonState.AddCall(NewOrderStash)
							commonState.NullAckmap()
							commonState.Ack()
							
						case RemoveCall:
							commonState.removeCall(RemoveOrderStash)
							commonState.NullAckmap()
							commonState.Ack()

						case StateChange:
							commonState.toHRAElevState(StateStash)
							commonState.NullAckmap()
							commonState.Ack()
						}	
						case commonState.Origin == config.Elevator_id  && stashed:
							stashed = false
							acking = false
						default:
							acking = false
					}
			
				case commonStatesEqual(commonState, arrivedCommonState): 
					arrivedCommonState.Ack()
					commonState = arrivedCommonState
					fmt.Println("ACKING IN SENDING SELF")

				default:
				}
			case peers := <-recieveFromPeerC: //bufferes lage stor kanal 64 feks
				commonState.makeElevUnav(peers)
				commonState.makeElevav()
			default:
			}
		}
		select {
		case <-heartbeatTimer.C:
			giverToNetwork <- commonState
		default:
		}
	}
}
