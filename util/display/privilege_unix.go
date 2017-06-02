// +build !windows

package display

import (
	"fmt"
)

// PrintRequiresPrivilege prints a message informing privilege escalation is required
func PrintRequiresPrivilege(reason string) {
	fmt.Printf("Root privileges are required %s. Your system password may be requested...\n", reason)
}
