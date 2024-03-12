package elevator

import (
	"root/config"
	"root/elevio"
)

// HandleButtonPress handles the button press event and updates the elevator state
// Make config file for parameters of elevator
// HUSK ikke ha med EmptyAssigner

type Assignments [config.NumFloors][config.NumButtons] bool

func (a Assignments) ReqInDirection(floor int, dir Direction) bool {
	switch dir {
		case Up:
			for f := floor + 1; f < config.NumFloors; f++ {
				for b := 0; b < config.NumButtons; b++ {
					if a[f][b] {
						return true
					}
				}
			}
			return false
		case Down:
			for f := 0; f < floor; f++ {
				for b := 0; b < config.NumButtons; b++ {
					if a[f][b] {
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