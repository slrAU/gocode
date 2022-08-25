package main

import (
	"math/rand"
	"time"

	"slrAU/tinygo_yl3_8x7"
)

func sleepMS(milliseconds uint) {
	time.Sleep(time.Millisecond * time.Duration(milliseconds))
}

// A frivolous demo to show how to display information on
// the YL-3 display.  It is a comination of random numbers and
// a simple 'led chaser' pattern to highlight how the device
// functions.
func displayDemo() {
	const circleDelay = time.Millisecond * time.Duration(250)
	const randomDelay = time.Second

	values := [6]uint{}
	setValues := func() {
		for i := range values {
			values[i] = uint(rand.Intn(10))
		}
	}

	circleSegment := 1 << 5
	startCircle := time.Now()
	startRandom := startCircle

	setValues()

	for {
		// display random digits in the first 6 panels
		for position := uint8(0); position < 6; position++ {
			tinygo_yl3_8x7.WriteDigit(position, tinygo_yl3_8x7.DIGITS()[values[position]])
		}

		// I like how in hex, a blank pattern reads as 'OFF' :)
		tinygo_yl3_8x7.WriteDigit(6, 0xFF)

		// Animate a circular pattern on the right most segment
		tinygo_yl3_8x7.WriteDigit(7, ^uint8(circleSegment))

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
			setValues()
			startRandom = time.Now()
		}
	}
}

func main() {
	tinygo_yl3_8x7.Initialise(9, 8, 7)

	// Go routines are ideal for concurrent processing on microcontrollers!!
	go displayDemo()

	for {
		// on the Pi Pico that I test with, this delay is required
		// otherwise TinyGo seems to forget to yeild to other functions
		sleepMS(5)
	}
}
