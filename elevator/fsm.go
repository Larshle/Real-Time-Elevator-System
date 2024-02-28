package elevator

import (

	"root/driver/elevio"
)

type State struct {
	Direction Direction
	Behaviour  Behaviour
	Floor int 
}

type Behaviour int

const (
	Idle Behaviour = iota
	Moving
	DoorOpen
)

func (b Behaviour) ToString() string {
	return map[Behaviour]string{Idle: "Idle", Moving: "Moving", DoorOpen: "DoorOpen"}[b]
}



func Elevator(eleveatorAssingmentC <-chan Assingments, stateC chan<- State, orderDelivered chan<- elevio.ButtonEvent){
		
		doorOpenC := make(chan bool, 16)
		doorClosedC := make(chan bool, 16)
		floorEnteredC := make(chan int)

		go Door(doorClosedC, doorOpenC)
		go elevio.PollFloorSensor(floorEnteredC)
		
		// Initialize elevator
		elevio.SetDoorOpenLamp(false)
		elevio.SetMotorDirection(elevio.MD_Down)
		state := State{Direction: Down, Behaviour:  Moving}
		
		var assingments Assingments



		for {
			select {
				case <- doorClosedC:
					switch state.Behaviour{
						case DoorOpen:
							switch{
								case assingments[state.Floor][state.Direction.toOpposite()]:
									elevio.SetMotorDirection(state.Direction.toOpposite().toMD())
									EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
									state.Direction = state.Direction.toOpposite()
									state.Behaviour = Moving
									stateC <- state

								case assingments.ReqInDirection(state.Floor, state.Direction):
									elevio.SetMotorDirection(state.Direction.toMD())
									state.Behaviour = Moving
									stateC <- state

								case assingments.ReqInDirection(state.Floor, state.Direction.toOpposite()):
									elevio.SetMotorDirection(state.Direction.toOpposite().toMD())
									state.Direction = state.Direction.toOpposite()
									state.Behaviour = Moving
									EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
									stateC <- state

								default:
									state.Behaviour = Idle
									stateC <- state
								}
							default:
							panic("DoorClosed in wrong state")
					}
				
				case f := <- floorEnteredC:
					state.Floor = f
					elevio.SetFloorIndicator(state.Floor)
					switch state.Behaviour{
						case Moving:
							switch {

								
								case assingments[state.Floor][state.Direction]:
									elevio.SetMotorDirection(elevio.MD_Stop)
									state.Behaviour = DoorOpen
									EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
									doorOpenC <- true

								case assingments[state.Floor][state.Direction] && assingments[state.Floor][elevio.BT_Cab]:
									elevio.SetMotorDirection(elevio.MD_Stop)
									state.Behaviour = DoorOpen
									EmptyAssingner(state.Floor, state.Direction.toOpposite(), assingments, orderDelivered)
									doorOpenC <- true

								case assingments[state.Floor][elevio.BT_Cab] && !assingments[state.Floor][state.Direction.toOpposite()]:
									elevio.SetMotorDirection(elevio.MD_Stop)
									EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
									state.Behaviour = DoorOpen
									doorOpenC <- true

								case assingments.ReqInDirection(state.Floor, state.Direction):

								case assingments[state.Floor][state.Direction.toOpposite()]:
									elevio.SetMotorDirection(elevio.MD_Stop)
									EmptyAssingner(state.Floor, state.Direction.toOpposite(), assingments, orderDelivered)
									state.Behaviour = DoorOpen
									doorOpenC <- true

								
								case assingments.ReqInDirection(state.Floor, state.Direction.toOpposite()):
									state.Direction = state.Direction.toOpposite()
									elevio.SetMotorDirection(state.Direction.toMD())

								default:
									elevio.SetMotorDirection(elevio.MD_Stop)
									state.Behaviour = Idle
									stateC <- state

							}
						default:
							panic("FloorEntered in wrong state")
					}

				case assingments = <- eleveatorAssingmentC:
					switch state.Behaviour{
						case Idle:
							switch{
								case assingments[state.Floor][state.Direction] || assingments[state.Floor][elevio.BT_Cab]:
									EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
									state.Behaviour = DoorOpen
									doorOpenC <- true
									stateC <- state

								case assingments[state.Floor][state.Direction.toOpposite()]:
									EmptyAssingner(state.Floor, state.Direction.toOpposite(), assingments, orderDelivered)
									state.Direction = state.Direction.toOpposite()
									state.Behaviour = DoorOpen
									doorOpenC <- true
									stateC <- state

								case assingments.ReqInDirection(state.Floor, state.Direction):
									elevio.SetMotorDirection(state.Direction.toMD())	
									state.Behaviour = Moving
									stateC <- state
								case assingments.ReqInDirection(state.Floor, state.Direction.toOpposite()):
									elevio.SetMotorDirection(state.Direction.toOpposite().toMD())
									state.Direction = state.Direction.toOpposite()
									state.Behaviour = Moving
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


