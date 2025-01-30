package main

import (
	"time"
	"fmt"
)

func (elevator *Elevator) TimerStart(duration time.Duration) {
    elevator.TimerStop() // Ensure old timer is stopped
    elevator.DoorTimer = time.NewTimer(duration)
    elevator.TimerActive = true
    fmt.Println("Door timer started")
}

func (elevator *Elevator) TimerStop() {
    if elevator.DoorTimer != nil {
        elevator.DoorTimer.Stop()
        select {
        case <-elevator.DoorTimer.C: // Drain channel
        default:
        }
    }
    elevator.TimerActive = false
}

func (elevator *Elevator) TimerTimedOut() bool {
    if !elevator.TimerActive || elevator.DoorTimer == nil {
        return false
    }

    select {
    case <-elevator.DoorTimer.C:
        fmt.Println("Timer timed out - door should close")
        elevator.TimerActive = false
        return true
    default:
        return false
    }
}

func (elevator *Elevator) ResetDoorTimer() {
    elevator.TimerStart(elevator.Config.DoorOpenDuration)
}
