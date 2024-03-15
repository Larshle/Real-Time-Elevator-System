<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
package main

import (
	"flag"
	"fmt"
	"root/assigner"
	"root/config"
	"root/distributor"
	"root/elevator"
	"root/elevio"
	"root/lights"
	"root/network/bcast"
	"root/network/peers"
	"strconv"
)

var Port int
var ElevatorID int

func main() {

	port := flag.Int("port", 15301, "<-- Default verdi, men kan overskrives som en command line argument ved bruk av -port=xxxxx")
	id := flag.Int("id", 0, "<-- Default verdi, men kan overskrives som en command line argument ved bruk av -id=xxxxx")
	flag.Parse()

	Port = *port
	ElevatorID = *id

	fmt.Println()
	elevio.Init("localhost:"+strconv.Itoa(Port), config.NumFloors)

	fmt.Println("Elevator initialized with ID", ElevatorID, "on port", Port)
	fmt.Println("System has", config.NumFloors, "floors and", config.NumElevators, "elevators.")

	newAssignmentC := make(chan elevator.Assignments, 10000)
	deliveredAssignmentC := make(chan elevio.ButtonEvent, 10000)
	newLocalElevStateC := make(chan elevator.State, 10000)
	giverToNetworkC := make(chan distributor.CommonState, 10000)      // Endre navn?
	receiverFromNetworkC := make(chan distributor.CommonState, 10000) // Endre navn? ja faktisk
	toAssignerC := make(chan distributor.CommonState, 10000)
	receiverPeersC := make(chan peers.PeerUpdate, 10000) // Endre navn?
	giverPeersC := make(chan bool, 10000)

	go peers.Receiver(config.PeersPortNumber, receiverPeersC)
	go peers.Transmitter(config.PeersPortNumber, ElevatorID, giverPeersC)

	go bcast.Receiver(config.BcastPortNumber, receiverFromNetworkC)
	go bcast.Transmitter(config.BcastPortNumber, giverToNetworkC)

	go distributor.Distributor(
		deliveredAssignmentC,
		newLocalElevStateC,
		giverToNetworkC,
		receiverFromNetworkC,
		toAssignerC,
		receiverPeersC,
		ElevatorID)

	go elevator.Elevator(
		newAssignmentC,
		newLocalElevStateC,
		deliveredAssignmentC)

	for {
		select {
		case cs := <-toAssignerC:
			localAssingment := assigner.CalculateOptimalAssignments(cs, ElevatorID)
			newAssignmentC <- localAssingment
			lights.SetLights(cs, ElevatorID)

		default:
			continue
		}

	}
}
=======
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
// midlertidig main fil for å teste coden vår

package main

import (

)


func main() {


<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
}
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
}
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
}
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
}
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
}
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
}
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
}
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
}
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
}
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======
}
>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
