//
package commands

import (
	"github.com/nanobox-io/nanobox/commands/box"
	"github.com/nanobox-io/nanobox/config"
	"github.com/spf13/cobra"
	"os"
)

var initCmd = &cobra.Command{
	Hidden: true,

	Use:   "init",
	Short: "Creates a nanobox-flavored Vagrantfile",
	Long:  ``,

	Run: initialize,
}

// initialize
func initialize(ccmd *cobra.Command, args []string) {

	// check to see if a box needs to be installed
	box.Install(nil, args)

	// creates a project folder at ~/.nanobox/apps/<name> (if it doesn't already
	// exists) where the Vagrantfile and .vagrant dir will live for each app
	if _, err := os.Stat(config.AppDir); err != nil {
		if err := os.Mkdir(config.AppDir, 0755); err != nil {
			Config.Fatal("[commands/init] os.Mkdir() failed", err.Error())
		}
	}

	// set up a dedicated vagrant logger
	Vagrant.NewLogger(config.AppDir + "/vagrant.log")

	// set up a dedicated server logger
	Server.NewLogger(config.AppDir + "/server.log")

	// 'parse' the .vmfile (either creating one, or parsing it)
	config.VMfile = Config.ParseVMfile()

	//
	// generate a Vagrantfile at ~/.nanobox/apps/<app-name>/Vagrantfile
	// only if one doesn't already exist (unless forced)
	if !config.Force {
		if _, err := os.Stat(config.AppDir + "/Vagrantfile"); err == nil {
			return
		}
	}

	Vagrant.Init()
}
