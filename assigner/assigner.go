package assigner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"root/driver/elevio"
	"root/elevator"
	"root/distributor"
	"runtime"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase



func toLocalAssingment(a map[string][][3]bool, elevatorID string) elevator.Assingments {
    var ea elevator.Assingments
    L, ok := a[elevatorID]
    if !ok {
        panic("elevator not here")
    }

	for f := 0; f < 4; f++ {
		for b := 0; b < 3; b++ {
			ea[f][b] = L[f][b]
		}
	}		
    return ea
}

func toLightsAssingment(cs distributor.HRAInput, elevatorID string) elevator.Assingments {
	var lights elevator.Assingments
	L, ok := cs.States[elevatorID]
    if !ok {
        panic("elevator not here")
    }
	for f := 0; f < 4; f++ {
		for b := 0; b < 2; b++ {
			lights[f][b] = cs.HallRequests[f][b]
			
		}
	}
	for f:= 0; f < 4; f++ {
		lights[f][elevio.BT_Cab] = L.CabRequests[f]
	}
	return lights
}

func Assingner(eleveatorAssingmentC chan<- elevator.Assingments, lightsAssingmentC chan<- elevator.Assingments , csToAssingerC <-chan HRAInput, elevatorID string){
	// MÅ finne noe her for å få tak i elevatorID
	// Må ha bruke noe for å gjøre om  fra cs til enkel order

	for{
		select{
		case cs := <- csToAssingerC:
			localAssingment := toLocalAssingment( CalculateHRA(cs), elevatorID)
			lightsAssingment:= toLightsAssingment(cs, elevatorID)
			lightsAssingmentC <- lightsAssingment
			eleveatorAssingmentC <- localAssingment
			
			
		}
	}
}

func CalculateHRA(cs distributor.HRAInput) map[string][][3]bool {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux", "darwin":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	input := distributor.HRAInput{
		ID: 1,
		Origin: "string",
		Ackmap: map[string]distributor.Ack_status {"ein": distributor.Acked, "to": distributor.Acked},
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
		States: map[string]distributor.HRAElevState{
			"one":{
				Behaviour:    "moving",
				Floor:       2,
				Direction:   "up",
				CabRequests: []bool{true, true, false, false},
			},
			"two":{
				Behaviour:    "idle",
				Floor:       3,
				Direction:   "stop",
				CabRequests: []bool{false, false, false, false},
			},
		},
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		panic("json.Marshal error")
	}

	ret, err := exec.Command("../assigner/"+hraExecutable, "-i", "--includeCab", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		panic("exec.Command error")
	}

	output := new(map[string][][3]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		panic("json.Unmarshal error")
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
	
	return *output
}
