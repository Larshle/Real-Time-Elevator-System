// to ensure packedges are sent and notified if lost
package handShake

import (
	"fmt"
	"root/network/network_modules/bcast"
	"root/network/network_modules/localip"
	"root/network/network_modules/peers"
)

type SYN_ACK struct {
	synNum int
	ackNum int
}
func watch(ElvTransEnable chan bool, SYN_ACK chan SYN_ACK, peerUpdate chan peers.PeerUpdate, id string) {
	for {
		select {
		case <-peerUpdate:
			//fmt.Println("PeerUpdate")
			ElvTransEnable <- false
		case <-SYN_ACK:
			//fmt.Println("SYN_ACK")
			ElvTransEnable <- false
		default:
			//fmt.Println("Default")
			ElvTransEnable <- true
		}
	}
}
SYN_ACK := make(chan SYN_ACK)


peerUpdate := make(chan peers.PeerUpdate)
	
ElvTransEnable := make(chan bool)

go peers.Transmitter(15600, id, ElvTransEnable)
go peers.Receiver(15600, peerUpdate)


go bcast.Transmitter(16969, SYN_ACK)
go bcast.Receiver(16969, SYN_ACK)










