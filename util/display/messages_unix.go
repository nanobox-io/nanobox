// +build !windows

package display

import (
	"fmt"
)

// prints a message informing the user that the terminal is invalid
func InvalidTerminal() {
	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()
	fmt.Println("Oops, this terminal doesn't fully support nanobox.")
	fmt.Println()
	fmt.Println("--------------------------------------------------")
	fmt.Println()
}
