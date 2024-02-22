package distributor

import (
	"root/driver/elevio"
)

type localAssignments struct {
	elevator_id string
	localCabAssignments [4]bool
	localHallAssignments [4][2]bool
}




func (a localAssignments) Add_Assingment(newAssignments elevio.ButtonEvent) localAssignments{
	if newAssignments.Button == elevio.BT_Cab {
		a.localCabAssignments[newAssignments.Floor] = true
	} else {
		a.localHallAssignments[newAssignments.Floor][newAssignments.Button] = true
	}
	return a
}

func (a localAssignments) Remove_Assingment( deliveredAssingement elevio.ButtonEvent) localAssignments{
	if deliveredAssingement.Button == elevio.BT_Cab {
		a.localCabAssignments[deliveredAssingement.Floor] = false
	} else {
		a.localHallAssignments[deliveredAssingement.Floor][deliveredAssingement.Button] = false
	}
	return a
}

func Update_Assingments(newAssingemntC <-chan elevio.ButtonEvent, deliveredAssingmentC <-chan elevio.ButtonEvent, updatedAssingmentsC chan<- localAssignments) {
	var localAssignments localAssignments
	for{
		select{
		case newAssingment := <- newAssingemntC:
			localAssignments = localAssignments.Add_Assingment(newAssingment)
			updatedAssingmentsC <- localAssignments
		case deliveredAssingment := <- deliveredAssingmentC:
			localAssignments = localAssignments.Remove_Assingment(deliveredAssingment)
		}
	}
}
