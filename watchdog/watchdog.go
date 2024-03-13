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
			barkC <- false

		case <-startMovingC:
			timer = time.NewTimer(time.Duration(seconds)*time.Second)
			barkC <- false

		case <-timer.C:
			fmt.Println("IM STUCK!!!")
			barkC <- true	
		}
	}
}
