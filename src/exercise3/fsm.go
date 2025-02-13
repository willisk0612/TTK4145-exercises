// Contains helper functions for main.go
package main

import (
	"io"
	"log"
	"main/dev-tools/driver-go/elevio"
)

func init() {
	// Disable logging
	log.SetOutput(io.Discard)
}

// Moves elevator if order floor is different from current floor
func (elevator *Elevator) HandleButtonPress(btn elevio.ButtonEvent, timerAction chan TimerAction) {
	elevator.Orders[btn.Floor][btn.Button] = 1
	elevio.SetButtonLamp(btn.Button, btn.Floor, true)

	switch elevator.Behaviour {
	case elevio.DoorOpen:
		if elevator.Floor == btn.Floor {
			timerAction <- RESET
		}
	case elevio.Moving:
		// Keep moving, orders handled at floor arrival
		return
	case elevio.Idle:
		if elevator.Floor == btn.Floor {
			elevator.Behaviour = elevio.DoorOpen
			elevio.SetDoorOpenLamp(true)
			timerAction <- RESET
		} else {
			pair := elevator.chooseDirection()
			elevator.Dir = pair.Dir
			elevator.Behaviour = pair.Behaviour
			elevio.SetMotorDirection(elevator.Dir)
		}
	//case elevio.Error:
	}
}

// Stops elevator at floor and opens door
func (elevator *Elevator) HandleFloorArrival(floor int, timerAction chan TimerAction) {
	elevator.Floor = floor
	elevio.SetFloorIndicator(floor)
	log.Printf("Arrived at floor %d\n", floor)

	if elevator.shouldStop() {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(true)
		timerAction <- RESET
		log.Println("Door opened and timer reset")
		elevator.clearOrdersAtFloor()
		elevator.Behaviour = elevio.DoorOpen
	} else {
		log.Println("Continuing to next floor")
	}
}

// Stops elevator and opens door
func (elevator *Elevator) HandleObstruction(obstruction bool, timerAction chan TimerAction) {
	elevator.Obstructed = obstruction

	if obstruction {
		elevio.SetMotorDirection(elevio.MD_Stop)
		if elevator.Behaviour == elevio.Moving {
			elevator.Behaviour = elevio.DoorOpen
		}
		elevio.SetDoorOpenLamp(true)
	} else {
		timerAction <- RESET
	}
}

// Stops elevator and clears all orders and button lamps
func (elevator *Elevator) HandleStop() {
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(false)
	for f := 0; f < elevio.N_FLOORS; f++ {
		for b := elevio.ButtonType(0); b < 3; b++ {
			elevator.Orders[f][b] = 0
			elevio.SetButtonLamp(b, f, false)
		}
	}
}

// Handles door timeout with obstruction check
func (elevator *Elevator) HandleDoorTimeout(timerAction chan TimerAction) {
	if elevator.Behaviour != elevio.DoorOpen {
		return
	}
	log.Println("HandleDoorTimeout: Timer expired")
	if elevator.Obstructed {
		log.Println("Door obstructed - keeping open")
		timerAction <- RESET
		return
	}

	log.Println("Closing door")
	elevio.SetDoorOpenLamp(false)
	elevator.clearOrdersAtFloor()
	pair := elevator.chooseDirection()
	elevator.Dir = pair.Dir
	elevator.Behaviour = pair.Behaviour

	if elevator.Behaviour == elevio.Moving {
		elevio.SetMotorDirection(elevator.Dir)
	}
}

// Helper function to count orders
func (elevator *Elevator) countOrders(startFloor int, endFloor int) (result int) {
	for floor := startFloor; floor < endFloor; floor++ {
		for btn := 0; btn < elevio.N_BUTTONS; btn++ {
			if elevator.Orders[floor][btn] != 0 {
				result++
			}
		}
	}
	return result
}

// Counts button orders above
func (elevator *Elevator) ordersAbove() (result int) {
	return elevator.countOrders(elevator.Floor+1, elevio.N_FLOORS)
}

// Counts button orders below
func (elevator *Elevator) ordersBelow() (result int) {
	return elevator.countOrders(0, elevator.Floor)
}

// Counts button orders at current floor
func (elevator *Elevator) ordersHere() (result int) {
	return elevator.countOrders(elevator.Floor, elevator.Floor+1)
}

// Chooses elevator direction based on current orders. Prio: Up > Down > Stop
func (elevator *Elevator) chooseDirection() DirnBehaviourPair {
	switch elevator.Dir {
	case elevio.MD_Up:
		if elevator.ordersAbove() > 0 {
			return DirnBehaviourPair{elevio.MD_Up, elevio.Moving}
		} else if elevator.ordersHere() > 0 {
			return DirnBehaviourPair{elevio.MD_Stop, elevio.DoorOpen}
		} else if elevator.ordersBelow() > 0 {
			return DirnBehaviourPair{elevio.MD_Down, elevio.Moving}
		}
	case elevio.MD_Down:
		if elevator.ordersBelow() > 0 {
			return DirnBehaviourPair{elevio.MD_Down, elevio.Moving}
		} else if elevator.ordersHere() > 0 {
			return DirnBehaviourPair{elevio.MD_Stop, elevio.DoorOpen}
		} else if elevator.ordersAbove() > 0 {
			return DirnBehaviourPair{elevio.MD_Up, elevio.Moving}
		}
	case elevio.MD_Stop:
		if elevator.ordersAbove() > 0 {
			return DirnBehaviourPair{elevio.MD_Up, elevio.Moving}
		} else if elevator.ordersBelow() > 0 {
			return DirnBehaviourPair{elevio.MD_Down, elevio.Moving}
		} else if elevator.ordersHere() > 0 {
			return DirnBehaviourPair{elevio.MD_Stop, elevio.DoorOpen}
		}
	}
	return DirnBehaviourPair{elevio.MD_Stop, elevio.Idle}
}

// Checks if the elevator should stop at the current floor
func (elevator *Elevator) shouldStop() bool {
	currentorders := elevator.Orders[elevator.Floor]
	// Elevator should always stop at top and bottom floor
	if elevator.Floor == 0 || elevator.Floor == elevio.N_FLOORS-1 {
		return true
	}
	switch elevator.Dir {
	case elevio.MD_Down:
		return currentorders[elevio.BT_HallDown] != 0 ||
			currentorders[elevio.BT_Cab] != 0 ||
			elevator.ordersBelow() == 0
	case elevio.MD_Up:
		return currentorders[elevio.BT_HallUp] != 0 ||
			currentorders[elevio.BT_Cab] != 0 ||
			elevator.ordersAbove() == 0
	case elevio.MD_Stop:
		return true
	default:
		return false
	}
}

// Clears orders at current floor
func (elevator *Elevator) clearOrdersAtFloor() {
	switch elevator.Config.ClearOrderVariant {
	case elevio.CV_All:
		for btn := 0; btn < elevio.N_BUTTONS; btn++ {
			elevator.Orders[elevator.Floor][btn] = 0
			elevio.SetButtonLamp(elevio.ButtonType(btn), elevator.Floor, false)
		}
	case elevio.CV_InDirn:
		elevator.clearOrderAndLamp(elevio.BT_Cab)

		switch elevator.Dir {
		case elevio.MD_Up:
			elevator.clearOrderAndLamp(elevio.BT_HallUp)
			if elevator.ordersAbove() == 0 {
				elevator.clearOrderAndLamp(elevio.BT_HallDown)
			}
		case elevio.MD_Down:
			elevator.clearOrderAndLamp(elevio.BT_HallDown)
			if elevator.ordersBelow() == 0 {
				elevator.clearOrderAndLamp(elevio.BT_HallUp)
			}
		}
	}
}

// Helper function to clear orders and lamps
func (elevator *Elevator) clearOrderAndLamp(btn elevio.ButtonType) {
	elevator.Orders[elevator.Floor][btn] = 0
	elevio.SetButtonLamp(btn, elevator.Floor, false)
}
