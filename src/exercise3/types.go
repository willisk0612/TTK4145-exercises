package main

import (
	"main/dev-tools/driver-go/elevio"
	"time"
)

// Elevator type with methods for handling fsm events
type Elevator struct {
	Floor       int
	Dir         elevio.MotorDirection
	Orders      [elevio.N_FLOORS][elevio.N_BUTTONS]int
	Behaviour   elevio.ElevatorBehaviour
	Config      elevio.ElevatorConfig
	DoorTimer   *time.Timer
	TimerChan   chan bool
	TimerActive bool
	Obstructed  bool
}

type DirnBehaviourPair struct {
	Dir       elevio.MotorDirection
	Behaviour elevio.ElevatorBehaviour
}
