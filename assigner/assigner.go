package assigner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"root/config"
	"root/distributor"
	"root/elevator"
	"root/elevio"
	"runtime"
	"strconv"
)

type CalculateOptimalAssignmentsFormat struct {
	HallRequests [config.NumFloors][2]bool             `json:"hallRequests"`
	States       map[string]distributor.LocalElevState `json:"states"`
}

func CalculateOptimalAssignments(cs distributor.CommonState, ElevatorID int) elevator.Assignments {

	// Convert from []LocalElevState to map[string]LocalElevState
	stateMap := make(map[string]distributor.LocalElevState)
	for i, v := range cs.States {
		if cs.Ackmap[i] == distributor.NotAvailable {
			continue
		} else {
			stateMap[strconv.Itoa(i)] = v
		}
	}

	hall_request_assignerInput := CalculateOptimalAssignmentsFormat{cs.HallRequests, stateMap}

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

	jsonBytes, err := json.Marshal(hall_request_assignerInput)
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

	outputContent:= *output

	var elevatorAssignments elevator.Assignments
	L, ok := outputContent[strconv.Itoa(ElevatorID)]

	if !ok {
		panic("elevator not here -local")
	}

	for f := 0; f < config.NumFloors; f++ {
		for b := 0; b < 3; b++ {
			elevatorAssignments[f][b] = L[f][b]
		}
	}
	return elevatorAssignments
}

func ToLocalAssingment(a map[string][][config.NumButtons]bool, ElevatorID int) elevator.Assignments {
	var ea elevator.Assignments
	L, ok := a[strconv.Itoa(ElevatorID)]

	if !ok {
		panic("elevator not here -local")
	}

	for f := 0; f < config.NumFloors; f++ {
		for b := 0; b < 3; b++ {
			ea[f][b] = L[f][b]
		}
	}
	return ea
}

func ToLightsAssingment(cs distributor.CommonState, ElevatorID int) elevator.Assignments {
	var lights elevator.Assignments
	myState := cs.States[ElevatorID]

	for f := 0; f < config.NumFloors; f++ {
		for b := 0; b < 2; b++ {
			lights[f][b] = cs.HallRequests[f][b]
		}
	}

	for f := 0; f < config.NumFloors; f++ {
		lights[f][elevio.BT_Cab] = myState.CabRequests[f]
	}

	return lights
}
