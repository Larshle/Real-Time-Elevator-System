package elevator

import (

	"root/driver/elevio"
)

type State struct {
	Direction Direction
	Behavior  Behavior
	Floor int 
}

type Behavior int

const (
	Idle Behavior = iota
	Moving
	DoorOpen
)


func Elevator_Fsm( assingmentsC <-chan Assingments, stateC chan<- State, orderDelivered chan<- elevio.ButtonEvent){
		
		doorOpenC := make(chan bool, 16)
		doorClosedC := make(chan bool, 16)
		floorEnteredC := make(chan int)

		go Door(doorClosedC, doorOpenC)
		go elevio.PollFloorSensor(floorEnteredC)
		
		// Initialize elevator
		elevio.SetDoorOpenLamp(false)
		elevio.SetMotorDirection(elevio.MD_Down)
		state := State{Direction: Down, Behavior:  Moving}
		
		var assingments Assingments



		for {
			select {
				case <- doorClosedC:
					switch state.Behavior{
						case DoorOpen:
							switch{
								case assingments[state.Floor][state.Direction.toOpposite()]:
									elevio.SetMotorDirection(state.Direction.toOpposite().toMD())
									EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
									state.Direction = state.Direction.toOpposite()
									state.Behavior = Moving
									stateC <- state

								case assingments.ReqInDirection(state.Floor, state.Direction):
									elevio.SetMotorDirection(state.Direction.toMD())
									state.Behavior = Moving
									stateC <- state

								case assingments.ReqInDirection(state.Floor, state.Direction.toOpposite()):
									elevio.SetMotorDirection(state.Direction.toOpposite().toMD())
									state.Direction = state.Direction.toOpposite()
									state.Behavior = Moving
									EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
									stateC <- state

								default:
									state.Behavior = Idle
									stateC <- state
								}
							default:
							panic("DoorClosed in wrong state")
					}
				
				case f := <- floorEnteredC:
					state.Floor = f
					elevio.SetFloorIndicator(state.Floor)
					switch state.Behavior{
						case Moving:
							switch {

								
								case assingments[state.Floor][state.Direction]:
									elevio.SetMotorDirection(elevio.MD_Stop)
									state.Behavior = DoorOpen
									EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
									doorOpenC <- true

								case assingments[state.Floor][state.Direction] && assingments[state.Floor][elevio.BT_Cab]:
									elevio.SetMotorDirection(elevio.MD_Stop)
									state.Behavior = DoorOpen
									EmptyAssingner(state.Floor, state.Direction.toOpposite(), assingments, orderDelivered)
									doorOpenC <- true

								case assingments[state.Floor][elevio.BT_Cab] && !assingments[state.Floor][state.Direction.toOpposite()]:
									elevio.SetMotorDirection(elevio.MD_Stop)
									EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
									state.Behavior = DoorOpen
									doorOpenC <- true

								case assingments.ReqInDirection(state.Floor, state.Direction):

								case assingments[state.Floor][state.Direction.toOpposite()]:
									elevio.SetMotorDirection(elevio.MD_Stop)
									EmptyAssingner(state.Floor, state.Direction.toOpposite(), assingments, orderDelivered)
									state.Behavior = DoorOpen
									doorOpenC <- true

								
								case assingments.ReqInDirection(state.Floor, state.Direction.toOpposite()):
									state.Direction = state.Direction.toOpposite()
									elevio.SetMotorDirection(state.Direction.toMD())

								default:
									elevio.SetMotorDirection(elevio.MD_Stop)
									state.Behavior = Idle
									stateC <- state

							}
						default:
							panic("FloorEntered in wrong state")
					}

				case assingments = <- assingmentsC:
					switch state.Behavior{
						case Idle:
							switch{
								case assingments[state.Floor][state.Direction] || assingments[state.Floor][elevio.BT_Cab]:
									EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
									state.Behavior = DoorOpen
									doorOpenC <- true
									stateC <- state

								case assingments[state.Floor][state.Direction.toOpposite()]:
									EmptyAssingner(state.Floor, state.Direction.toOpposite(), assingments, orderDelivered)
									state.Direction = state.Direction.toOpposite()
									state.Behavior = DoorOpen
									doorOpenC <- true
									stateC <- state

								case assingments.ReqInDirection(state.Floor, state.Direction):
									elevio.SetMotorDirection(state.Direction.toMD())	
									state.Behavior = Moving
									stateC <- state
								case assingments.ReqInDirection(state.Floor, state.Direction.toOpposite()):
									elevio.SetMotorDirection(state.Direction.toOpposite().toMD())
									state.Direction = state.Direction.toOpposite()
									state.Behavior = Moving
									stateC <- state
								default:
							}
						
							case DoorOpen:
								switch{
									case assingments[state.Floor][state.Direction] || assingments[state.Floor][elevio.BT_Cab]:
										EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
										doorOpenC <- true
									
								}
									
							case Moving:
							
							default:
								panic("Assingments in wrong state")
						}
				}
			}
}


