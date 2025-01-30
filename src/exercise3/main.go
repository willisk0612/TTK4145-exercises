package main

import (
	"fmt"
	"main/dev-tools/driver-go/elevio"
)

func main() {
	elevio.Init("localhost:15657", elevio.N_FLOORS)
	elevator := InitElevator()

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
			elevator.HandleButtonPress(btn)

		case floor := <-drv_floors:
			elevator.HandleFloorArrival(floor)

		case obstruction := <-drv_obstr:
			elevator.HandleObstruction(obstruction)

		case <-drv_stop:
			elevator.HandleStop()
		default:
			elevator.HandleDoorTimeout()
		}
	}
}
