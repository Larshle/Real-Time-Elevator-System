package distributor

import (
	"fmt"
	"root/assigner"
)

var commonstate = assigner.HRAInput{
	HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
	States: map[string]assigner.HRAElevState{
		"one": assigner.HRAElevState{
			Behavior:    "moving",
			Floor:       2,
			Direction:   "up",
			CabRequests: []bool{false, false, false, true},
		},
		"two": assigner.HRAElevState{
			Behavior:    "idle",
			Floor:       3,
			Direction:   "stop",
			CabRequests: []bool{true, false, false, false},
		},
	},
}

func updateCommonstate() {
	fmt.Println("Hello, World!")
}

func main() {
	// Now you can use exported identifiers from the assigner package
	state := assigner.HRAInput{
		// Initialize the struct
	}

	fmt.Println(state)
	fmt.Println("Hello, World!")
}
