package config

import (
	"time"
)

const (
	NumFloors       = 4
	NumElevators    = 3
	NumButtons      = 3
	PeersPortNumber = 58735
	BcastPortNumber = 58750
	DisconnectTime  = 5 * time.Second
	DoorOpenDuration = 2 * time.Second
)
