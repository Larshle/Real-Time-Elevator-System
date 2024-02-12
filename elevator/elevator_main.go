package main

import(
	"root/driver/elevio"
)

func main() {
	elevio.Init("localhost:15657", 4)
	println("Elevator staretd")

}