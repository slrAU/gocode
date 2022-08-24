package main

import (
	"machine"
	"math/rand"
	"time"
)

// These devices are often common cathode.
// All segments that are zeroed will be lit.
// MSB to LSB = decimal, center, top left, then clockwise to top segment at LSB
var DIGITS = [12]uint8{
	0b11000000, // 0
	0b11111001, // 1
	0b10100100, // 2
	0b10110000, // 3
	0b10011001, // 4
	0b10010010, // 5
	0b10000010, // 6
	0b11111000, // 7
	0b10000000, // 8
	0b10010000, // 9
	0b01111111, // decimal
	0b11111111, // off
}

var data = machine.Pin(machine.GP9)  // DIO Pin
var clock = machine.Pin(machine.GP8) // SCK Pin
var latch = machine.Pin(machine.GP7) // RCK Pin

func sleepMS(milliseconds uint) {
	time.Sleep(time.Millisecond * time.Duration(milliseconds))
}

func sleepUS(microseconds uint) {
	time.Sleep(time.Microsecond * time.Duration(microseconds))
}

func initOutputPin(pin machine.Pin, state bool) {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	if state {
		pin.High()
	} else {
		pin.Low()
	}
}

// Toggle a GPIO pin high then low.
//
// This is used both to shift bits on to the display's
// shift registers and to trigger the latch to complete
// the display process.
func pulsePin(pin machine.Pin) {
	pin.High()
	sleepUS(150)
	pin.Low()
}

// Push a single byte a bit at a time to the shift registers
// on the YL-3.
func shiftOut(value uint8) {
	// write MSBs first
	for bit := 0; bit < 8; bit++ {
		if value&0x80 == 0x80 {
			data.High()
		} else {
			data.Low()
		}
		pulsePin(clock)
		value = value << 0x01
	}
}

// Write a character pattern to a single 7-segment section
//
//	position: index 0-7 representing sections from left to right
//	pattern:  8-bit segment pattern to write out
//
// The YL-3 has a pair of shift registers that combine to represent
// a single digit/character on the display.  It expects up to 16 bits
// to be shifted out before triggering the latch to display the data.
// the 16-bit data is split into 2 bytes, with the first representing
// the character position to update, and the next representing the
// character leds to activate.
func writeDigit(position uint8, pattern uint8) {
	// YL-3 expects the 2 shift registers to be populated in turn.
	// First, the 7-segment panel to activate
	shiftOut(1 << position)
	// Next, the led pattern to activate
	shiftOut(pattern)
	// Finally, pulse the latch pin to display the pattern
	pulsePin(latch)
}

// A frivolous demo to show how to display information on
// the YL-3 display.  It is a comination of random numbers and
// a simple 'led chaser' pattern to highlight how the device
// functions.
func displayDemo() {
	const circleDelay = time.Millisecond * time.Duration(250)
	const randomDelay = time.Second

	values := [6]uint{}
	popValues := func() {
		for i := range values {
			values[i] = uint(rand.Intn(10))
		}
	}

	circleSegment := 1 << 5
	startCircle := time.Now()
	startRandom := startCircle

	popValues()

	for {
		// display random digits in the first 6 panels
		for position := uint8(0); position < 6; position++ {
			writeDigit(position, DIGITS[values[position]])
		}

		// I like how in hex, a blank pattern reads as 'OFF' :)
		writeDigit(6, 0xFF)

		// Animate a circular pattern on the right most segment
		writeDigit(7, ^uint8(circleSegment))

		if time.Since(startCircle) > circleDelay {
			if circleSegment > 1 {
				circleSegment = circleSegment >> 1
			} else {
				circleSegment = circleSegment << 5
			}

			startCircle = time.Now()
		}

		// randomise the digits every second
		if time.Since(startRandom) > randomDelay {
			popValues()
			startRandom = time.Now()
		}
	}
}

func main() {
	initOutputPin(latch, false)
	initOutputPin(clock, false)
	initOutputPin(data, false)

	// I prefer to use go routines to keep each component's operation
	// discrete, and minimise the work to maintain concurrent device
	// operations
	go displayDemo()

	for {
		// on the Pi Pico that I test with, this delay is required
		// otherwise TinyGo seems to forget to yeild to other functions
		sleepMS(5)
	}
}
