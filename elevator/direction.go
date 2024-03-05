package elevator

import (
	"root/driver/elevio"
)

type Direction int

const (
	Up Direction = iota
	Down
	Stop
)

func (d Direction) toMD() elevio.MotorDirection{
	return map[Direction]elevio.MotorDirection{Up: elevio.MD_Up, Down: elevio.MD_Down}[d]
}

func (d Direction) toBT() elevio.ButtonType{
	return map[Direction]elevio.ButtonType{Up: elevio.BT_HallUp, Down: elevio.BT_HallDown}[d]
}

func (d Direction) toOpposite() Direction{
	return map[Direction]Direction{Up: Down, Down: Up}[d]

}


func (d Direction) ToString() string{
	return map[Direction]string{Up: "up", Down: "down", Stop: "stop"}[d]
}

