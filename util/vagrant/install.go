//
package vagrant

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
	"os"
	"os/exec"
)

// Install downloads the nanobox vagrant and adds it to the list of vagrant boxes
func Install() error {

	// only download and install nanobox/boot2docker if it doesn't exist
	if _, err := os.Stat(config.Home + "/.vagrant.d/boxes/nanobox-VAGRANTSLASH-boot2docker"); err != nil {
		fmt.Printf(stylish.Bullet("Installing nanobox/boot2docker..."))

		// add nanobox/boot2docker
		return run(exec.Command("vagrant", "box", "add", "--force", "nanobox/boot2docker"))
	}

	fmt.Printf(stylish.Bullet("nanobox/boot2docker already installed"))

	return nil
}
