//
package dev

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/box"
	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util/server"
	"github.com/nanobox-io/nanobox/util/vagrant"
)

var initCmd = &cobra.Command{
	Hidden: true,

	Use:   "init",
	Short: "Creates a nanobox-flavored Vagrantfile",
	Long:  ``,

	PreRun: runnable,
	Run:    initialize,
}

// initialize
func initialize(ccmd *cobra.Command, args []string) {

	// PreRun: runnable

	// check to see if a box needs to be installed
	box.Install(nil, args)

	// creates a project folder at ~/.nanobox/apps/<name> where the Vagrantfile and
	// .vagrant dir will live for each app
	if err := os.MkdirAll(config.AppDir, 0755); err != nil {
		config.Fatal("[commands/init] os.Mkdir() failed", err.Error())
	}

	// set up a dedicated vagrant logger
	vagrant.NewLogger(config.AppDir + "/vagrant.log")

	// set up a dedicated server logger
	server.NewLogger(config.AppDir + "/server.log")

	// 'parse' the .vmfile (either creating one, or parsing it)
	config.VMfile = config.ParseVMfile()

	//
	// generate a Vagrantfile at ~/.nanobox/apps/<app-name>/Vagrantfile
	// only if one doesn't already exist (unless forced)
	if !config.Force {
		if _, err := os.Stat(config.AppDir + "/Vagrantfile"); err == nil {
			return
		}
	}

	vagrant.Init()
}
