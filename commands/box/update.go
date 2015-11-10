//
package box

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "",
	Long:  ``,

	Run: Update,
}

// Update
func Update(ccmd *cobra.Command, args []string) {

	if err := checkInstall(); err != nil {
		Config.Fatal("[commands/box/update] checkInstall() failed - ", err.Error())
	}

	//
	match, err := Util.MD5sMatch(Config.Root()+"/nanobox-boot2docker.md5", "https://s3.amazonaws.com/tools.nanobox.io/boxes/vagrant/nanobox-boot2docker.md5")
	if err != nil {
		Config.Fatal("[commands/box/update] Util.MD5sMatch() failed - ", err.Error())
	}

	// if the local md5 does not match the remote md5 download the newest nanobox
	// image
	if !match {
		fmt.Printf(stylish.Bullet("Updating nanobox image..."))

		// update the nanobox vagrant image
		if err := Vagrant.Update(); err != nil {
			Config.Fatal("[commands/box/update] Vagrant.Update() failed - ", err.Error())
		}
	}
}
