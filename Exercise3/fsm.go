package fsm

import (
	"fmt"
	"time"
	"../Driver-go/elevio"
)

type Elevator struct {
    Floor     int
    Dir      int
    Requests  [elevio.N_FLOORS][elevio.N_BUTTONS]int
    Behaviour elevio.ElevatorBehaviour
    Config    elevio.ElevatorConfig
	DoorTimer *time.Timer
    timerActive bool
}


func (e *Elevator) RequestsAbove(floor int) bool {
    for f := floor + 1; f < elevio.N_FLOORS; f++ {
        for btn := 0; btn < elevio.N_BUTTONS; btn++ {
            if e.Requests[f][btn] != 0 {
                return true
            }
        }
    }
    return false
}

func (e *Elevator) RequestsBelow(floor int) bool {
    for f := 0; f < floor; f++ {
        for btn := 0; btn < elevio.N_BUTTONS; btn++ {
            if e.Requests[f][btn] != 0 {
                return true
            }
        }
    }
    return false
}

func (e *Elevator) RequestsHere(floor int) bool {
    for btn := 0; btn < elevio.N_BUTTONS; btn++ {
        if e.Requests[floor][btn] != 0 {
            return true
        }
    }
    return false
}

func (e *Elevator) chooseDirection(floor int) (int, elevio.ElevatorBehaviour) {
    if e.RequestsHere(floor) {
        return 0, elevio.DoorOpen
    } else if e.RequestsAbove(floor) {
        return 1, elevio.Moving
    } else if e.RequestsBelow(floor) {
        return -1, elevio.Moving
    }
    
    return 0, elevio.Idle
}

func (e *Elevator) timerStart(duration time.Duration) {
    e.resetDoorTimer()
    e.timerActive = true
    e.DoorTimer = time.NewTimer(duration)
}

func (e *Elevator) timerStop() {
    if e.DoorTimer != nil {
        e.DoorTimer.Stop()
    }
    e.timerActive = false
}

func (e *Elevator) timerTimedOut() bool {
    if !e.timerActive {
        return false
    }
    select {
    case <-e.DoorTimer.C:
        e.timerActive = false
        return true
    default:
        return false
    }
}

func (e *Elevator) resetDoorTimer() {
    if e.DoorTimer != nil {
        e.DoorTimer.Stop()
    }
    // set doortimer after openduration reaches onDoorTimeout
    e.timerStart(e.Config.DoorOpenDuration)
}

    

func (e *Elevator) OnDoorTimeout() {
    fmt.Println("Door timeout")
    switch e.Behaviour {
    case elevio.DoorOpen:
        e.Dir, e.Behaviour = e.chooseDirection(e.Floor)
        
        switch e.Behaviour {
        case elevio.Moving:
            elevio.SetDoorOpenLamp(false)
            elevio.SetMotorDirection(elevio.MotorDirection(e.Dir))
        case elevio.DoorOpen:
            e.resetDoorTimer()
        case elevio.Idle:
            elevio.SetDoorOpenLamp(false)
        }
    }
}