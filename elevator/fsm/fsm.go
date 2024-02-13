package fsm

import (
	"fmt"
	"root/driver/elevio"
	"root/elevator/localElevator"
	"root/elevator/requests"
	"time"
)

func HandleButtonPress(e *localElevator.Elevator, outputDevice *ElevOutputDevice, floor int, elevio.ButtonType Button) {
    // Handle button press event and update elevator state

	switch e.Behavior {
	case localElevator.EB_Idle:
		if e.CurrentFloor == floor {
			elevio.SetButtonLamp(floor, Button, true)
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetDoorOpenLamp(true)
			e.Behavior = localElevator.EB_DoorOpen
			time.Sleep(3 * time.Second)
			elevio.SetDoorOpenLamp(false)
			elevio.SetMotorDirection(elevio.MD_Up)
			e.Behavior = localElevator.EB_Moving
		} else {
			e.Assingments[floor][Button] = true
			elevio.SetButtonLamp(floor, Button, true)
			e.GetElevatorDirection()
			elevio.SetMotorDirection(e.MotorDirection)
			e.Behavior = localElevator.EB_Moving
		}
		
	case localElevator.EB_DoorOpen:
		if e.CurrentFloor == floor {
			elevio.SetButtonLamp(floor, Button, true)
			elevio.SetDoorOpenLamp(true)
			time.Sleep(3 * time.Second)
			elevio.SetDoorOpenLamp(false)
			requests.EmptyHall(e)
			elevio.SetButtonLamp(floor, Button, false)
		} else {
			e.Assingments[floor][Button] = true
			elevio.SetButtonLamp(floor, Button, true)
		}
	
	case localElevator.EB_Moving:
		e.Assingments[floor][Button] = true
		elevio.SetButtonLamp(floor, Button, true)

	}

}




func(e *localElevator.Elevator) GetElevatorDirection(){
	e.MotorDirection = int(dir)
	return e.MotorDirection
}

func main() {
	e =: localElevator.initializeElevator()



	e.Direction := elevio.GetMotorDirection()
	fmt.Println(motorDir)

	time.Sleep(100 * time.Millisecond)

}