package elevator

import (
	"root/elevio"
)

// HandleButtonPress handles the button press event and updates the elevator state
// Make config file for parameters of elevator
// HUSK ikke ha med EmptyAssigner

type Assignments [4][3] bool

func (a Assignments) ReqInDirection(floor int, dir Direction) bool {
	switch dir {
		case Up:
			for i := floor + 1; i < elevio.NumFloors; i++ {
				for j := 0; j < 3; j++ {
					if a[i][j] {
						return true
					}
				}
			}
			return false
		case Down:
			for i := 0; i < floor; i++ {
				for j := 0; j < 3; j++ {
					if a[i][j] {
						return true
					}
				}
			}
			return false
		default:
			panic("Invalid direction")
		}
}
	


func EmptyAssigner(floor int, dir Direction, a Assignments, orderDoneC chan<- elevio.ButtonEvent) {
	if a[floor][elevio.BT_Cab] {
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
	}
	if a[floor][dir] {
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: dir.toBT()}
	}
}