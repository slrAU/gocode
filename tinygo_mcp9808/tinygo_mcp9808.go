package tinygo_mcp9808

import "machine"

/*
Return a temperature value from an MCP9808 temperature Sensor.

	i2c *machine.I2C
		- machine.I2C0 or machine.I2C1. Must already be configured
		  with Frequency, and SDA/SCL pins set.
	address byte
		- a value between 0..7 representing the i2c bus address
		  of the temperature sensor.
*/
func ReadTemperature(i2c *machine.I2C, address byte) (float32, error) {
	address = 0b00011000 | (address & 0x07) // from MCP9808 datasheet, first 0 not used, 0011=address code, 000=unjumpered address
	const REG_TEMP = 0b00000101             // from MCP9808 datasheet, first 0 not used, must zero bits 7-4, 0101=temp register
	const TEMP_MASK = 0b00001111
	const SIGN_MASK = 0b00010000
	var temperature float32
	tempdata := []byte{0, 0}

	err := i2c.ReadRegister(address, REG_TEMP, tempdata)

	if err == nil {
		temperature = float32((float32(tempdata[0]&TEMP_MASK)*16 + float32(tempdata[1])/16))

		if tempdata[0]&SIGN_MASK == SIGN_MASK {
			temperature = 256 - temperature
		}
	}

	return temperature, err
}
