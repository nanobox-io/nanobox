//
package vagrant

import "os/exec"

// Resume runs a vagrant resume
func Resume() error {
	return runInContext(exec.Command("vagrant", "resume"))
}
