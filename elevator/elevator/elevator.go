package elevator

import (
	"root/driver/elevio/"
)




type elevatorBehavior int

const (
	EB_Idle elevatorBehavior = iota
	EB_DoorOpen
	EB_Moving
)

type localElevator struct {
	currentFloor int
	Direction elevio.MotorDirection
	Behavior elevatorBehavior
	assingments [][]bool

}


func initializeElevator() localElevator{
	e := localElevator{}
	e.currentFloor = 0
	e.Direction = 0
	e.Behavior = EB_Idle
	e.assingments = make([][]bool, 4)
	for i := range e.assingments {
		e.assingments[i] = make([]bool, 3)
	}
	return e
}

