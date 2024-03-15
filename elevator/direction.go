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

func (d Direction) toOpposite() Direction {
	return map[Direction]Direction{Up: Down, Down: Up}[d]
}	

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
<<<<<<< HEAD
<<<<<<< HEAD
func (d Direction) ToString() string {
	return map[Direction]string{Up: "up", Down: "down"}[d]
}
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
=======

>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======

>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
=======

>>>>>>> parent of 34a4414 (Merge pull request #2 from Larshle/UpdatingNEWassignments)
