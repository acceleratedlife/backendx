package main

import (
	"time"
)

// schedule a function to be called every minute
// returns channel to stop repeating
// to set up
//  runEveryMinute(func(time2 time.Time) {
//	  lgr.Printf("INFO Tick2 at", time2)
//   })
// to stop:   done <- true

func runEveryMinute(handler func(time2 time.Time)) (done chan bool) {
	ticker := time.NewTicker(60 * time.Second)
	done = make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				handler(t)
			}
		}
	}()

	return done
}
