package main

import (
	"main/dev-tools/driver-go/elevio"
	"time"
)

const (
	DOOR_OPEN_DURATION= 3 * time.Second
)



// Initializes elevator with default values. Moves elevator down if between floors.
func InitElevator() Elevator {
	elevator := Elevator{
		Dir:       elevio.MD_Stop,
		Orders:    [elevio.N_FLOORS][elevio.N_BUTTONS]int{},
		Behaviour: elevio.Idle,

	}
	// Move down to find a floor if starting between floors
	if elevio.GetFloor() == -1 {
		elevator.Dir = elevio.MD_Down
		elevator.Behaviour = elevio.Moving
		elevio.SetMotorDirection(elevator.Dir)
	}
	return elevator
}
