// Always start with 'go mod init [path/]projectname' at the command line
// to generate the project's MOD file.

package main

import (
	"machine"
	"time"
)

const BLINK_RATE_MIN = 100
const BLINK_RATE_MAX = 500
const DEBOUNCE_DELAY = 100
const IDLE_TIME = time.Millisecond * 250
const LED_GREEN = machine.LED
const LED_WHITE = machine.GP15
const MOMENTARY_1 = machine.GP16
const MOMENTARY_2 = machine.GP17

var ledActive = true
var blinkRate = 250

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
func lowHighSwitched(pin machine.Pin, rate int) {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	for {
		if ledActive {
			pin.Low()
			time.Sleep(time.Millisecond * time.Duration(blinkRate))
			pin.High()
			time.Sleep(time.Millisecond * time.Duration(blinkRate))
		} else {
			pin.Low()
			time.Sleep(time.Millisecond * time.Duration(blinkRate))
		}
	}
}

var debouncePin machine.Pin = machine.GP0
var debounceStart int64

func hasDebounced(p machine.Pin) bool {
	var result bool

	if debouncePin != p {
		debouncePin = p
		debounceStart = time.Now().UnixMilli()
		result = false
	} else {
		result = (time.Now().UnixMilli() - debounceStart) >= DEBOUNCE_DELAY
	}

	if result {
		debouncePin = machine.GP0
	}

	return result
}

func onChangeRate(p machine.Pin) {
	if hasDebounced(p) {
		if blinkRate <= BLINK_RATE_MIN {
			blinkRate = BLINK_RATE_MAX
		} else {
			blinkRate -= BLINK_RATE_MIN
		}
	}
}

func configureSwitches() {
	toggleBtn := machine.Pin(MOMENTARY_1)

	toggleBtn.Configure(machine.PinConfig{
		Mode: machine.PinInputPulldown,
	})

	// Example of a Pin interrupt using an anonymous function
	toggleBtn.SetInterrupt(
		machine.PinRising,
		func(p machine.Pin) {
			ledActive = !ledActive
		},
	)

	rateBtn := machine.Pin(MOMENTARY_2)
	rateBtn.Configure(machine.PinConfig{
		Mode: machine.PinInputPulldown,
	})

	// Example of a Pin interrupt using a callback function
	rateBtn.SetInterrupt(
		machine.PinFalling,
		onChangeRate,
	)
}

func main() {
	configureSwitches()

	// Trigger functions to to call concurrently via a go routine.
	// Ideal if they  can run independently unless needing to synchronise comms/data with other routines.
	go lowHigh(LED_GREEN, BLINK_RATE_MAX)
	go lowHighSwitched(LED_WHITE, blinkRate)

	for {
		// Endless Main Loop required to ensure the program doesn't exit
		// on hardware that is intended to run endlessly, but must be given
		// something to do or it will not yeild to go routines.
		time.Sleep(IDLE_TIME)
	}
}
