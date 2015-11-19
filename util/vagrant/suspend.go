//
package vagrant

import (
	"fmt"
	"os/exec"

	"github.com/nanobox-io/nanobox-golang-stylish"
)

// Suspend runs a "vagrant suspend"
func Suspend() error {

	// suspend the vm
	fmt.Printf("\n%s", stylish.Bullet("Suspending nanobox..."))
	if err := runInContext(exec.Command("vagrant", "suspend")); err != nil {
		return err
	}
	fmt.Printf(stylish.Bullet("Exiting"))

	return nil
}
