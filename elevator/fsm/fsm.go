package fsm

import (
	"root/driver/elevio"
	"root/elevator/elevator"
	"fmt"
	"time"

)

func HandleButtonPress(elevator *localElevator, outputDevice *ElevOutputDevice, floor int, buttonType Button) {
    // Handle button press event and update elevator state
}


func(e *localElevator) ElevatorInit() {
	e.currentFloor = 0
	e.MotorDirection = 0
	e.doorsOpen = false
}

func(e *elevator) GetElevatorDirection(){
	e.MotorDirection = int(dir)
}

func main() {
	motorDir := elevio.GetMotorDirection()
	fmt.Println(motorDir)

	time.Sleep(100 * time.Millisecond)

}