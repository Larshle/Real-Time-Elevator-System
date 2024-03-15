package distributor

import (
	"fmt"
	"root/config"
	"root/elevator"
	"root/elevio"
	"root/network/peers"
	"time"
)

type StashType int

const (
	None StashType = iota
	RemoveCall
	AddCall
	StateChange
)

func Distributor(
	deliveredAssignmentC <-chan elevio.ButtonEvent,
	newLocalStateC <-chan elevator.State,
	networkTx chan<- CommonState,
	networkRx <-chan CommonState,
	confirmedCommonstateC chan<- CommonState,
	peersC <-chan peers.PeerUpdate,
	id int) {

	elevioOrdersC := make(chan elevio.ButtonEvent, 10000)

	go elevio.PollButtons(elevioOrdersC)

	var stateStash elevator.State
	var newOrderStash elevio.ButtonEvent
	var removeOrderStash elevio.ButtonEvent
	var stashType StashType
	var peers peers.PeerUpdate
	var cs CommonState

	disconnectTimer := time.NewTimer(config.DisconnectTime)
	heartbeatTimer := time.NewTicker(config.HeartbeatTime)

	idle := true
	aloneOnNetwork := false

	for {
		select {
		case <-disconnectTimer.C:
			cs.makeOthersUnavailable(id)
			fmt.Println("Lost connection to network")
			aloneOnNetwork = true

		case P := <-peersC:
			peers = P
			cs.makeOthersUnavailable(id)
			idle = false

		case <-heartbeatTimer.C:
			networkTx <- cs

		default:
		}

		switch {
		case idle:
			select {
			case newOrder := <-elevioOrdersC:
				stashType = AddCall
				newOrderStash = newOrder
				cs.prepNewCs(id)
				cs.addAssignments(newOrder, id)
				cs.Ackmap[id] = Acked
				idle = false

			case removeOrder := <-deliveredAssignmentC:
				stashType = RemoveCall
				removeOrderStash = removeOrder
				cs.prepNewCs(id)
				cs.removeAssignments(removeOrder, id)
				cs.Ackmap[id] = Acked
				idle = false

			case newLocalState := <-newLocalStateC:
				stashType = StateChange
				stateStash = newLocalState
				cs.prepNewCs(id)
				cs.updateLocalState(newLocalState, id)
				cs.Ackmap[id] = Acked
				idle = false

			case arrivedCs := <-networkRx:
				disconnectTimer = time.NewTimer(config.DisconnectTime)
				if (arrivedCs.Origin > cs.Origin && arrivedCs.Seq == cs.Seq) || arrivedCs.Seq > cs.Seq {
					cs = arrivedCs
					cs.makeLostPeersUnavailable(peers)
					cs.Ackmap[id] = Acked
					idle = false
				}

			default:
			}

		case aloneOnNetwork:
			select {
			case <-networkRx:
				if cs.States[id].CabRequests == [config.NumFloors]bool{} {
					fmt.Println("Regained connection to network")
					aloneOnNetwork = false
				} else {
					cs.Ackmap[id] = NotAvailable
				}

			case newOrder := <-elevioOrdersC:
				if cs.States[id].State.Motorstop {
					break
				}
				cs.Ackmap[id] = Acked
				cs.addCabCall(newOrder, id)
				confirmedCommonstateC <- cs

			case removeOrder := <-deliveredAssignmentC:
				cs.Ackmap[id] = Acked
				cs.removeAssignments(removeOrder, id)
				confirmedCommonstateC <- cs

			case newLocalState := <-newLocalStateC:
				if !(newLocalState.Obstructed || newLocalState.Motorstop) {
					cs.Ackmap[id] = Acked
					cs.updateLocalState(newLocalState, id)
					confirmedCommonstateC <- cs
				}

			default:
			}

		case !idle:
			select {
			case arrivedCs := <-networkRx:
				if arrivedCs.Seq < cs.Seq {
					break
				}
				disconnectTimer = time.NewTimer(config.DisconnectTime)

				switch {
				case (arrivedCs.Origin > cs.Origin && arrivedCs.Seq == cs.Seq) || arrivedCs.Seq > cs.Seq:
					cs = arrivedCs
					cs.Ackmap[id] = Acked
					cs.makeLostPeersUnavailable(peers)

				case arrivedCs.fullyAcked(id):
					cs = arrivedCs
					confirmedCommonstateC <- cs
					switch {
					case cs.Origin != id && stashType != None:
						cs.prepNewCs(id)

						switch stashType {
						case AddCall:
							cs.addAssignments(newOrderStash, id)
							cs.Ackmap[id] = Acked

						case RemoveCall:
							cs.removeAssignments(removeOrderStash, id)
							cs.Ackmap[id] = Acked

						case StateChange:
							cs.updateLocalState(stateStash, id)
							cs.Ackmap[id] = Acked
						}

					case cs.Origin == id && stashType != None:
						stashType = None
						idle = true

					default:
						idle = true
					}

				case cs.equals(arrivedCs):
					cs = arrivedCs
					cs.Ackmap[id] = Acked
					cs.makeLostPeersUnavailable(peers)

				default:
				}
			default:
			}
		}
	}
}
