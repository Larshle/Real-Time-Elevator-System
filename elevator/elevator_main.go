package elevator_main

import(
	"root/driver/elevio"
	"root/elevator/localElevator"
)

func main() {
	elevio.Init("localhost:15657", 4)
	println("Elevator staretd")

	switch localElevator.Elevator.Assingments {
	case condition:
		
	}

}