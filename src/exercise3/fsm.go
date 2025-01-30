// Contains helper functions for main.go
package main

import (
	"main/dev-tools/driver-go/elevio"
	"fmt"
)

// Moves elevator if order floor is different from current floor
func (elevator *Elevator) HandleButtonPress(btn elevio.ButtonEvent) {
	elevator.Orders[btn.Floor][btn.Button] = 1
	elevio.SetButtonLamp(btn.Button, btn.Floor, true)

	switch elevator.Behaviour {
	case elevio.DoorOpen:
		if elevator.Floor == btn.Floor {
			elevator.ResetDoorTimer()
		}
	case elevio.Moving:
		// Keep moving, orders handled at floor arrival
		return
	case elevio.Idle:
		if elevator.Floor == btn.Floor {
			elevator.Behaviour = elevio.DoorOpen
			elevio.SetDoorOpenLamp(true)
			elevator.TimerStart(elevator.Config.DoorOpenDuration)
		} else {
			pair := elevator.chooseDirection()
			elevator.Dir = pair.Dir
			elevator.Behaviour = pair.Behaviour
			elevio.SetMotorDirection(elevator.Dir)
		}
	}
}

// Stops elevator at floor and opens door
func (elevator *Elevator) HandleFloorArrival(floor int) {
	elevator.Floor = floor
	elevio.SetFloorIndicator(floor)

	if elevator.shouldStop() {
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetDoorOpenLamp(true)
		elevator.TimerStart(elevator.Config.DoorOpenDuration)
		elevator.clearOrdersAtFloor()

	}
}

// Stops elevator and opens door
func (elevator *Elevator) HandleObstruction(obstruction bool) {
	elevator.Obstructed = obstruction

	if obstruction {
	elevio.SetMotorDirection(elevio.MD_Stop)
		if elevator.Behaviour == elevio.Moving {
			elevator.Behaviour = elevio.DoorOpen
		}
	elevio.SetDoorOpenLamp(true)
	} else {
		elevator.TimerStart(elevator.Config.DoorOpenDuration)
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
func (elevator *Elevator) HandleDoorTimeout() {
    if elevator.Behaviour != elevio.DoorOpen {
        return
    }
	// Start timer if obstruction is removed
	elevator.ResetDoorTimer()

    if elevator.TimerTimedOut() {
        fmt.Println("HandleDoorTimeout: Timer expired")
        if elevator.Obstructed {
            fmt.Println("Door obstructed - keeping open")
            elevator.ResetDoorTimer()
            return
        }

        fmt.Println("Closing door")
        elevio.SetDoorOpenLamp(false)
        elevator.clearOrdersAtFloor()
        pair := elevator.chooseDirection()
        elevator.Dir = pair.Dir
        elevator.Behaviour = pair.Behaviour

        if elevator.Behaviour == elevio.Moving {
            elevio.SetMotorDirection(elevator.Dir)
        }
    }
}

// Counts button orders above
func (elevator *Elevator) ordersAbove() (result int) {
	for floor := elevator.Floor + 1; floor < elevio.N_FLOORS; floor++ {
		for btn := 0; btn < elevio.N_BUTTONS; btn++ {
			if elevator.Orders[floor][btn] != 0 {
				// Count the number of orders above the current floor
				result++
			}
		}
	}
	return result
}

// Counts button orders below
func (elevator *Elevator) ordersBelow() (result int) {
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < elevio.N_BUTTONS; btn++ {
			if elevator.Orders[floor][btn] != 0 {
				result++
			}
		}
	}
	return result
}

// Counts button orders at current floor
func (elevator *Elevator) ordersHere() (result int) {
	for btn := 0; btn < elevio.N_BUTTONS; btn++ {
		if elevator.Orders[elevator.Floor][btn] != 0 {
			result++
		}
	}
	return result
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
		// Always clear cab calls
		elevator.Orders[elevator.Floor][elevio.BT_Cab] = 0
		elevio.SetButtonLamp(elevio.BT_Cab, elevator.Floor, false)

		switch elevator.Dir {
		case elevio.MD_Up:
			elevator.Orders[elevator.Floor][elevio.BT_HallUp] = 0
			elevio.SetButtonLamp(elevio.BT_HallUp, elevator.Floor, false)
			if elevator.ordersAbove() == 0 {
				elevator.Orders[elevator.Floor][elevio.BT_HallDown] = 0
				elevio.SetButtonLamp(elevio.BT_HallDown, elevator.Floor, false)
			}
		case elevio.MD_Down:
			elevator.Orders[elevator.Floor][elevio.BT_HallDown] = 0
			elevio.SetButtonLamp(elevio.BT_HallDown, elevator.Floor, false)
			if elevator.ordersBelow() == 0 {
				elevator.Orders[elevator.Floor][elevio.BT_HallUp] = 0
				elevio.SetButtonLamp(elevio.BT_HallUp, elevator.Floor, false)
			}
		}
	}
}
