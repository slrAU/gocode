package main

import (
	"fmt"

	"somemath"
	"somemath/special"
	"somepackage"
)

func main() {
	fmt.Println("Output from Main() function")
	somepackage.SomePackagePrint()
	fmt.Println("sub package maybe? ", special.GetOneHundred())
	fmt.Println("10 + 2 = ", somemath.Add(10, 2))
	fmt.Println("10 - 2 = ", somemath.Subtract(10, 2))
	fmt.Println("10 x 2 = ", somemath.Subtract(10, 2))
	fmt.Println("10 / 2 = ", somemath.Subtract(10, 2))
}
