package fsm

import (
	"fmt"
	"root/driver/elevio"
	"root/elevator/localElevator"
	"root/elevator/requests"
	"time"
)

func Fsm( c
	//Add channels here
	){
		localElev = localElevator.initializeElevator()
		e = &localElev

		for{
			select{
			e.CurrentFloor = <- ch_arrivedAtFloors
			case e.CurrentFloor = -1:
				elevio.SetMotorDirection(elevio.MD_Down)
			}
			case e.CurrentFloor != -1:
				e.Direction = elevio.MD_Stop
				elevio.SetMotorDirection(elevio.MD_Down)

		}

	}




	go elevio.PollFloorSensor(ch_arrivedAtFloors)
	go elevio.PollObstructionSwitch(ch_obstruction)
	go elevio.PollButtons(ch_newLocalOrder)


	ch_orderChan chan elevio.ButtonEvent,
	ch_elevatorState chan<- elevator.Elevator,
	ch_clearLocalHallOrders chan bool,
	ch_arrivedAtFloors chan int,
	ch_obstruction chan bool,
	ch_timerDoor chan bool) {

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