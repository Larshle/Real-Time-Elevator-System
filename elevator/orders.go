package elevator

import (
	"root/config"
	"root/elevio"
)

<<<<<<< HEAD:elevator/orders.go
type Orders [config.NumFloors][config.NumButtons]bool
=======
// HandleButtonPress handles the button press event and updates the elevator state
>>>>>>> a346d18674c6e853af0a1af46bf3dbee428466b7:elevator/requests.go

func (a Orders) OrderInDirection(floor int, dir Direction) bool {
	switch dir {
<<<<<<< HEAD
	case Up:
		for f := floor + 1; f < config.NumFloors; f++ {
			for b := 0; b < config.NumButtons; b++ {
				if a[f][b] {
					return true
				}
			}
=======
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

<<<<<<< HEAD:elevator/orders.go
func OrderDone(floor int, dir Direction, a Orders, orderDoneC chan<- elevio.ButtonEvent) {
=======
<<<<<<< HEAD
func EmptyAssigner(floor int, dir Direction, a Assignments, orderDoneC chan<- elevio.ButtonEvent) {
=======

func EmptyAssingner(floor int, dir Direction, a Assingments, orderDoneC chan<- elevio.ButtonEvent) bool {
>>>>>>> a346d18674c6e853af0a1af46bf3dbee428466b7:elevator/requests.go
	if a[floor][elevio.BT_Cab] {
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: elevio.BT_Cab}
	}
	if a[floor][dir] {
		orderDoneC <- elevio.ButtonEvent{Floor: floor, Button: dir.toBT()}
	}
	return false
}





