// +build !windows

package display

import (
	"fmt"
)

// PrintRequiresPrivilege prints a message informing privilege escalation is required
func PrintRequiresPrivilege(reason string) {
	fmt.Println("Root privileges are required %s, your password may be requested...", reason)
}
