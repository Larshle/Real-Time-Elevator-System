package elevator

import (
	"github.com/vebjornwn/Sanntid-Prosjekt/Driver-go/elevio/elevio"
)




type elevatorBehavior int

const (
	EB_Idle elevatorBehavior = iota
	EB_DoorOpen
	EB_Moving
)

type elevator struct {
	currentFloor int
	MotorDirection int
	Behavior elevatorBehavior
	assingments [][]bool

}





func(eb elevatorBehavior) String() string {
	switch eb {
	case EB_Idle:
		return "Idle"
	case EB_DoorOpen:
		return "DoorOpen"
	case EB_Moving:
		return "Moving"
	default:
		return "Unknown"
	}	
}