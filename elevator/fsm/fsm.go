package elevator_fsm

import (
	"fmt"
	"root/driver/elevio"
	"root/elevator/direction"
	"root/elevator/requests"
	"time"
)

type State struct {
	Direction direction.Direction
	Behavior  Behavior
	Floor int 
}

type Behavior int

const (
	Idle Behavior = iota
	Moving
	DoorOpen
)


func Elevator_Fsm( assingmentsC <-chan Assingments, 
		stateC chan<- localElevator.ElevatorBehavior,
		orederDelivered chan<- bool,
		arrivedAtFloor chan<- int,
		obstruction chan<- bool){
		
		doorOpenC := make(chan bool, config.ChanSize)
		doorClosedC := make(chan bool, config.ChanSize)
		floorEnteredC := make(chan int)

		go door.Door(doorOpenC, doorClosedC)
		go elevio.PollFloorSensor(floorEnteredC)
		
		
		elevio.SetDoorOpenLamp(false)
		elevio.SetMotorDirection(elevio.MD_Down)
		state = State{Direction: direction.DownDown,Behavior:  Moving}
		
		var assingments Assingments



		for {
			select {
				case <- doorClosedC:
					switch state.Behavior{
						case DoorOpen:
							switch{
								case requests.HallCallUp(e):
									e.Direction = elevio.MD_Up
									e.Behavior = localElevator.EB_Moving
									elevio.SetMotorDirection(elevio.MD_Up)
								case requests.HallCallDown(e):
									e.Direction = elevio.MD_Down
									e.Behavior = localElevator.EB_Moving
									elevio.SetMotorDirection(elevio.MD_Down)

								case requests.CabCall(e):
							


								default:
									state.Behavior = Idle
									stateC <- state
								}
							default:
							panic("DoorClosed in wrong state")
						}
				
				case state.Floor = <- floorEnteredC:
					elevio.SetFloorIndicator(state.Floor)
					switch state.Behavior{
						case Moving:
							switch {

							}
						default:
							panic("FloorEntered in wrong state")
							stateC <- state
					}
				
				case assingments = <- assingmentsC:
					


	
				
				
			select {
			case o := <- orederDelivered:{

			}
			}
			}
		
		

		}


	}
}




