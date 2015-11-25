//
package vagrant

import (
	"os/exec"

	"github.com/nanobox-io/nanobox/config"
)

// Reload runs a vagrant reload
func Reload() error {

	// gain sudo privilages; not handling error here because worst case scenario
	// this fails and just prompts for a password later
	if config.OS != "windows" {
		exec.Command("sudo", "ls").Run()
	}

	return runInContext(exec.Command("vagrant", "reload", "--provision"))
}
