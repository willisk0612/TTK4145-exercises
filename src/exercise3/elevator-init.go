package main

import (
	"main/dev-tools/driver-go/elevio"
	"time"
)

// Initializes elevator with default values. Moves elevator down if between floors.
func InitElevator() Elevator {
	elevator := Elevator{
		Floor:       0,
		Dir:         elevio.MD_Stop,
		Orders:      [elevio.N_FLOORS][elevio.N_BUTTONS]int{},
		Behaviour:   elevio.Idle,
		Config: elevio.ElevatorConfig{
			DoorOpenDuration: 3 * time.Second,
		},
		DoorTimer:   time.NewTimer(0),
		TimerChan:   make(chan bool),
	}
	// Move down to find a floor if starting between floors
	if elevio.GetFloor() == -1 {
		elevator.Dir = elevio.MD_Down
		elevator.Behaviour = elevio.Moving
		elevio.SetMotorDirection(elevator.Dir)
	}
	return elevator
}
