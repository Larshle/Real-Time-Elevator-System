package fsm

import (
	"Driver-go/elevio/elevatorio"
	"fmt"
	"time"

)

func HandleButtonPress(elevator *Elevator, outputDevice *ElevOutputDevice, floor int, buttonType Button) {
    // Handle button press event and update elevator state
}


func(e *elevator) ElevatorInit() {
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