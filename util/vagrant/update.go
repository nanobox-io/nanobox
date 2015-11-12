//
package vagrant

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"os/exec"
)

// Update downloads the nanobox vagrant and adds it to the list of vagrant boxes
func Update() error {
	fmt.Printf(stylish.Bullet("Update nanobox/boot2docker..."))

	// update nanobox/boot2docker
	return run(exec.Command("vagrant", "box", "update"))
}
