package fsm

import (
	"main/dev-tools/driver-go/elevio"
	"time"
)

type Elevator struct {
	Floor       int
	Dir         elevio.MotorDirection
	Orders      [elevio.N_FLOORS][elevio.N_BUTTONS]int
	Behaviour   elevio.ElevatorBehaviour
	Config      elevio.ElevatorConfig
	DoorTimer   *time.Timer
	TimerActive bool
}

type DirnBehaviourPair struct {
	Dir       elevio.MotorDirection
	Behaviour elevio.ElevatorBehaviour
}
