// Somepackage is just file to import into another go file
//
// So far it looks like you need the package to be in a separate folder with
// the same name as the the package statement declaration
package somepackage

import (
	"fmt"
	"otherpackage"
)

// Function to try and call within the same 'main' package
func SomePackagePrint() {
	fmt.Println("Printed from another source file in the somepackage... er... package.")

	otherpackage.PrintOther()
}
