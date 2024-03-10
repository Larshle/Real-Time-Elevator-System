package elevator

import (
	"root/driver/elevio"
)

type State struct {
	Direction Direction
	Behaviour Behaviour
	Floor     int
}

type Behaviour int

const (
	Idle Behaviour = iota
	Moving
	DoorOpen
)

func (b Behaviour) ToString() string {
	return map[Behaviour]string{Idle: "idle", Moving: "moving", DoorOpen: "doorOpen"}[b]
}

func Elevator(eleveatorAssingmentC <-chan Assingments, stateC chan<- State, orderDelivered chan<- elevio.ButtonEvent) {
	doorOpenC := make(chan bool, 16)
	doorClosedC := make(chan bool, 16)
	floorEnteredC := make(chan int)

	go Door(doorClosedC, doorOpenC)
	go elevio.PollFloorSensor(floorEnteredC)

	// Initialize elevator
	elevio.SetDoorOpenLamp(false)
	elevio.SetMotorDirection(elevio.MD_Down)
	state := State{Direction: Down, Behaviour: Moving}

	var assingments Assingments

	for {
		select {
		case <-doorClosedC:
			switch state.Behaviour {
			case DoorOpen:
				switch {

				case assingments.ReqInDirection(state.Floor, state.Direction):
					elevio.SetMotorDirection(state.Direction.toMD())
					state.Behaviour = Moving
					stateC <- state

				case assingments[state.Floor][state.Direction.toOpposite()]:
					doorOpenC <- true
					state.Direction = state.Direction.toOpposite()
					EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
					stateC <- state

				case assingments.ReqInDirection(state.Floor, state.Direction.toOpposite()):
					state.Direction = state.Direction.toOpposite()
					elevio.SetMotorDirection(state.Direction.toOpposite().toMD())
					state.Behaviour = Moving
					stateC <- state

				default:
					state.Behaviour = Idle
					stateC <- state
				}
			default:
				panic("DoorClosed in wrong state")
			}

		case state.Floor = <-floorEnteredC:
			elevio.SetFloorIndicator(state.Floor)
			switch state.Behaviour {
			case Moving:
				switch {
				case assingments[state.Floor][state.Direction]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
					state.Behaviour = DoorOpen

				case assingments[state.Floor][elevio.BT_Cab] && assingments.ReqInDirection(state.Floor, state.Direction):
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
					state.Behaviour = DoorOpen
					


				case assingments[state.Floor][elevio.BT_Cab] && !assingments[state.Floor][state.Direction.toOpposite()]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
					state.Behaviour = DoorOpen


				case assingments.ReqInDirection(state.Floor, state.Direction):

				case assingments[state.Floor][state.Direction.toOpposite()]:
					elevio.SetMotorDirection(elevio.MD_Stop)
					doorOpenC <- true
					state.Direction = state.Direction.toOpposite()
					EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
					state.Behaviour = DoorOpen
					

				case assingments.ReqInDirection(state.Floor, state.Direction.toOpposite()):
					state.Direction = state.Direction.toOpposite()
					elevio.SetMotorDirection(state.Direction.toMD())

				default:
					elevio.SetMotorDirection(elevio.MD_Stop)
					state.Behaviour = Idle
				}
			default:
				panic("FloorEntered in wrong state")
			}
			stateC <- state

		case assingments = <-eleveatorAssingmentC:
			switch state.Behaviour {
			case Idle:
				switch {
				case assingments[state.Floor][state.Direction] || assingments[state.Floor][elevio.BT_Cab]:
					doorOpenC <- true
					EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
					state.Behaviour = DoorOpen
					stateC <- state

				case assingments[state.Floor][state.Direction.toOpposite()]:
					doorOpenC <- true
					state.Direction = state.Direction.toOpposite()
					EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)
					state.Behaviour = DoorOpen
					stateC <- state

				case assingments.ReqInDirection(state.Floor, state.Direction):
					elevio.SetMotorDirection(state.Direction.toMD())
					state.Behaviour = Moving
					stateC <- state

				case assingments.ReqInDirection(state.Floor, state.Direction.toOpposite()):
					state.Direction = state.Direction.toOpposite()
					elevio.SetMotorDirection(state.Direction.toMD())
					state.Behaviour = Moving
					stateC <- state
				default:
				}

			case DoorOpen:
				switch {
				case assingments[state.Floor][state.Direction] || assingments[state.Floor][elevio.BT_Cab]:
					doorOpenC <- true
					EmptyAssingner(state.Floor, state.Direction, assingments, orderDelivered)

				}

			case Moving:

			default:
				panic("Assingments in wrong state")
			}
		}
	}
}
