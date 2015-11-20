//
package vagrant

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/config"
)

// Update downloads the nanobox vagrant and adds it to the list of vagrant boxes
func Update() error {

	// ensure nanobox/boot2docker has already been installed
	if _, err := os.Stat(config.Home + "/.vagrant.d/boxes/nanobox-VAGRANTSLASH-boot2docker"); err != nil {
		fmt.Printf(stylish.ErrBullet("nanobox/boot2docker is not installed. Please run 'nanobox box install'"))
		return nil
	}

	fmt.Printf(stylish.Bullet("Update nanobox/boot2docker..."))
	return runInContext(exec.Command("vagrant", "box", "update"))
}
