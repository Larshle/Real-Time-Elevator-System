package peers

import (
	"root/network/conn"
	"fmt"
	"net"
	"sort"
	"strconv"
	"time"
)

type PeerUpdate struct {
	Peers []int // Changed from []string to []int
	New   int   // Changed from string to int
	Lost  []int // Changed from []string to []int
}

const interval = 15 * time.Millisecond
const timeout = 500 * time.Millisecond

func Transmitter(port int, id int, transmitEnable <-chan bool) { // Changed id type to int

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			idStr := strconv.Itoa(id) // Convert int ID to string for transmission
			conn.WriteTo([]byte(idStr), addr)
		}
	}
}

func Receiver(port int, peerUpdateCh chan<- PeerUpdate) {

	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[int]time.Time) // Key changed from string to int

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		idStr := string(buf[:n])
		id, err := strconv.Atoi(idStr) // Convert received ID from string to int
		if err != nil {
			continue // Skip invalid IDs
		}

		// Adding new connection
		p.New = 0
		if id != 0 {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// Removing dead connection
		p.Lost = make([]int, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {
			p.Peers = make([]int, 0, len(lastSeen))

			for k := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Ints(p.Peers) // Use sort.Ints for integer slices
			sort.Ints(p.Lost)
			peerUpdateCh <- p
		}
	}
}
