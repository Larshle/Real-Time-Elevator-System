package distributor

import (
	"root/driver/elevio"
	
)

type Ass int

const (
	untouched Ass = iota
	add   
	remove 
)

type localAssignments struct {
	localCabAssignments  [4]Ass
	localHallAssignments [4][2]Ass
}

func (a *localAssignments) Add_Assingment(newAssignments elevio.ButtonEvent) {
	if newAssignments.Button == elevio.BT_Cab {
		a.localCabAssignments[newAssignments.Floor] = add
	}
	if newAssignments.Button == elevio.BT_HallDown {
		a.localHallAssignments[newAssignments.Floor][elevio.BT_HallDown] = add
	}
	if newAssignments.Button == elevio.BT_HallUp {
		a.localHallAssignments[newAssignments.Floor][elevio.BT_HallUp] = add
	}
}

func (a localAssignments) Remove_Assingment(deliveredAssingement elevio.ButtonEvent) localAssignments {
	if deliveredAssingement.Button == elevio.BT_Cab {
		a.localCabAssignments[deliveredAssingement.Floor] = remove
	}
	if deliveredAssingement.Button == elevio.BT_HallDown {
		a.localHallAssignments[deliveredAssingement.Floor][elevio.BT_HallDown] = remove
	}
	if deliveredAssingement.Button == elevio.BT_HallUp {
		a.localHallAssignments[deliveredAssingement.Floor][elevio.BT_HallUp] = remove
	}
	return a
}


func Update_Assingments(newAssingemntC <-chan elevio.ButtonEvent, deliveredAssingmentC <-chan elevio.ButtonEvent, updatedAssingmentsC chan<- localAssignments) {
	var localAssignments localAssignments
	for {
		select {
		case newAssingment := <-newAssingemntC:
			localAssignments.Add_Assingment(newAssingment)
			updatedAssingmentsC <- localAssignments
		case deliveredAssingment := <-deliveredAssingmentC:
			localAssignments = localAssignments.Remove_Assingment(deliveredAssingment)
			updatedAssingmentsC <- localAssignments
		}
	}
}
