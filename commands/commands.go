// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/commands/box"
	"github.com/nanobox-io/nanobox-cli/commands/engine"
	"github.com/nanobox-io/nanobox-cli/commands/production"
	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util/server"
	"github.com/nanobox-io/nanobox-cli/util/vagrant"
	"github.com/nanobox-io/nanobox-golang-stylish"
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
			if fVersion || config.Verbose {
				fmt.Printf("nanobox v%s\n", config.VERSION)
				os.Exit(0)
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}

	// flags
	fAddEntry    bool   //
	fCount       int    //
	fLevel       string //
	fOffset      int    //
	fRebuild     bool   //
	fRemove      bool   //
	fRemoveEntry bool   //
	fRun         bool   //
	fStream      bool   //
	fVersion     bool   //
	fWatch       bool   //
	fWrite       bool   //
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
	NanoboxCmd.Flags().BoolVarP(&fVersion, "version", "", false, "Display the current version of this CLI")

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

// boot
func boot(ccmd *cobra.Command, args []string) {

	// ensure a Vagrantfile is available before attempting to boot the VM
	initialize(nil, args)

	// get the status to know what needs to happen with the VM
	status := vagrant.Status()

	fmt.Println("STATUS?", status)

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
	server.Lock()

	// if the background flag is passed, set the mode to "background"
	if config.Background {
		config.VMfile.ModeIs("background")
	}
}

// save
func save(ccmd *cobra.Command, args []string) {

	//
	server.Unlock()

	//
	if err := server.Suspend(); err != nil {
		config.Fatal("[commands/commands] failed - ", err.Error())
	}

	//
	if err := vagrant.Suspend(); err != nil {
		config.Fatal("[commands/nanoboxDown] failed - ", err.Error())
	}
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
		config.Fatal("[utils/exec]", err.Error())
	}
}
