package distributor

import (
	"fmt"
	"root/driver/elevio"
)

type Ass int

const (
	add    Ass = 1
	remove Ass = 2
)

type localAssignments struct {
	localCabAssignments  [4]Ass
	localHallAssignments [4][2]Ass
}

func (a *localAssignments) Add_Assingment(newAssignments elevio.ButtonEvent) {
	if newAssignments.Button == elevio.BT_Cab {
		a.localCabAssignments[newAssignments.Floor] = add
		fmt.Println("HER SJEKKER VI 1")
	}
	if newAssignments.Button == elevio.BT_HallDown {
		a.localHallAssignments[newAssignments.Floor][elevio.BT_HallDown] = add
		fmt.Println("HER SJEKKER VI 2")
	}
	if newAssignments.Button == elevio.BT_HallUp {
		a.localHallAssignments[newAssignments.Floor][elevio.BT_HallUp] = add
		fmt.Println("HER SJEKKER VI 3")
	}
}

func (a localAssignments) Remove_Assingment(deliveredAssingement elevio.ButtonEvent) localAssignments {
	fmt.Println("Knappen: ", deliveredAssingement.Button)
	if deliveredAssingement.Button == elevio.BT_Cab {
		a.localCabAssignments[deliveredAssingement.Floor] = remove
		fmt.Println("HER SJEKKER VI 4")
	}
	if deliveredAssingement.Button == elevio.BT_HallDown {
		a.localHallAssignments[deliveredAssingement.Floor][elevio.BT_HallDown] = remove
		fmt.Println("HER SJEKKER VI 5")
	}
	if deliveredAssingement.Button == elevio.BT_HallUp {
		a.localHallAssignments[deliveredAssingement.Floor][elevio.BT_HallUp] = remove
		fmt.Println("HER SJEKKER VI 6")
	}
	return a
}

func Update_Assingments(newAssingemntC <-chan elevio.ButtonEvent, deliveredAssingmentC <-chan elevio.ButtonEvent, updatedAssingmentsC chan<- localAssignments) {
	var localAssignments localAssignments
	fmt.Println("Got message from elev")
	for {
		select {
		case newAssingment := <-newAssingemntC:
			fmt.Println("jdfhgbd")
			localAssignments.Add_Assingment(newAssingment)
			updatedAssingmentsC <- localAssignments
		case deliveredAssingment := <-deliveredAssingmentC:
			fmt.Println("Got delivered assingment")
			localAssignments = localAssignments.Remove_Assingment(deliveredAssingment)
			updatedAssingmentsC <- localAssignments
		}
	}
}
