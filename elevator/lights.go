package elevator

import (
	"root/elevio"
	"root/config"
	"root/distributor"
)


func SetLights(cs distributor.CommonState, ElevatorID int){
	for f := 0; f < config.NumFloors ; f++ {
		for b := 0; b < config.NumButtons; b++ {
			if cs.HallRequests[f][b]{ 
				elevio.SetButtonLamp( elevio.ButtonType(b), f, true)
			} else {
				elevio.SetButtonLamp( elevio.ButtonType(b), f, false)
			}

		}
	}
	for f := 0; f < config.NumFloors ; f++ {
		if cs.States[ElevatorID].CabRequests[f]{
		elevio.SetButtonLamp( elevio.ButtonType(2), f, true)
		} else {
			elevio.SetButtonLamp( elevio.ButtonType(2), f, false)
		}

	}
}
