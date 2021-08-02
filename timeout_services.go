package main

import "time"

func calcTime(perfectTime time.Duration) time.Duration {
	var timeOut time.Duration
	if perfectTime == 0 {
		timeOut = time.Minute * 15
	} else if perfectTime.Minutes() > 5 {
		timeOut = time.Duration(perfectTime.Minutes() * 1.2)
	} else if perfectTime.Minutes() > 1 {
		timeOut = time.Duration(perfectTime.Minutes() * 1.5)
	} else {
		timeOut = perfectTime*2 + (2 * time.Second)
	}
	return timeOut
}
