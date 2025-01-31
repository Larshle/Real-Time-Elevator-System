package elevator

import (
	"root/elevio"
)

type Direction int

const (
	Down Direction = iota
	Up
)

func (d Direction) toMD() elevio.MotorDirection {
	return map[Direction]elevio.MotorDirection{Up: elevio.MD_Up, Down: elevio.MD_Down}[d]
}

func (d Direction) toBT() elevio.ButtonType {
	return map[Direction]elevio.ButtonType{Up: elevio.BT_HallUp, Down: elevio.BT_HallDown}[d]
}

func (d Direction) Opposite() Direction {
	return map[Direction]Direction{Up: Down, Down: Up}[d]
}	


