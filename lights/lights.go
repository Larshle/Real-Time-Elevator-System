package lights

import (
	"root/elevio"
	"root/elevator"
	"root/config"
)

func SetLights(lightAssignment elevator.Assignments) {
	for floor := 0; floor < config.NumFloors; floor++ {
		for button := 0; button < 3; button++ {
			if lightAssignment[floor][button] {
				elevio.SetButtonLamp( elevio.ButtonType(button), floor, true)
			} else {
				elevio.SetButtonLamp( elevio.ButtonType(button), floor, false)
			}
		}
	}
}

func Lights(lightsAssignmentC <-chan elevator.Assignments) {
	for {
		select {
		case a := <-lightsAssignmentC:
			SetLights(a)
		}
	}
}