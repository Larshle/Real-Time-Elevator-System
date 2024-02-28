package distributor

import (
	"root/driver/elevio"
)

type Ass int

const(
	add Ass = 1
	remove
)


type localAssignments struct {
	localCabAssignments [4]Ass
	localHallAssignments [4][2]Ass
}

func (a localAssignments) Add_Assingment(newAssignments elevio.ButtonEvent) localAssignments{
	if newAssignments.Button == elevio.BT_Cab {
		a.localCabAssignments[newAssignments.Floor] = add
	} else {
		a.localHallAssignments[newAssignments.Floor][newAssignments.Button] = add
	}
	return a
}

func (a localAssignments) Remove_Assingment( deliveredAssingement elevio.ButtonEvent) localAssignments{
	if deliveredAssingement.Button == elevio.BT_Cab {
		a.localCabAssignments[deliveredAssingement.Floor] = remove
	} else {
		a.localHallAssignments[deliveredAssingement.Floor][deliveredAssingement.Button] = remove
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
