package watchdog

import (
	"fmt"
	"root/config"
	"time"
)

func Watchdog(barkC chan<- bool, startMovingC <-chan bool, stopMovingC <-chan bool) {
	timer := time.NewTimer(time.Duration(config.WatchdogTime) * time.Second)
	timer.Stop()

	for {
		select {
		case <-stopMovingC:
			timer.Stop()
			barkC <- false

		case <-startMovingC:
			timer = time.NewTimer(time.Duration(config.WatchdogTime) * time.Second)
			barkC <- false

		case <-timer.C:
			fmt.Println("IM STUCK!")
			barkC <- true

		}
	}
}
