// Keeps track of the events of the elevator system
package fsm

import (
	"main/dev-tools/driver-go/elevio"
)

func (elev *Elevator) OrdersAbove() (result int) {
	for floor := elev.Floor + 1; floor < elevio.N_FLOORS; floor++ {
		for btn := 0; btn < elevio.N_BUTTONS; btn++ {
			if elev.Orders[floor][btn] > 0 {
				// Count the number of orders above the current floor
				result++
			}
		}
	}
	return result
}

func (elev *Elevator) OrdersBelow() (result int) {
	for floor := 0; floor < elev.Floor; floor++ {
		for btn := 0; btn < elevio.N_BUTTONS; btn++ {
			if elev.Orders[floor][btn] > 0 {
				result++
			}
		}
	}
	return result
}

func (elev *Elevator) OrdersHere() (result int) {
	for btn := 0; btn < elevio.N_BUTTONS; btn++ {
		if elev.Orders[elev.Floor][btn] > 0 {
			result++
		}
	}
	return result
}

func (elev *Elevator) ChooseDirection() DirnBehaviourPair {
	switch elev.Dir {
	case elevio.MD_Stop:
		// If we have orders above, go up
		if elev.OrdersAbove() > 0 {
			return DirnBehaviourPair{elevio.MD_Up, elevio.Moving}
		}
		// If we have orders below, go down
		if elev.OrdersBelow() > 0 {
			return DirnBehaviourPair{elevio.MD_Down, elevio.Moving}
		}
		// If we have orders here, stay
		if elev.OrdersHere() > 0 {
			return DirnBehaviourPair{elevio.MD_Stop, elevio.DoorOpen}
		}
	}
	return DirnBehaviourPair{elevio.MD_Stop, elevio.Idle}
}

func (elev *Elevator) ShouldStop() bool {
	switch elev.Dir {
	case elevio.MD_Down:
		if elev.Orders[elev.Floor][elevio.BT_HallDown] > 0 ||
			elev.Orders[elev.Floor][elevio.BT_Cab] > 0 ||
			elev.OrdersBelow() == 0 {
			return true
		}
		return false
	case elevio.MD_Up:
		if elev.Orders[elev.Floor][elevio.BT_HallUp] > 0 ||
			elev.Orders[elev.Floor][elevio.BT_Cab] > 0 ||
			elev.OrdersAbove() == 0 {
			return true
		}
		return false
	case elevio.MD_Stop:
		return true
	}
	return false
}
