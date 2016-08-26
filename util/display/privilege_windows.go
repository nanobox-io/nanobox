// +build windows

package display

import (
	"fmt"
)

// PrintRequiresPrivilege prints a message informing privilege escalation is required
func PrintRequiresPrivilege(reason string) {
	fmt.Println("Administrator privileges are required %s.", reason)
	fmt.Println("Another window will be opened as the Administrator and permission may be requested.")

	// block here until the user hits enter. It's not ideal, but we need to make
	// sure they see the new window open if it requests their password.
	fmt.Println("Enter to continue:")
	var input string
	fmt.Scanln(&input)
}
