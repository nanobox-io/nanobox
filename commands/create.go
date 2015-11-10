//
package commands

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/config"
	"github.com/spf13/cobra"
	"os"
)

var (

	//
	createCmd = &cobra.Command{
		Hidden: true,

		Use:   "create",
		Short: "Creates a new nanobox",
		Long:  ``,

		PreRun: initialize,
		Run:    create,
	}

	//
	addEntry bool // does an entry need to be added to the /etc/hosts file
)

//
func init() {
	createCmd.Flags().BoolVarP(&addEntry, "add-entry", "", false, "")
	createCmd.Flags().MarkHidden("add-entry")
}

//
// create
func create(ccmd *cobra.Command, args []string) {

	// PreRun: initialize

	// if the command is being run with the "add" flag, it means an entry needs to
	// be added to the hosts file and execution yielded back to the parent
	if addEntry {
		Hosts.AddDomain()
		os.Exit(0) // this exits the sudoed (child) created, not the parent proccess
	}

	// boot the vm
	fmt.Printf(stylish.Bullet("Creating a nanobox"))
	if err := Vagrant.Up(); err != nil {
		Config.Fatal("[commands/create] vagrant.Up() failed - ", err.Error())
	}

	// after the machine boots, update the docker images
	updateImages(nil, args)

	// add the entry if needed
	if !Hosts.HasDomain() {
		sudo("create --add-entry", fmt.Sprintf("Adding %v domain to hosts file", config.Nanofile.Domain))
	}

	// if devmode is detected, the machine needs to be rebooted for devmode to take
	// effect
	if config.Devmode {
		fmt.Printf(stylish.Bullet("Rebooting machine to finalize 'devmode' configuration..."))
		reload(nil, args)
	}
}
