package assigner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"root/distributor"
	"root/elevator"
	"root/elevio"
	"runtime"
	"strconv"
)

func ToLocalAssingment(a map[string][][3]bool, ElevatorID int) elevator.Assignments {
	var ea elevator.Assignments
	L, ok := a[strconv.Itoa(ElevatorID)]
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

func ToLightsAssingment(cs distributor.CommonState, ElevatorID int) elevator.Assignments {
	var lights elevator.Assignments

	L := cs.States[ElevatorID]

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

// func RemoveUnavailableElevators(cs distributor.CommonState, ElevatorID int) distributor.CommonState {

// 	var newSlice []distributor.LocalElevState

// 	for index := range cs.States {
// 		if index != ElevatorID && cs.Ackmap[ElevatorID] == distributor.NotAvailable{
// 			continue
// 		} else{
// 			newSlice
// 		}

// 		if index != ElevatorID && cs.Ackmap[ElevatorID] == distributor.NotAvailable{
// 			newSlice = append(cs.States[:index], cs.States[index+1:]...)
// 			fmt.Println("Assigner: Removed unavailable elevators")
// 		}
// 	}
// 	cs.States = newSlice
// 	return cs
// }

type hra struct {
	HallRequests [][2]bool                             `json:"hallRequests"`
	States       map[string]distributor.LocalElevState `json:"states"`
}

func CalculateHRA(cs distributor.CommonState) map[string][][3]bool {

	m := make(map[string]distributor.LocalElevState)
	for i, v := range cs.States {
		if cs.Ackmap[i] == distributor.NotAvailable {
			continue
		} else {
			m[strconv.Itoa(i)] = v
		}
	}

	newHRA := hra{cs.HallRequests, m}

	// fmt.Println("\n")
	// fmt.Println(cs)
	// fmt.Println("\n")

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "darwin":
		hraExecutable = "hall_request_assigner_mac"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	jsonBytes, err := json.Marshal(newHRA)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
		panic("json.Marshal error")
	}

	ret, err := exec.Command("assigner/executables/"+hraExecutable, "-i", "--includeCab", string(jsonBytes)).CombinedOutput()
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

	//fmt.Printf("output: \n")
	//for k, v := range *output {
	//fmt.Printf("%6v :  %+v\n", k, v)
	//}

	return *output
}
