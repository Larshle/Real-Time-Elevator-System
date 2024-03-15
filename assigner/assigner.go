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
	HallRequests [config.NumFloors][2]bool `json:"hallRequests"`
	States       map[string]hraState       `json:"states"`
}

type hraState struct {
	Behaviour   string                 `json:"behaviour"`
	Floor       int                    `json:"floor"`
	Direction   string                 `json:"direction"`
	CabRequests [config.NumFloors]bool `json:"cabRequests"`
}

func CalculateOptimalAssignments(cs distributor.CommonState, ElevatorID int) elevator.Assignments {

	stateMap := make(map[string]hraState)
	for i, v := range cs.States {
		if cs.Ackmap[i] == distributor.NotAvailable || v.State.Stuck { // Remove not-available and stuck elevators from stateMap
			continue
		} else {
			stateMap[strconv.Itoa(i)] = hraState{
				Behaviour:   v.State.Behaviour.ToString(),
				Floor:       v.State.Floor,
				Direction:   v.State.Direction.ToString(),
				CabRequests: v.CabRequests,
			}
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

	output := new(map[string]elevator.Assignments)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		panic("json.Unmarshal error")
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
 	fmt.Printf("%6v :  %+v\n", k, v)
	}

	return (*output)[strconv.Itoa(ElevatorID)]
}

func ToLightsAssingment(cs distributor.CommonState, ElevatorID int) elevator.Assignments {
	var lights elevator.Assignments

	for f := 0; f < config.NumFloors; f++ {
		for b := 0; b < 2; b++ {
			lights[f][b] = cs.HallRequests[f][b]
		}
	}

	for f := 0; f < config.NumFloors; f++ {
		lights[f][elevio.BT_Cab] = cs.States[ElevatorID].CabRequests[f]
	}

	return lights
}
