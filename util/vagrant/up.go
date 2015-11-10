//
package vagrant

import "os/exec"

// Up runs a vagrant up
func Up() error {
	return runInContext(exec.Command("vagrant", "up"))
}
