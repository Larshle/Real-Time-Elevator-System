package watchdog

import (
	"time"

)

func Watchdog(seconds int, petC <-chan bool, barkC chan <- bool) {
	TheDawg := time.NewTimer(time.Duration(seconds) * time.Second)
	for {
		select {
		case <- petC:
			TheDawg.Reset(time.Duration(seconds) * time.Second)

		case <- TheDawg.C:
			barkC <- true
			TheDawg.Reset(time.Duration(seconds) * time.Second)
	 	}
 	}	
}