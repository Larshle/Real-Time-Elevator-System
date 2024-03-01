package lights

import (
	"fmt"
	"root/driver/elevio"
	"root/elevator"
	"root/config"
)

func SetLights(light_assignment elevator.Assingments) {
	for floor := 0; floor < config.N_floors; floor++ {
		for button := 0; button < 3; button++ {
			if light_assignment[floor][button] {
				elevio.SetButtonLamp( elevio.ButtonType(button), floor, true)
			} else {
				elevio.SetButtonLamp( elevio.ButtonType(button), floor, false)
			}
		}
	}
}

func Lights(lightsAssingmentChan <-chan elevator.Assingments) {
	fmt.Print("Lights started\n")
	for {
		select {
		case a := <-lightsAssingmentChan:
			SetLights(a)
		}
	}
}