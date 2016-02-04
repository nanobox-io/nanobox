//
package dev

import (
	"fmt"
	"net/url"
	"os"

	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/config"
	"github.com/nanobox-io/nanobox/util"
	engineutil "github.com/nanobox-io/nanobox/util/engine"
	"github.com/nanobox-io/nanobox/util/server"
	mistutil "github.com/nanobox-io/nanobox/util/server/mist"
	"github.com/nanobox-io/nanobox/util/vagrant"
)

var (

	//
	DevCmd = &cobra.Command{
		Use:   "dev",
		Short: "Starts the nanobox, provisions app, & opens an interactive terminal",
		Long:  ``,

		PreRun:  boot,
		Run:     dev,
		PostRun: halt,
	}

	//
	devconfig string // sets the type of environment to be configured on the guest vm
	nobuild   bool   // force skip a deploy
	rebuild   bool   // force a deploy
)

//
func init() {
	DevCmd.Flags().StringVarP(&devconfig, "dev-config", "", config.Nanofile.DevConfig, "The environment to configure on the guest vm")
	DevCmd.Flags().BoolVarP(&nobuild, "no-build", "", false, "Force skip a rebuild")
	DevCmd.Flags().BoolVarP(&rebuild, "rebuild", "", false, "Force a rebuild")

	// 'hidden' commands
	DevCmd.AddCommand(buildCmd)
	DevCmd.AddCommand(createCmd)
	DevCmd.AddCommand(deployCmd)
	DevCmd.AddCommand(execCmd)
	DevCmd.AddCommand(initCmd)
	DevCmd.AddCommand(logCmd)
	DevCmd.AddCommand(resumeCmd)
	DevCmd.AddCommand(sshCmd)
	DevCmd.AddCommand(watchCmd)

	// 'nanobox dev' commands
	DevCmd.AddCommand(bootstrapCmd)
	DevCmd.AddCommand(runCmd)
	DevCmd.AddCommand(reloadCmd)
	DevCmd.AddCommand(stopCmd)
	DevCmd.AddCommand(destroyCmd)
	DevCmd.AddCommand(infoCmd)
	DevCmd.AddCommand(consoleCmd)
	DevCmd.AddCommand(updateImagesCmd)
}

// dev
func dev(ccmd *cobra.Command, args []string) {

	// PreRun: boot

	// check to see if the devconfig option is one of our predetermined values. If
	// not indicate as much and return
	switch devconfig {
	case "mount", "copy", "none":
		break
	default:
		fmt.Printf("--dev-config option '%s' is not accepted. Please choose either 'mount', 'copy', or 'none'\n", devconfig)
		os.Exit(1)
	}

	// stream log output; this is done here because there might be output from hooks
	// that needs to be streamed to the client. This will also capture any output
	// that would come from a deploy (if one is run).
	mist, err := mistutil.Connect([]string{"log", "deploy"}, mistutil.PrintLogStream)
	if err != nil {
		config.Fatal("[commands/dev] mistutil.Connect() failed", err.Error())
	}

	// don't rebuild
	if !nobuild {

		// if the vm has no been created or deployed, the rebuild flag, or the VM has
		// recently been reloaded do a deploy
		if vagrant.Status() == "not created" || !config.VMfile.HasDeployed() || rebuild || config.VMfile.HasReloaded() {

			fmt.Printf(stylish.Bullet("Deploying codebase..."))

			// remount the engine file at ~/.nanobox/apps/<app>/<engine> so any new scripts
			// will be used during the deploy
			if err := engineutil.RemountLocal(); err != nil {
				config.Debug("No engine mounted (not found locally).")
			}

			// run a deploy
			if err := server.Deploy(""); err != nil {
				server.Fatal("[commands/dev] server.Deploy() failed", err.Error())
			}

			// listen for status updates
			errch := make(chan error)
			go func() {
				errch <- mistutil.Listen([]string{"job", "deploy"}, mistutil.DeployUpdates)
			}()

			// wait for a status update (blocking)
			err := <-errch

			//
			if err != nil {
				fmt.Printf(err.Error())
				return
			}

			// reset "reloaded" to false after a successful deploy so as NOT to deploy
			// on subsequent runnings of "nanobox dev"
			config.VMfile.ReloadedIs(false)
		}
	}

	//
	v := url.Values{}
	v.Add("dev_config", devconfig)

	//
	if err := server.Develop(v.Encode(), mist); err != nil {
		server.Error("[commands/dev] server.Develop failed", err.Error())
	}

	// PostRun: halt
}

// runnable ensures all dependencies are satisfied before running dev commands
func runnable(ccmd *cobra.Command, args []string) {

	// ensure vagrant exists
	if exists := vagrant.Exists(); !exists {
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

	// ensures the cli can run before trying to boot vm
	runnable(nil, args)

	// ensure a Vagrantfile is available before attempting to boot the VM
	initialize(nil, args)

	// get the status to know what needs to happen with the VM
	status := vagrant.Status()

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
		config.VMfile.BackgroundIs(true)
	}
}

// halt
func halt(ccmd *cobra.Command, args []string) {

	//
	server.Unlock()

	//

	if err := server.Suspend(); err != nil {
		config.Fatal("[commands/halt] server.Suspend() failed", err.Error())
	}

	//
	if err := vagrant.Suspend(); err != nil {
		config.Fatal("[commands/halt] vagrant.Suspend() failed", err.Error())
	}

	//
	// os.Exit(0)
}
