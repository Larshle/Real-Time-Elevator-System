package requests

import (
	"Driver-go/elevio"
	"time"

	"github.                                        com/vebjornwn/Sanntid-Prosjekt/Driver-go/elevio/elevio"
	"github.com/vebjornwn/Sanntid-Prosjekt/elevator/elevator/elevator"
)

// HandleButtonPress handles the button press event and updates the elevator state

func hallUp(e elevator.localElevator) bool {
	for i:= 0; i < elevio.numFloors; i++ {
		if e.assingments[i][elevio.BT_HallUp] {
			return true
		}
	}
	return false
}

func hallDown(e elevator.localElevator) bool {
	for i:= 0; i < elevio.numFloors; i++ {
		if e.assingments[i][elevio.BT_HallDown] {
			return true
		}
	}
	return false
}

func cabCall(e elevator.localElevator) bool {
	for i:= 0; i < elevio.numFloors; i++ {
		if e.assingments[i][elevio.BT_Cab] {
			return true
		}
	}
	return false
}

func emptyHallUp(e elevator.localElevator) bool {
	for i:= 0; i < elevio.numFloors; i++ {
		if !e.assingments[i][elevio.BT_HallUp] {
			return true
		}
	}
	return false
}

