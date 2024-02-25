package assigner

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"root/driver/elevio"
	"root/elevator"
	"runtime"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type Ack_status int
const (
	NotAcked Ack_status = iota
	Acked
	NotAvailable
)

type HRAElevState struct {
	Behaviour   string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	ID int
	Origin string
	Ackmap map[string]Ack_status
	HallRequests [][2]bool               `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

func (a HRAInput) toLocalAssingment(elev string) elevator.Assingments {
    var ea elevator.Assingments
    L, ok := a.States[elev]
    if !ok {
        panic("elevator not here")
    }
    for c := 0; c < elevio.NumFloors; c++ {
        if c < len(L.CabRequests) {
            ea[elevio.BT_Cab][c] = L.CabRequests[c]
        }
        ea[c][elevio.BT_HallUp] = a.HallRequests[c][elevio.BT_HallUp]
        ea[c][elevio.BT_HallDown] = a.HallRequests[c][elevio.BT_HallDown]
    }
    return ea
}

func Assingner(eleveatorAssingmentC chan<- elevator.Assingments, lightsAssingmentC chan<- elevio.ButtonEvent, csToAssingerC <-chan HRAInput){
	var cs HRAInput
	var elevatorID string

	// MÅ finne noe her for å få tak i elevatorID
	// Må ha bruke noe for å gjøre om  fra cs til enkel order

	for{
		select{
		case cs := <- csToAssingerC:
			localAssingment := cs.toLocalAssingment(elevatorID)
			eleveatorAssingmentC <- localAssingment
			
		}
	}
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
		Ackmap: map[string]Ack_status {"en": Acked, "to": NotAcked, "tre": NotAvailable},
		HallRequests: [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}},
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

	output := new(map[string][][3]bool)
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
