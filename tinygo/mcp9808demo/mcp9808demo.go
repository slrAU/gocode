package main

import (
	"machine"
	"math"
	temperature "slrAU/tinygo_mcp9808"
	display "slrAU/tinygo_yl3_8x7"
	"time"
)

// var temperature float32

func showDigits(value float32) {
	temp := int(math.Trunc(float64(value) * 100))
	digits := [8]uint8{}

	for position := 7; position >= 0; position-- {
		reduced := temp / 10
		digits[position] = uint8(temp - reduced*10)
		temp = reduced
	}

	for index, digit := range digits {
		display.WriteDigit(uint8(index), display.DIGITS()[digit])
	}
}

func main() {
	i2c := machine.I2C0
	i2c.Configure(
		machine.I2CConfig{
			SCL:       machine.Pin(machine.GP21),
			SDA:       machine.Pin(machine.GP20),
			Frequency: machine.TWI_FREQ_400KHZ,
		},
	)

	display.Initialise(9, 8, 7)

	// const ADDRESS = 0b00011000  // from MCP9808 datasheet, 0 not used, 0011=address code, 000=unjumpered address, 0=read
	// const REG_TEMP = 0b00000101 // from MCP9808 datasheet, 0 not used, must zero bits 7-4, 0101=requests temperature register
	// const TEMP_MASK = 0b00001111
	// const SIGN_MASK = 0b00010000
	// tempdata := []byte{0, 0}

	for {
		// 	err := i2c.ReadRegister(ADDRESS, REG_TEMP, tempdata)

		// 	if err != nil {
		// 		temperature = 911000000 // Unrecognised state
		// 		if err == machine.ErrInvalidI2CBaudrate {
		// 			temperature = 91100001
		// 		}
		// 		if err == machine.ErrInvalidTgtAddr {
		// 			temperature = 91100002
		// 		}
		// 		if err == machine.ErrI2CGeneric {
		// 			temperature = 91100003
		// 		}
		// 		if err == machine.ErrRP2040I2CDisable {
		// 			temperature = 91100004
		// 		}
		// 	} else {

		// 		temperature = float32((float32(tempdata[0]&TEMP_MASK)*16 + float32(tempdata[1])/16))

		// 		if tempdata[0]&SIGN_MASK == SIGN_MASK {
		// 			temperature = 256 - temperature
		// 		}
		// 	}

		// 	showDigits(temperature)

		temp, _ := temperature.ReadTemperature(i2c, 0)
		showDigits(temp)

		time.Sleep(time.Millisecond * 5)
	}
}
