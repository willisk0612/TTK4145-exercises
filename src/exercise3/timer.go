package main

import (
	"time"
)

type TimerAction int

const (
	RESET TimerAction = iota
	STOP
)

// Starts, stops or resets a door timer for a specified time
func (elevator *Elevator) DoorTimer(duration *time.Timer, timeout chan bool, action <-chan TimerAction) {
	for {
		select {
		case a := <-action:
			switch a {
			case RESET:
				resetTimer(duration)
			case STOP:
				duration.Stop()
			}
		case <-duration.C:
			timeout <- true
		}
	}
}

// Stops the timer and resets it
func resetTimer(t *time.Timer) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
    t.Reset(DOOR_OPEN_DURATION)
}
