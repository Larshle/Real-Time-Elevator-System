package elevator

import (
	"Driver-go/elevio/elevatorio"
	"fmt"
	"time"

)

type elevator struct {
	currentFloor int
	MotorDirection int
	doorsOpen bool
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