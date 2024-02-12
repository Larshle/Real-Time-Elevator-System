package requests

import (
	"root/driver/elevio"
	"root/elevator/localElevator"

)

// HandleButtonPress handles the button press event and updates the elevator state

func hallUp(e localElevator.Elevator) bool {
	for i:= 0; i < elevio._numFloors; i++ {
		if e.assingments[i][elevio.BT_HallUp] {
			return true
		}
	}
	return false
}

func hallDown(e localElevator.Elevator) bool {
	for i:= 0; i < elevio._numFloors; i++ {
		if e.assingments[i][elevio.BT_HallDown] {
			return true
		}
	}
	return false
}

func cabCall(e localElevator.Elevator) bool {
	for i:= 0; i < elevio._numFloors; i++ {
		if e.assingments[i][elevio.BT_Cab] {
			return true
		}
	}
	return false
}

func emptyHallUp(e localElevator.Elevator) bool {
	for i:= 0; i < elevio._numFloors; i++ {
		if !e.assingments[i][elevio.BT_HallUp] {
			return true
		}
	}
	return false
}

