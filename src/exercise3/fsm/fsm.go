// Finite state machine logic for the elevator system
package fsm

import (
	"time"
)

func (elev *Elevator) TimerStart(duration time.Duration) {
	elev.ResetDoorTimer()
	elev.TimerActive = true
	elev.DoorTimer = time.NewTimer(duration)
}

func (elev *Elevator) TimerStop() {
	if elev.DoorTimer != nil {
		elev.DoorTimer.Stop()
	}
	elev.TimerActive = false
}

func (elev *Elevator) TimerTimedOut() bool {
	if !elev.TimerActive {
		return false
	}
	select {
	case <-elev.DoorTimer.C:
		elev.TimerActive = false
		return true
	default:
		return false
	}
}

func (elev *Elevator) ResetDoorTimer() {
	if elev.DoorTimer != nil {
		elev.DoorTimer.Stop()
	}
	// set doortimer after openduration reaches onDoorTimeout
	elev.TimerStart(elev.Config.DoorOpenDuration)
}
