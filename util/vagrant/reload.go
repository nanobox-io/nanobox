//
package vagrant

import "os/exec"

// Reload runs a vagrant reload
func Reload() error {
	// gain sudo privilages
	cmd := exec.Command("sudo", "ls")
	cmd.CombinedOutput()
	
	return runInContext(exec.Command("vagrant", "reload", "--provision"))
}
