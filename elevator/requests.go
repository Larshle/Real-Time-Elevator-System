package elevator

import (
	"root/driver/elevio"
)

// HandleButtonPress handles the button press event and updates the elevator state
// Make config file for parameters of elevator

type Assingments [4][3] bool

func (a Assingments) ReqInDirection(floor int, dir Direction) bool {
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
	


func EmptyAssingner(floor int, dir Direction, a Assingments, orderDoneC chan<- elevio.ButtonEvent) {
	if a[floor][elevio.BT_Cab] {
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
	}
	if a[floor][dir] {
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: dir.toBT()}
	}
}