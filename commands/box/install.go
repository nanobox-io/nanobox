//
package box

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "",
	Long:  ``,

	Run: Install,
}

// Install
func Install(ccmd *cobra.Command, args []string) {
	if err := checkInstall(); err != nil {
		Config.Fatal("[commands/box/install] checkInstall() failed - ", err.Error())
	}
}

// checkInstall
func checkInstall() (err error) {
	// install the nanobox vagrant image only if it isn't already available
	if !Vagrant.HaveImage() {
		fmt.Printf(stylish.Bullet("Installing nanobox image..."))

		// install the nanobox vagrant image
		err = Vagrant.Install()
	}
	return
}
