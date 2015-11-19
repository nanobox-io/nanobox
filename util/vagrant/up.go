//
package vagrant

import "os/exec"

// Up runs a vagrant up
func Up() error {
	// gain sudo privilages
	cmd := exec.Command("sudo", "ls")
	cmd.CombinedOutput()
	
	return runInContext(exec.Command("vagrant", "up"))
}
