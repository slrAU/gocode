// Always start with 'go mod init [path/]projectname' at the command line
// to generate the project's MOD file.

package main

import (
	"machine"
	"time"
)

// Blink the Raspberry Pi Pico's on board LED on and off.
// Could apply to any pin output requiring a Low/High double toggle
func lowHigh(pin machine.Pin, rate int) {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		pin.Low()
		time.Sleep(time.Millisecond * time.Duration(rate))
		pin.High()
		time.Sleep(time.Millisecond * time.Duration(rate))
	}
}

func main() {
	// Trigger functions to to call concurrently via a go routine.
	// Ideal if they  can run independently unless needing to synchronise comms/data with other routines.
	go lowHigh(machine.LED, 500)
	go lowHigh(machine.GP15, 250)

	for {
		// Endless Main Loop required to ensure the program doesn't exit
		// on hardware that is intended to run endlessly, but must be given
		// something to do or it will not yeild to go routines.
		time.Sleep(time.Millisecond * 250)
	}
}
