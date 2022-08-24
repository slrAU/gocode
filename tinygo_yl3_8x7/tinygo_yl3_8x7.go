// Driver library for a YL-3 8x7-segment display panel.
//
// These devices often consist of two 4x7-segment panels and
// a pair of 74HC165 (or equivalent) shift registers with a
// 5-pin connector.
package tinygo_yl3_8x7

import (
	"machine"
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

var data machine.Pin  // DIO
var clock machine.Pin // SCK
var latch machine.Pin // RCK

var pulseDelay = uint16(150)

func sleepUS(microseconds uint16) {
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
	sleepUS(pulseDelay)
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

// Initialise assigns and configures the microcontroller pins used
// to drive the YL-3 display
//
// data uint8  - sometimes labelled `DIO`
// clock uint8 - sometimes labelled `SCK`
// latch uint8 - sometimes labelled `RCK`
func Initialise(dataPin uint8, clockPin uint8, latchPin uint8, pulseMicroseconds ...uint16) {
	data = machine.Pin(dataPin)
	clock = machine.Pin(clockPin)
	latch = machine.Pin(latchPin)
	initOutputPin(latch, false)
	initOutputPin(clock, false)
	initOutputPin(data, false)

	if len(pulseMicroseconds) >= 1 {
		pulseDelay = pulseMicroseconds[0]
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
func WriteDigit(position uint8, pattern uint8) {
	// YL-3 expects the 2 shift registers to be populated in turn.
	// First, the 7-segment panel to activate
	shiftOut(1 << position)
	// Next, the led pattern to activate
	shiftOut(pattern)
	// Finally, pulse the latch pin to display the pattern
	pulsePin(latch)
}
