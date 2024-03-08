package assigner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"root/config"
	"root/distributor"
	"root/driver/elevio"
	"root/elevator"
	"runtime"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

func toLocalAssingment(a map[string][][3]bool) elevator.Assingments {
	var ea elevator.Assingments
	L, ok := a[config.Elevator_id]
	if !ok {
		panic("elevator not here -local")
	}

	for f := 0; f < 4; f++ {
		for b := 0; b < 3; b++ {
			ea[f][b] = L[f][b]
		}
	}
	return ea
}

func toLightsAssingment(cs distributor.HRAInput) elevator.Assingments {
	var lights elevator.Assingments
	L, ok := cs.States[config.Elevator_id]
	if !ok {
		panic("elevator not here -lights")
	}
	for f := 0; f < 4; f++ {
		for b := 0; b < 2; b++ {
			lights[f][b] = cs.HallRequests[f][b]

		}
	}
	for f := 0; f < 4; f++ {
		lights[f][elevio.BT_Cab] = L.CabRequests[f]
	}
	return lights
}

func removeUnavailableElevators(cs distributor.HRAInput) (distributor.HRAInput) {
	for k := range cs.States {
		if k != config.Elevator_id && cs.Ackmap[k] == distributor.NotAvailable {
			delete(cs.States, k)
		}
	}
	return cs
}

func Assigner(
	eleveatorAssingmentC chan<- elevator.Assingments,
	lightsAssingmentC chan<- elevator.Assingments,
	messageToAssinger <-chan distributor.HRAInput) {

	for {
		select {
		case cs := <-messageToAssinger:
			fmt.Println("Assigner: Received commonstate")
			distributor.PrintCommonState(cs)
			cs = removeUnavailableElevators(cs)
			localAssingment := toLocalAssingment(CalculateHRA(cs))
			lightsAssingment := toLightsAssingment(cs)
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

	jsonBytes, err := json.Marshal(cs)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		panic("json.Marshal error")
	}

	ret, err := exec.Command("assigner/"+hraExecutable, "-i", "--includeCab", string(jsonBytes)).CombinedOutput()
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

	// fmt.Printf("output: \n")
	// for k, v := range *output {
	// 	fmt.Printf("%6v :  %+v\n", k, v)
	// }

	return *output
}
