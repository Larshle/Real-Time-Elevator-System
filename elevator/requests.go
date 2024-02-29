package elevator

import (
	"fmt"
	"root/driver/elevio"
)

// HandleButtonPress handles the button press event and updates the elevator state
// Make config file for parameters of elevator

type Assingments [4][3] bool

func (a Assingments) ReqInDirection(floor int, dir Direction) bool {
	switch dir {
		case Up:
			for i := floor + 1; i < elevio.NumFloors; i++ {
				if a[i][elevio.BT_HallUp] || a[i][elevio.BT_Cab] {
					return true
				}
			}
			return false
		case Down:
			for i := 0; i < floor; i++ {
				if a[i][elevio.BT_HallDown] || a[i][elevio.BT_Cab] {
					return true
				}
			}
			return false
		default:
			panic("Invalid direction")
		}
}
	


func EmptyAssingner(floor int, dir Direction, a Assingments, orderDoneC chan<- elevio.ButtonEvent) {
	fmt.Println("Emptying assingner")
	fmt.Println(a)
	fmt.Println("se oppe")
	if a[floor][elevio.BT_Cab] {
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
		fmt.Println("Cab order done")
	}
	if a[floor][dir] {
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: dir.toBT()}
		fmt.Println("Hall order done")
	}
}





