//
package vagrant

import "os/exec"

// Destroy runs a vagrant destroy
func Destroy() error {
	return runInContext(exec.Command("vagrant", "destroy", "--force"))
}
