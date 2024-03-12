package assigner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"root/distributor"
	"root/elevator"
	"runtime"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

func ToLocalAssingment(a map[int][][3]bool, ElevatorID int) elevator.Assignments {
	var ea elevator.Assignments
	L, ok := a[ElevatorID]
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


func RemoveUnavailableElevators(cs distributor.CommonState, ElevatorID int) distributor.CommonState {
	for k := range cs.States {
		if k != ElevatorID && cs.Ackmap[k] == distributor.NotAvailable {
			delete(cs.States, k)
		}
	}

	return cs
}


func CalculateHRA(cs distributor.CommonState) map[int][][3]bool {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner_linux"
	case "darwin":
		hraExecutable = "hall_request_assigner_mac"
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

	ret, err := exec.Command("assigner/executables/"+hraExecutable, "-i", "--includeCab", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		panic("exec.Command error")
	}

	output := new(map[int][][3]bool)
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
