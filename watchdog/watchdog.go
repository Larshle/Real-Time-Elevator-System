package watchdog

import (
	"fmt"
	"time"
)

func MotorWatchdog(WatchdogTime time.Duration, barkC chan<- bool, startMovingC <-chan bool, stopMovingC <-chan bool) {
	timer := time.NewTimer(WatchdogTime)
	timer.Stop()

	for {
		select {
		case <-stopMovingC:
			timer.Stop()
			barkC <- false

		case <-startMovingC:
			timer = time.NewTimer(WatchdogTime)
			barkC <- false

		case <-timer.C:
			fmt.Println("IM STUCK!")
			barkC <- true
		}
	}
}
