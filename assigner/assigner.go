package assigner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
	Behaviour   string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	ID int
	Origin string
	Ackmap map[string]int
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func main() {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux", "darwin":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	input := HRAInput{
		ID: 1,
		Origin: "string",
<<<<<<< HEAD
		Ackmap: map[string]int {"ein": "true", "to": "false"},
		HallRequests: [][2]bool{{false, false}, {true, false}, {false, false}, {false, true}},
=======
		Ackmap: map[string]string {"ein": "true", "to": "false"},
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
>>>>>>> c270e87a40ee3cc388cea87420d6e5fb9ae8378d
		States: map[string]HRAElevState{
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
		return
	}

	ret, err := exec.Command("../assigner/"+hraExecutable, "-i", "--includeCab", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
		return
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
		return
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
}
