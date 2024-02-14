package localElevator

import (
	"root/driver/elevio"
)




type ElevatorBehavior int

const (
	EB_Idle ElevatorBehavior = iota
	EB_DoorOpen
	EB_Moving
)

type Elevator struct {
	CurrentFloor int
	Direction elevio.MotorDirection
	Behavior ElevatorBehavior
	Assingments [][]bool
	Obstruction bool

}


func initializeElevator() Elevator{
	e := Elevator{}
	e.CurrentFloor = 0
	e.Direction = 0
	e.Behavior = EB_Idle
	e.Obstruction = false
	e.Assingments = make([][]bool, 4)
	for i := range e.Assingments {
		e.Assingments[i] = make([]bool, 3)
	}
	return e
}

