// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox-cli/config"
	"github.com/nanobox-io/nanobox-cli/util"
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
		Run: func(ccmd *cobra.Command, args []string) {

			// hijack the verbose flag (-v), and use it to display the version of the
			// CLI
			if fVersion || fVerbose {
				fmt.Printf("nanobox %s\n", config.VERSION)
				os.Exit(0)
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}

	// subcommands
	boxCmd        = &cobra.Command{Use: "box", Short: "", Long: ``}
	engineCmd     = &cobra.Command{Use: "engine", Short: "", Long: ``}
	imagesCmd     = &cobra.Command{Use: "images", Short: "", Long: ``}
	productionCmd = &cobra.Command{Use: "production", Short: "", Long: ``}

	// persistent (global) flags
	fBackground bool //
	fDevmode    bool //
	fForce      bool //
	fVerbose    bool //

	// flags
	fAddEntry    bool   //
	fCount       int    //
	fFile        string //
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

//
type Service struct {
	CreatedAt time.Time
	Name      string
	Password  string
	Ports     []int
	Username  string
	UID       string
}

// init creates the list of available nanobox commands and sub commands
func init() {

	// internal flags
	NanoboxCmd.PersistentFlags().BoolVarP(&fBackground, "background", "", false, "")
	NanoboxCmd.PersistentFlags().MarkHidden("background")

	NanoboxCmd.PersistentFlags().BoolVarP(&fDevmode, "dev", "", false, "")
	NanoboxCmd.PersistentFlags().MarkHidden("dev")

	// persistent flags
	NanoboxCmd.PersistentFlags().BoolVarP(&fForce, "force", "f", false, "Forces a command to run (effects vary per command).")
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

	NanoboxCmd.AddCommand(buildCmd)
	NanoboxCmd.AddCommand(createCmd)
	NanoboxCmd.AddCommand(deployCmd)
	NanoboxCmd.AddCommand(initCmd)
	NanoboxCmd.AddCommand(logCmd)
	NanoboxCmd.AddCommand(reloadCmd)
	NanoboxCmd.AddCommand(resumeCmd)
	NanoboxCmd.AddCommand(sshCmd)
	NanoboxCmd.AddCommand(watchCmd)

	// 'public' commands
	NanoboxCmd.AddCommand(updateCmd)

	// 'nanobox' commands
	NanoboxCmd.AddCommand(nanoboxRunCmd)
	NanoboxCmd.AddCommand(bootstrapCmd)
	NanoboxCmd.AddCommand(nanoboxDevCmd)
	NanoboxCmd.AddCommand(nanoboxInfoCmd)
	NanoboxCmd.AddCommand(nanoboxConsoleCmd)
	NanoboxCmd.AddCommand(nanoboxExecCmd)
	NanoboxCmd.AddCommand(nanoboxDownCmd)
	NanoboxCmd.AddCommand(nanoboxDestroyCmd)
	NanoboxCmd.AddCommand(nanoboxPublishCmd)

	// 'box' subcommand
	NanoboxCmd.AddCommand(boxCmd)
	boxCmd.AddCommand(boxInstallCmd)
	boxCmd.AddCommand(boxUpdateCmd)

	// 'engine' subcommand
	NanoboxCmd.AddCommand(engineCmd)
	engineCmd.AddCommand(engineFetchCmd)
	engineCmd.AddCommand(engineNewCmd)
	engineCmd.AddCommand(enginePublishCmd)

	// 'images' subcommand
	NanoboxCmd.AddCommand(imagesCmd)
	imagesCmd.AddCommand(imagesUpdateCmd)

	// 'production' subcommand
	NanoboxCmd.AddCommand(productionCmd)
	// productionCmd.AddCommand(deployCmd)
}

// PRERUN COMMANDS

// vmIsRunning
func vmIsRunning(ccmd *cobra.Command, args []string) {
	if util.VagrantStatus() != "running" {
		fmt.Printf("Your nanobox VM is not running. Run 'nanobox up' first")
		os.Exit(1)
	}
}

// bootVM
func bootVM(ccmd *cobra.Command, args []string) {

	// check to see if a box needs to be installed
	boxInstall(nil, args)

	// ensure a Vagrantfile is available before attempting to boot the VM
	nanoInit(nil, args)

	// get the status to know what needs to happen with the VM
	status := util.VagrantStatus()
	switch status {

	// vm is running - do nothing
	case "running":
		fmt.Printf(stylish.Bullet("Nanobox VM already running..."))
		break

	// vm is suspended - resume it
	case "saved":
		nanoResume(nil, args)

	// vm is not created - create it
	case "not created":
		nanoCreate(nil, args)

	// vm is in some unknown state - reload it
	default:
		fmt.Printf(stylish.Bullet("Nanobox VM is in an unknown state."))
		nanoReload(nil, args)
	}

	// open a 'lock' with the server; this is done so we can know how many clients
	// are currently connected to the server
	// NOTE: the connection is NOT closed here. It is closed when saving the VM
	conn, err := net.Dial("tcp", config.ServerURI)
	if err != nil {
		config.Fatal("[commands/commands] new.Dial() failed", err.Error())
	}

	conn.Write([]byte(fmt.Sprintf("PUT /lock? HTTP/1.1\r\n\r\n")))

	//
	config.Lock = conn

	// after the VM is running updated the .vmfile
	config.VMfile.StatusIs(status)
	config.VMfile.UUIDIs(util.VagrantUUID())

	if fBackground {
		config.VMfile.ModeIs("background")
	}
}

// saveVM
func saveVM(ccmd *cobra.Command, args []string) {

	// close the connection to the server (indicating that a client is disconnecting)
	if config.Lock != nil {
		config.Lock.Close()
	}

	// if the CLI is running in background mode dont suspend the VM
	if config.VMfile.IsMode("background") {
		fmt.Printf(stylish.Bullet("Nanobox VM not suspended (running in background). To suspend the VM run 'nanobox down'"))
		os.Exit(0)
	}

	// check to see if the VM is able to be suspended
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/suspend", config.ServerURL), nil)
	if err != nil {
		config.Fatal("[commands/commands] http.NewRequest() failed", err.Error())
	}

	//
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		config.Fatal("[commands/commands] http.DefaultClient.Do() failed", err.Error())
	}
	defer res.Body.Close()

	//
	switch res.StatusCode / 100 {

	// anything but 200 CANNOT suspend
	default:
		config.VMfile.SuspendableIs(false)

	// ok to suspend
	case 2:
		config.VMfile.SuspendableIs(true)
		break
	}

	// suspend the machine if not active consoles are connected and the command was
	// not run in background mode
	switch {
	case !config.VMfile.IsSuspendable():
		fmt.Printf("\n   Note: The nanobox VM has not been suspended because there are other active console sessions.\n\n")
		break
	case config.VMfile.IsMode("background"):
		fmt.Printf(stylish.Bullet("\n   Note: the nanobox VM not suspended (running in background). To suspend the VM run 'nanobox down'\n\n"))
		break
	default:
		nanoboxDown(nil, args)
	}
}
