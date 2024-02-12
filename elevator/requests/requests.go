package requests

import (
	"root/driver/elevio"
	"root/elevator/localElevator"

)

// HandleButtonPress handles the button press event and updates the elevator state

func HallCallUp(e localElevator.Elevator) bool {
	for i:= 0; i < elevio.NumFloors; i++ {
		if e.Assingments[i][elevio.BT_HallUp] {
			return true
		}
	}
	return false
}

func HallCallDown(e localElevator.Elevator) bool {
	for i:= 0; i < elevio.NumFloors; i++ {
		if e.Assingments[i][elevio.BT_HallDown] {
			return true
		}
	}
	return false
}

func CabCall(e localElevator.Elevator) bool {
	for i:= 0; i < elevio.NumFloors; i++ {
		if e.Assingments[i][elevio.BT_Cab] {
			return true
		}
	}
	return false
}

func EmptyHallUp(e localElevator.Elevator) bool {
    for i := 0; i < elevio.NumFloors; i++ {
        if e.CurrentFloor == i && e.Assingments[i][elevio.BT_HallUp] {
            e.Assingments[i][elevio.BT_HallUp] = false
            return true
        }
    }
    return false
}

func EmptyHallDown(e localElevator.Elevator) bool {
    for i := 0; i < elevio.NumFloors; i++ {
        if e.CurrentFloor == i && e.Assingments[i][elevio.BT_HallUp] {
            e.Assingments[i][elevio.BT_HallUp] = false
            return true
        }
    }
    return false
}



