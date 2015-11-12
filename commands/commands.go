// package commands ...
package commands

import (
	"fmt"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/nanobox-io/nanobox/commands/box"
	"github.com/nanobox-io/nanobox/commands/engine"
	"github.com/nanobox-io/nanobox/commands/production"
	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

//
var (

	//
	NanoboxCmd = &cobra.Command{
		Use:   "nanobox",
		Short: "",
		Long:  ``,

		//
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {

			// ensures all dependencies are satisfied before running nanobox commands
			runnable(ccmd, args)

			// if the verbose flag is used, up the log level (this will cascade into
			// all subcommands since this is the root command)
			if config.Verbose {
				config.LogLevel = "debug"
			}
		},

		//
		Run: func(ccmd *cobra.Command, args []string) {

			// hijack the verbose flag (-v), and use it to display the version of the
			// CLI
			if version || config.Verbose {
				fmt.Printf("nanobox v%s\n", config.VERSION)
				os.Exit(0)
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}

	//
	version bool // display the version of nanobox
)

// init creates the list of available nanobox commands and sub commands
func init() {

	// internal flags
	NanoboxCmd.PersistentFlags().BoolVarP(&config.Devmode, "dev", "", false, "")
	NanoboxCmd.PersistentFlags().MarkHidden("dev")

	// persistent flags
	NanoboxCmd.PersistentFlags().BoolVarP(&config.Background, "background", "", false, "Stops nanobox from auto-suspending.")
	NanoboxCmd.PersistentFlags().BoolVarP(&config.Force, "force", "f", false, "Forces a command to run (effects vary per command).")
	NanoboxCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "Increase command output from 'info' to 'debug'.")

	// local flags
	NanoboxCmd.Flags().BoolVarP(&version, "version", "", false, "Display the current version of this CLI")

	// 'hidden' commands
	NanoboxCmd.AddCommand(buildCmd)
	NanoboxCmd.AddCommand(createCmd)
	NanoboxCmd.AddCommand(deployCmd)
	NanoboxCmd.AddCommand(execCmd)
	NanoboxCmd.AddCommand(initCmd)
	NanoboxCmd.AddCommand(logCmd)
	NanoboxCmd.AddCommand(reloadCmd)
	NanoboxCmd.AddCommand(resumeCmd)
	NanoboxCmd.AddCommand(sshCmd)
	NanoboxCmd.AddCommand(watchCmd)

	// 'nanobox' commands
	NanoboxCmd.AddCommand(runCmd)
	NanoboxCmd.AddCommand(devCmd)
	NanoboxCmd.AddCommand(bootstrapCmd)
	NanoboxCmd.AddCommand(infoCmd)
	NanoboxCmd.AddCommand(consoleCmd)
	NanoboxCmd.AddCommand(destroyCmd)
	NanoboxCmd.AddCommand(publishCmd)
	NanoboxCmd.AddCommand(stopCmd)
	NanoboxCmd.AddCommand(updateCmd)
	NanoboxCmd.AddCommand(updateImagesCmd)

	// subcommands
	NanoboxCmd.AddCommand(box.BoxCmd)
	NanoboxCmd.AddCommand(engine.EngineCmd)
	NanoboxCmd.AddCommand(production.ProductionCmd)
}

// runnable ensures all dependencies are satisfied before running box commands
func runnable(ccmd *cobra.Command, args []string) {

	// ensure vagrant exists
	if exists := Vagrant.Exists(); !exists {
		fmt.Println("Missing dependency 'Vagrant'. Please download and install it to continue (https://www.vagrantup.com/).")
		os.Exit(1)
	}

	// ensure virtualbox exists
	if exists := util.VboxExists(); !exists {
		fmt.Println("Missing dependency 'Virtualbox'. Please download and install it to continue (https://www.virtualbox.org/wiki/Downloads).")
		os.Exit(1)
	}
}

// boot
func boot(ccmd *cobra.Command, args []string) {

	// ensure vagrant exists before trying to run any nanobox command
	if exists := Vagrant.Exists(); !exists {
		fmt.Println("Missing dependency 'Vagrant'. Please download and install it to continue (https://www.vagrantup.com/).")
		os.Exit(1)
	}

	// ensure virtualbox exists before trying to run any nanobox command
	if exists := util.VboxExists(); !exists {
		fmt.Println("Missing dependency 'Virtualbox'. Please download and install it to continue (https://www.virtualbox.org/wiki/Downloads).")
		os.Exit(1)
	}

	// ensure a Vagrantfile is available before attempting to boot the VM
	initialize(nil, args)

	// get the status to know what needs to happen with the VM
	status := Vagrant.Status()

	switch status {

	// vm is running - do nothing
	case "running":
		fmt.Printf(stylish.Bullet("Nanobox is already running"))
		break

	// vm is suspended - resume it
	case "saved":
		resume(nil, args)

	// vm is not created - create it
	case "not created":
		create(nil, args)

	// vm is in some unknown state - reload it
	default:
		fmt.Printf(stylish.Bullet("Nanobox is in an unknown state (%s). Reloading...", status))
		reload(nil, args)
	}

	//
	Server.Lock()

	// if the background flag is passed, set the mode to "background"
	if config.Background {
		config.VMfile.BackgroundIs(true)
	}
}

// halt
func halt(ccmd *cobra.Command, args []string) {

	//
	Server.Unlock()

	//

	if err := Server.Suspend(); err != nil {
		Config.Fatal("[commands/halt] server.Suspend() failed - ", err.Error())
	}

	//
	if err := Vagrant.Suspend(); err != nil {
		Config.Fatal("[commands/halt] vagrant.Suspend() failed - ", err.Error())
	}

	//
	// os.Exit(0)
}

// sudo runs a command as sudo
func sudo(command, msg string) {
	fmt.Printf(stylish.Bullet(msg))

	//
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo %v %v", os.Args[0], command))

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run command
	if err := cmd.Run(); err != nil {
		Config.Fatal("[commands/halt]", err.Error())
	}
}
