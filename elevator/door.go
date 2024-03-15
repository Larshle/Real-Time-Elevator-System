package elevator

import (
	"root/config"
	"root/elevio"
	"time"
)
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
const(
	DoorOpenDuration = 3*time.Second
)
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)

type DoorState int

const (
	Open DoorState = iota
	Closed
	Obstructed
	InCountDown
)

<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool, barkC chan<- bool) {

=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
func Door(doorClosedC chan<- bool, doorOpenC <-chan bool) {
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
	elevio.SetDoorOpenLamp(false)
	obstructionC := make(chan bool)
	go elevio.PollObstructionSwitch(obstructionC)

	obstruction := false
	timeCounter := time.NewTimer(time.Hour)
	var ds DoorState = Closed

	for {
		select {
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
		case obstruction = <-obstructionC:
			if !obstruction && ds == Obstructed {
				elevio.SetDoorOpenLamp(false)
				doorClosedC <- true
				ds = Closed
			}
			if obstruction {
				barkC <- true
			} else {
				barkC <- false
			}

		case <-doorOpenC:
			if obstruction {
				barkC <- true
			}
			switch ds {
			case Closed:
				elevio.SetDoorOpenLamp(true)
				timeCounter = time.NewTimer(config.DoorOpenDuration)
				ds = InCountDown
			case InCountDown:
				timeCounter = time.NewTimer(config.DoorOpenDuration)

			case Obstructed:
				timeCounter = time.NewTimer(config.DoorOpenDuration)
				ds = InCountDown
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}

>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)

=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
			case obstruction = <-obstructionC:
				if !obstruction && ds == Obstructed{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true 
				}


>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
			case <-doorOpenC:
				switch ds{
					case InCountDown:
						timeCounter = time.NewTimer(DoorOpenDuration)
					case Obstructed:
						timeCounter = time.NewTimer(DoorOpenDuration)
						ds = InCountDown
					case Closed:
						elevio.SetDoorOpenLamp(true)
						timeCounter = time.NewTimer(DoorOpenDuration)
						ds = InCountDown
					default:
						panic("Door state not implemented")
				}
			case <-timeCounter.C:
				if ds != InCountDown{
					panic("Door state not implemented")
				}
				if obstruction{
					ds = Obstructed
				}else{
					elevio.SetDoorOpenLamp(false)
					doorClosedC <- true
					ds = Closed
				}
		}
	}
}
