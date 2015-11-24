//
package vagrant

import (
	"os/exec"
	"runtime"
)

// Reload runs a vagrant reload
func Reload() error {

	// gain sudo privilages; not handling error here because worst case scenario
	// this fails and just prompts for a password later
	if runtime.GOOS != "windows" {
		exec.Command("sudo", "ls").Run()
	}

	return runInContext(exec.Command("vagrant", "reload", "--provision"))
}
