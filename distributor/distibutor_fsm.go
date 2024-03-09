package distributor

import (
	// "root/config"
	"root/driver/elevio"
	"root/elevator"
	"root/network/network_modules/peers"
	"time"
	"fmt"
)
type StatshType int 
const (
	RemoveCall StatshType = iota
	AddCall
	StateChange
)


type State int 
const (
	Idle State = iota
	Acking
	SendingSelf
	AckingOtherWhileTryingToSendSelf
	Isolated
	UnableToMove
)

func Distributor(
	deliveredOrderC <-chan elevio.ButtonEvent,
	newElevStateC <-chan elevator.State,
	giverToNetwork chan<- HRAInput,
	receiveFromNetworkC <-chan HRAInput,
	messageToAssinger chan<- HRAInput) {

	elevioOrdersC := make(chan elevio.ButtonEvent, 64)
	newAssingemntC := make(chan localAssignments, 64)
	peerUpdateC := make(chan peers.PeerUpdate, 64)

	var cs HRAInput
	var StateStash elevator.State
	//var AssignmentStash localAssignments
	var state State = Idle
	var StashType StatshType

	var NewOrderStash elevio.ButtonEvent
	var RemoveOrderStash elevio.ButtonEvent

	cs = HRAInput{
		Origin: "peer-10.22.229.227-22222",
		Seq:    0,
		Ackmap: map[string]Ack_status{},
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
		},
	}

	go elevio.PollButtons(elevioOrdersC)
	//go Update_Assingments(elevioOrdersC, deliveredOrderC, newAssingemntC)

	timeCounter := time.NewTimer(time.Hour)
	heartbeatTimer := time.NewTicker(15 * time.Millisecond)
	selfLostNetworkDuratio := 1 * time.Second

	for {

		select{
			case <- timeCounter.C:
				state = Isolated
			default:
		}

		switch state {
			case Idle:
				select {

					case newOrder := <-elevioOrdersC:
						StashType = AddCall
						NewOrderStash = newOrder
						cs.AddCall(newOrder)
						cs.Ack()
						state = SendingSelf
		
					case removeOrder := <-deliveredOrderC:
						StashType = RemoveCall
						RemoveOrderStash = removeOrder
						cs.removeCall(removeOrder)
						cs.Ack()
						state = SendingSelf

					// case assingmentUpdate := <-newAssingemntC:
					// 	StashType = AssingmetChange
					// 	AssignmentStash = assingmentUpdate
					// 	cs.Update_Assingments(assingmentUpdate)
					// 	cs.Ack()
					// 	state = SendingSelf
					// 	fmt.Println("Idle: assingmentUpdate")

					case newElevState := <-newElevStateC:
						StashType = StateChange
						StateStash = newElevState
						cs.toHRAElevState(newElevState)
						cs.Ack()
						state = SendingSelf
						fmt.Println("Idle: newElevState")

					case arrivedCommonState := <-receiveFromNetworkC:
						timeCounter = time.NewTimer(selfLostNetworkDuratio) 

						switch {
							case higherPriority(cs, arrivedCommonState):
								cs = arrivedCommonState
								arrivedCommonState.Ack()
								cs = arrivedCommonState
								state = Acking
							default:
								break
						}
					case peers := <- peerUpdateC:
						cs.makeElevUnav(peers)
					default:
				}
			case SendingSelf:
				select {
				case arrivedCommonState := <-receiveFromNetworkC:
					timeCounter = time.NewTimer(selfLostNetworkDuratio) 
					switch {
						case higherPriority(cs, arrivedCommonState):
							cs = arrivedCommonState
							cs.Ack()
							state = AckingOtherWhileTryingToSendSelf
							fmt.Println("SendingSelf: higherPriority")

						case Fully_acked(arrivedCommonState.Ackmap):
							cs = arrivedCommonState
							messageToAssinger <- cs
							cs.NullAckmap()
							state = Idle
							fmt.Println("SendingSelf: Fully_acked")
					}
				case peers := <- peerUpdateC:
					cs.makeElevUnav(peers)
					if Fully_acked(cs.Ackmap){
						state = Idle
						messageToAssinger <- cs
						cs.NullAckmap()
					}
				default:
					// sender cs på heartbeat
				}
					
				
			case Acking:
				select {
				case arrivedCommonState := <-receiveFromNetworkC:
					timeCounter = time.NewTimer(selfLostNetworkDuratio)
					switch {
					
					case Fully_acked(arrivedCommonState.Ackmap):
						state = Idle
						cs = arrivedCommonState
						messageToAssinger <- cs
						cs.NullAckmap()

					case higherPriority(cs, arrivedCommonState):
						arrivedCommonState.Ack()
						cs = arrivedCommonState
				
					case !higherPriority(cs, arrivedCommonState):
						break
					}

				case peers := <- peerUpdateC:
					cs.makeElevUnav(peers)
					if Fully_acked(cs.Ackmap){
						state = Idle
						messageToAssinger <- cs
						cs.NullAckmap()
					}
					
				default:
			}
			case AckingOtherWhileTryingToSendSelf:
				select {
				case arrivedCommonState := <-receiveFromNetworkC:
					fmt.Println("AckingOtherWhileTryingToSendSelf: receiveFromNetworkC")
					timeCounter = time.NewTimer(selfLostNetworkDuratio)
					switch {
					case !higherPriority(cs, arrivedCommonState):
						fmt.Println("AckingOtherWhileTryingToSendSelf: !higherPriority")
						break

					case higherPriority(cs, arrivedCommonState):
						fmt.Println("AckingOtherWhileTryingToSendSelf: higherPriority")
						cs = arrivedCommonState
						cs.Ack()
				
					case Fully_acked(arrivedCommonState.Ackmap):
						fmt.Println("AckingOtherWhileTryingToSendSelf: Fully_acked")
						state = SendingSelf
						switch StashType {

							case AddCall:
								arrivedCommonState.AddCall(NewOrderStash)
		
							case RemoveCall:
								arrivedCommonState.removeCall(RemoveOrderStash)
								

							case StateChange:
								arrivedCommonState.toHRAElevState(StateStash)
								fmt.Println("AckingOtherWhileTryingToSendSelf: Fully_acked: StateChange")
								
						}
						cs = arrivedCommonState
						messageToAssinger <- cs
						cs.NullAckmap()
					}
				case peers := <- peerUpdateC:
					cs.makeElevUnav(peers)
					if Fully_acked(cs.Ackmap){
						state = SendingSelf
						messageToAssinger <- cs
						cs.NullAckmap()
					}
				default:
			}
			case Isolated:
				select{

				case <-receiveFromNetworkC:
					state = Idle
				
				case assingmentUpdate := <-newAssingemntC:
					cs.makeElevUnavExceptOrigin()
					cs.UpdateCabAssignments(assingmentUpdate)
					messageToAssinger <- cs


				case newElevState := <-newElevStateC:
					cs.toHRAElevState(newElevState)
					cs.makeElevUnavExceptOrigin()
					messageToAssinger <- cs

		
				default:
				}
				
			}
		
			select {
			case <-heartbeatTimer.C:
				giverToNetwork <- cs
			default:
				}


		}	
	}

	
