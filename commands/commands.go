// Package commands ...
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-golang-stylish"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	_ "github.com/nanobox-io/nanobox/processor/app"
	_ "github.com/nanobox-io/nanobox/processor/code"
	_ "github.com/nanobox-io/nanobox/processor/dev"
	_ "github.com/nanobox-io/nanobox/processor/link"
	_ "github.com/nanobox-io/nanobox/processor/platform"
	_ "github.com/nanobox-io/nanobox/processor/provider"
	_ "github.com/nanobox-io/nanobox/processor/service"
	_ "github.com/nanobox-io/nanobox/processor/env"
	_ "github.com/nanobox-io/nanobox/processor/env/dns"
	_ "github.com/nanobox-io/nanobox/processor/env/netfs"
	_ "github.com/nanobox-io/nanobox/processor/sim"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
)

//
var (

	//
	version bool // display the version of nanobox

	// NanoboxCmd ...
	NanoboxCmd = &cobra.Command{
		Use:   "nanobox",
		Short: "",
		Long:  ``,

		// if the verbose flag is used, up the log level (this will cascade into
		// all subcommands since this is the root command)
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {

			// NOTE: each time nanobox is run we want to check to see if there are
			// updates available. If we need to update, UpdateCmd is run. After an
			// update nanobox will os.Exit(0) (see processor/update.go)

			// stat the ~.nanobox/.update file to get its last modified time; we don't
			// need to handle the error here becase config.UpdateFile() will either return
			// safely if the file exists, or create it if it doesn't.
			fi, _ := os.Stat(config.UpdateFile())

			// get the time since ./nanbox/.update was last updated
			lastUpdated := time.Since(fi.ModTime())

			//
			switch {

			// if lastUpdated is less than [<= 1 second] ago then we'll assume that
			// the file was just created and we'll prompt for an update; this case is
			// for people who probably used the installer and most likely have an old
			// version of nanobox, or are using nanobox for the first time
			case lastUpdated.Seconds() <= 1:
				fmt.Printf(stylish.Bullet("First time running nanobox - checking for updates..."))
				processor.DefaultControl.Force = true
				UpdateCmd.Run(nil, nil)

			// if lastUpdated is more than [14 days] ago, then we'll run the auto-update
			// process, prompting the user if they want to update
			case lastUpdated.Hours()/24 >= 14.0:
				fmt.Printf(stylish.Bullet("14 days since last update - checking for updates ..."))
				processor.DefaultControl.Force = true
				UpdateCmd.Run(nil, nil)
			}

			// set verbose
			if processor.DefaultControl.Verbose {
				// close the existing loggers
				lumber.Close()
				// create a new multi logger
				multiLogger := lumber.NewMultiLogger()

				fileLogger, err := lumber.NewTruncateLogger(filepath.ToSlash(filepath.Join(config.GlobalDir(), "nanobox.log")))
				if err != nil {
					fmt.Println("logging error:", err)
				}

				// now logs go to the console as well as a file
				multiLogger.AddLoggers(fileLogger, lumber.NewConsoleLogger(lumber.DEBUG))
				lumber.SetLogger(multiLogger)
				lumber.Level(lumber.DEBUG)
			}
		},

		//
		Run: func(ccmd *cobra.Command, args []string) {

			// hijack the verbose flag (-v), and use it to display the version of the
			// CLI
			if version || processor.DefaultControl.Verbose {
				fmt.Printf("nanobox v%v\n", util.VERSION)
				return
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}
)

// init creates the list of available nanobox commands and sub commands
func init() {

	// persistent flags
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultControl.Debug, "debug", "", false, "run nanobox in debug mode")
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultControl.Force, "force", "f", false, "Forces a command to run (effects vary per command).")
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultControl.Verbose, "verbose", "v", false, "Increases display output.")
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultControl.Quiet, "quiet", "q", false, "Decreases display output.")

	// local flags
	NanoboxCmd.Flags().BoolVarP(&version, "version", "", false, "Displays the current version of this CLI.")

	// nanobox commands
	NanoboxCmd.AddCommand(UpdateCmd)

	// subcommands
	NanoboxCmd.AddCommand(InspectCmd)
	NanoboxCmd.AddCommand(DeployCmd)
	NanoboxCmd.AddCommand(ConsoleCmd)
	NanoboxCmd.AddCommand(LinkCmd)
	NanoboxCmd.AddCommand(LoginCmd)
	NanoboxCmd.AddCommand(LogoutCmd)
	NanoboxCmd.AddCommand(BuildCmd)
	NanoboxCmd.AddCommand(DevCmd)
	NanoboxCmd.AddCommand(SimCmd)
	NanoboxCmd.AddCommand(EnvCmd)
	NanoboxCmd.AddCommand(DestroyCmd)
	NanoboxCmd.AddCommand(StartCmd)
	NanoboxCmd.AddCommand(StopCmd)
	
}
