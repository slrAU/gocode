// API for a YL-3 8x7-segment display panel.
//
// These devices often consist of two 4x7-segment panels and
// a pair of 74HC165 (or equivalent) shift registers with a
// 5-pin connector.
//
// These devices are also often common cathode.
// All segments that are zeroed will be lit.
// Segments are referenced MSB to LSB: decimal, center, top left, then clockwise to top segment at LSB
// Panels are referenced left to right: as bit values 1, 2, 4, 8, 16, 32, 64, 128
package tinygo_yl3_8x7

import (
	"machine"
	"time"
)

// 7-segment Character Constant
const (
	DIGIT_0       uint8 = 0b11000000
	DIGIT_1       uint8 = 0b11111001
	DIGIT_2       uint8 = 0b10100100
	DIGIT_3       uint8 = 0b10110000
	DIGIT_4       uint8 = 0b10011001
	DIGIT_5       uint8 = 0b10010010
	DIGIT_6       uint8 = 0b10000010
	DIGIT_7       uint8 = 0b11111000
	DIGIT_8       uint8 = 0b10000000
	DIGIT_9       uint8 = 0b10010000
	DIGIT_DECIMAL uint8 = 0b01111111
	DIGIT_OFF     uint8 = 0b11111111
)

// Return all 7-segment character constants as an array
func DIGITS() []uint8 {
	return []uint8{
		DIGIT_0,
		DIGIT_1,
		DIGIT_2,
		DIGIT_3,
		DIGIT_4,
		DIGIT_5,
		DIGIT_6,
		DIGIT_7,
		DIGIT_8,
		DIGIT_9,
		DIGIT_DECIMAL,
		DIGIT_OFF,
	}
}

var data machine.Pin  // DIO
var clock machine.Pin // SCK
var latch machine.Pin // RCK

// Represents a delay time to apply in between pulling a pin high
// then low again. This is a customisable value to allow tweaking
// of the display, allowing optimisation of brightness/flicker
// to tweak runtime performance.
var pulseDelay = uint16(150)

func initialiseOutputPin(pin machine.Pin, state bool) {
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
	time.Sleep(time.Microsecond * time.Duration(pulseDelay))
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

// Configure the YL-3 display prior to accepting calls to WriteDigit().
//
//		data uint8
//			GPIO pin index - sometimes labelled `DIO`
//		clock uint8
//			GPIO pin index - sometimes labelled `SCK`
//		latch uint8
//			GPIO pin index - sometimes labelled `RCK`
//		pulseMicroseconds uint16 (optional)
//	        The delay in microseconds to apply to when
//	        pulsing a pin high then low during comms.
func Initialise(dataPin uint8, clockPin uint8, latchPin uint8, pulseMicroseconds ...uint16) {
	data = machine.Pin(dataPin)
	initialiseOutputPin(data, false)
	clock = machine.Pin(clockPin)
	initialiseOutputPin(clock, false)
	latch = machine.Pin(latchPin)
	initialiseOutputPin(latch, false)

	if len(pulseMicroseconds) >= 1 {
		pulseDelay = pulseMicroseconds[0]
	}
}

// Write a character pattern to a single 7-segment section.
// This function simplifies how the display panels and led patterns
// are referenced.  Panels are indexed from left to right. Patterns should
// use the constants provided, but can represent any 8-bit value as required.
//
//	position uint8
//		panel position index from 0 (left-most) to 7 (right-most)
//	pattern uint8
//		8-bit segment pattern to write out (0xFF = off/blank)
func WriteDigit(position uint8, pattern uint8) {
	// YL-3 expects the 2 shift registers to be populated in turn.
	// First, the 7-segment panel to activate
	shiftOut(1 << position)
	// Next, the led pattern to activate
	shiftOut(pattern)
	// Finally, pulse the latch pin to display the pattern
	pulsePin(latch)
}
