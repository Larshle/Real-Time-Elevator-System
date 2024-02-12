package main

import(
	"fmt"
	"time"
	"github.com/vebjornwn/Sanntid-Prosjekt/Driver-go/elevio/elevio"
)

func main() {
	elevio.Init("localhost:15657", 4)
	println("Elevator staretd")

}