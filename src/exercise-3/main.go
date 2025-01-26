/*
features:
1. Uttilizaes fsm package to create single elevator logic
2. Prints events from the driver
3. Stops the elevator motor if obstruction switch is active, and opens the door
4. Stops the elevator and turns off all lights if stop button is pressed
*/
package main

import (
	"fmt"
	"main/dev-tools/driver-go/elevio"
	"main/src/exercise-3/fsm"
)

func main() {
	numFloors := 4
	elevio.Init("localhost:15657", numFloors)

	elevator := fsm.Elevator{
		Floor:     0,
		Dir:       elevio.MD_Stop,
		Behaviour: elevio.Idle,
		Orders:    [elevio.N_FLOORS][elevio.N_BUTTONS]int{}, // Initialize orders array
	}

	// Get initial floor
	initFloor := elevio.GetFloor()
	if initFloor != -1 {
		elevator.Floor = initFloor
	}

	// Same channels as before
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	fmt.Println("Driver started")
		for {
		select {
			case btn := <-drv_buttons:
			fmt.Printf("Button press: %+v\n", btn)
			fmt.Printf("Current floor: %d\n", elevator.Floor)
			elevator.Orders[btn.Floor][btn.Button] = 1
			elevio.SetButtonLamp(btn.Button, btn.Floor, true)

			dir := elevator.ChooseDirection()
			var dirStr string
			if dir.Dir == elevio.MD_Up {
				dirStr = "Up"
			} else if dir.Dir == elevio.MD_Down {
				dirStr = "Down"
			} else {
				dirStr = "Stop"
			}
			fmt.Printf("Direction chosen: %s\n", dirStr)

			elevator.Dir = dir.Dir
			elevator.Behaviour = dir.Behaviour
			elevio.SetMotorDirection(dir.Dir)

			// Always check for new direction when order received
			if elevator.Behaviour == elevio.Idle {
				dir := elevator.ChooseDirection()
				elevator.Dir = dir.Dir
				elevator.Behaviour = dir.Behaviour
				if dir.Dir != elevio.MD_Stop {
					elevio.SetMotorDirection(dir.Dir)
					fmt.Printf("Moving %v towards floor %d\n", dir.Dir, btn.Floor)
				}
			}

		case floor := <-drv_floors:
		elevator.Floor = floor
		elevio.SetFloorIndicator(floor)

		if elevator.ShouldStop() {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevator.Dir = elevio.MD_Stop
			elevator.Behaviour = elevio.DoorOpen
			elevio.SetDoorOpenLamp(true)

			// Clear completed orders
			for b := elevio.ButtonType(0); b < elevio.N_BUTTONS; b++ {
				elevator.Orders[floor][b] = 0
				elevio.SetButtonLamp(b, floor, false)
			}
		}

		// Edge case handling
		if floor == 0 {
			elevator.Orders[floor][elevio.BT_HallDown] = 0
		}
		if floor == numFloors-1 {
			elevator.Orders[floor][elevio.BT_HallUp] = 0
		}

		case obstruction := <-drv_obstr:
			if obstruction {
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(true)
			}

		case stop := <-drv_stop:
			fmt.Printf("%+v\n", stop)
			if stop {
				elevio.SetMotorDirection(elevio.MD_Stop)
				// Clear all button lamps
				for f := 0; f < numFloors; f++ {
					for b := elevio.ButtonType(0); b < 3; b++ {
						elevator.Orders[f][b] = 0
						elevio.SetButtonLamp(b, f, false)
					}
				}
			}
		}

		// Check if door timer has timed out
		if elevator.TimerTimedOut() {
			elevio.SetDoorOpenLamp(false)
			dir := elevator.ChooseDirection()
			elevator.Dir = dir.Dir
			elevator.Behaviour = dir.Behaviour
			elevio.SetMotorDirection(elevator.Dir)
		}
	}
}
