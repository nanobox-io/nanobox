//
package vagrant

import (
	"github.com/nanobox-io/nanobox/config"
	"os"
	"os/exec"
)

// SSH is run manually (vs Run) because the output needs to be hooked up differntly
func SSH() error {

	//
	setContext(config.AppDir)

	cmd := exec.Command("vagrant", "ssh")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output
	return cmd.Run()
}
