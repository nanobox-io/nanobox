//
package vagrant

import "os/exec"

// Reload runs a vagrant reload
func Reload() error {
	return runInContext(exec.Command("vagrant", "reload", "--provision"))
}
