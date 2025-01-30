package main

import (
	"main/dev-tools/driver-go/elevio"
)

// Elevator type with methods for handling fsm events
type Elevator struct {
	Floor       int
	Dir         elevio.MotorDirection
	Orders      [elevio.N_FLOORS][elevio.N_BUTTONS]int
	Behaviour   elevio.ElevatorBehaviour
	Config      elevio.ElevatorConfig
	Obstructed  bool
}

type DirnBehaviourPair struct {
	Dir       elevio.MotorDirection
	Behaviour elevio.ElevatorBehaviour
}
