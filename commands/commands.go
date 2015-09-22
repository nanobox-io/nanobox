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
	"time"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-cli/util"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var (

	//
	NanoboxCmd = &cobra.Command{
		Use:   "nanobox",
		Short: "",
		Long:  ``,

		//
		Run: func(ccmd *cobra.Command, args []string) {

			// hijack the verbose flag (-v), and use it to display the version of the
			// CLI
			if fVersion || fVerbose {
				fmt.Printf("nanobox %v\n", config.Version.String())
				os.Exit(0)
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}

	//
	engineCmd = &cobra.Command{
		Use:   "engine",
		Short: "",
		Long:  ``,

		//
		// Run: func(cmd *cobra.Command, args []string) {},
	}

	//
	imagesCmd = &cobra.Command{
		Use:   "images",
		Short: "",
		Long:  ``,

		//
		// Run: func(cmd *cobra.Command, args []string) {},
	}

	//
	productionCmd = &cobra.Command{
		Use:   "production",
		Short: "",
		Long:  ``,

		//
		// Run: func(cmd *cobra.Command, args []string) {},
	}

	// persistent (global) flags
	fDebug   bool //
	fDevmode bool //
	fForce   bool //
	fVerbose bool //

	// local flags
	fCount   int    //
	fFile    string //
	fLevel   string //
	fOffset  int    //
	fRemove  bool   //
	fRun     bool   //
	fSandbox bool   //
	fStream  bool   //
	fTunnel  string //
	fVersion bool   //
	fWatch   bool   //
	fWrite   bool   //
)

//
type Service struct {
	CreatedAt time.Time
	Name      string
	Password  string
	Ports     []int
	Username  string
}

// init creates the list of available nanobox commands and sub commands
func init() {

	// internal flags
	NanoboxCmd.PersistentFlags().BoolVarP(&fDebug, "debug", "", false, "")
	NanoboxCmd.PersistentFlags().MarkHidden("debug")

	NanoboxCmd.PersistentFlags().BoolVarP(&fDevmode, "dev", "", false, "")
	NanoboxCmd.PersistentFlags().MarkHidden("dev")

	// persistent flags
	NanoboxCmd.PersistentFlags().BoolVarP(&fForce, "force", "f", false, "Forces a command to run (effects very per command).")
	NanoboxCmd.PersistentFlags().BoolVarP(&fVerbose, "verbose", "v", false, "Increase command output from 'info' to 'debug'.")

	// local flags
	NanoboxCmd.Flags().BoolVarP(&fVersion, "version", "", false, "Display the current version of this CLI")

	//
	// NanoboxCmd.SetHelpCommand(helpCmd)
	// NanoboxCmd.SetHelpFunc(nanoHelp)
	// NanoboxCmd.SetHelpTemplate("")
	// NanoboxCmd.SetUsageFunc(usageCmd)
	// NanoboxCmd.SetUsageTemplate("")

	// all available nanobox commands

	// 'hidden' commands
	NanoboxCmd.AddCommand(createCmd)
	NanoboxCmd.AddCommand(initCmd)
	NanoboxCmd.AddCommand(sshCmd)

	// 'public' commands
	NanoboxCmd.AddCommand(bootstrapCmd)
	NanoboxCmd.AddCommand(buildCmd)
	NanoboxCmd.AddCommand(consoleCmd)
	NanoboxCmd.AddCommand(deployCmd)
	NanoboxCmd.AddCommand(destroyCmd)
	NanoboxCmd.AddCommand(execCmd)
	NanoboxCmd.AddCommand(haltCmd)
	NanoboxCmd.AddCommand(logCmd)
	NanoboxCmd.AddCommand(newCmd)
	NanoboxCmd.AddCommand(publishCmd)
	NanoboxCmd.AddCommand(reloadCmd)
	NanoboxCmd.AddCommand(resumeCmd)
	NanoboxCmd.AddCommand(statusCmd)
	NanoboxCmd.AddCommand(suspendCmd)
	NanoboxCmd.AddCommand(tunnelCmd)
	NanoboxCmd.AddCommand(upCmd)
	NanoboxCmd.AddCommand(updateCmd)
	NanoboxCmd.AddCommand(watchCmd)

	// 'engine' subcommand
	NanoboxCmd.AddCommand(engineCmd)
	engineCmd.AddCommand(engineFetchCmd)
	engineCmd.AddCommand(enginePublishCmd)

	// 'images' subcommand
	NanoboxCmd.AddCommand(imagesCmd)

	// 'production' subcommand
	NanoboxCmd.AddCommand(productionCmd)
	// productionCmd.AddCommand(deployCmd)
}

// PRERUN COMMANDS

// VMIsRunning
func VMIsRunning(ccmd *cobra.Command, args []string) {
	if util.GetVMStatus() != "running" {
		fmt.Printf(stylish.ErrBullet("Please start your virtual machine before running this command."))
		os.Exit(1)
	}
}

// ProjectIsCreated
func ProjectIsCreated(ccmd *cobra.Command, args []string) {
	if fi, _ := os.Stat(config.AppDir); fi == nil {
		fmt.Printf(stylish.ErrBullet("Please create your project before running this command."))
		os.Exit(1)
	}
}
