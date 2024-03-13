package watchdog

import (
	"time"
	"fmt"
)

func Watchdog(seconds int, barkC chan<- bool, startMovingC <-chan bool, stopMovingC <-chan bool) {
	timer := time.NewTimer(time.Duration(seconds) * time.Second)
	timer.Stop() 

	for {
		select {
		case <-stopMovingC:
			timer.Stop() 
			//fmt.Println("Stopped Moving: Timer reset.")
			barkC <- false

		case <-startMovingC:
			timer = time.NewTimer(time.Duration(seconds)*time.Second)
			//fmt.Println("Started Moving: Timer started.")
			barkC <- false

		case <-timer.C:
			fmt.Println("IM STUCK!!!")
			barkC <- true	
		}
	}
}
