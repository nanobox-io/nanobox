// Package commands ...
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	_ "github.com/nanobox-io/nanobox/processor/code"
	_ "github.com/nanobox-io/nanobox/processor/platform"
	_ "github.com/nanobox-io/nanobox/processor/provider"
	_ "github.com/nanobox-io/nanobox/processor/service"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/validate"
)

// VERSION ...
const VERSION = "1.0.0"

//
var (
	//
	NanoboxCmd = &cobra.Command{
		Use:   "nanobox",
		Short: "",
		Long:  ``,

		// if the verbose flag is used, up the log level (this will cascade into
		// all subcommands since this is the root command)
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {

			// each time nanobox is run we want to check to see if there are updates
			// available (if we haven't checked for [14 days]). If we need to update
			// UpdateCmd is run with the Force option set to true, this will trigger
			// the update processor to run the automatic update. After the update nanobox
			// will os.Exit(0) (see update processor)
			if fi, _ := os.Stat(util.UpdateFile()); time.Since(fi.ModTime()).Hours() >= 336.0 {
				processor.DefaultConfig.Force = true
				UpdateCmd.Run(nil, nil)
			}

			// set verbose
			if processor.DefaultConfig.Verbose {
				// close the existing loggers
				lumber.Close()
				// create a new multi logger
				multiLogger := lumber.NewMultiLogger()

				fileLogger, err := lumber.NewTruncateLogger(filepath.ToSlash(filepath.Join(util.GlobalDir(), "nanobox.log")))
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
			if version || processor.DefaultConfig.Verbose {
				fmt.Printf("nanobox v%v\n", VERSION)
				return
			}

			// fall back on default help if no args/flags are passed
			ccmd.HelpFunc()(ccmd, args)
		},
	}

	//
	version bool // display the version of nanobox
	app     string
)

// init creates the list of available nanobox commands and sub commands
func init() {

	// internal flags
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultConfig.DevMode, "dev", "", false, "")
	NanoboxCmd.PersistentFlags().MarkHidden("dev")

	// persistent flags
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultConfig.Background, "background", "", false, "Stops nanobox from auto-suspending.")
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultConfig.Force, "force", "f", false, "Forces a command to run (effects vary per command).")
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultConfig.Verbose, "verbose", "v", false, "Increase display output.")
	NanoboxCmd.PersistentFlags().BoolVarP(&processor.DefaultConfig.Quiet, "quiet", "q", false, "Decrease display output.")

	// local flags
	NanoboxCmd.Flags().BoolVarP(&version, "version", "", false, "Display the current version of this CLI")

	// nanobox commands
	NanoboxCmd.AddCommand(UpdateCmd)

	// subcommands
	NanoboxCmd.AddCommand(DeployCmd)
	NanoboxCmd.AddCommand(LinkCmd)
	NanoboxCmd.AddCommand(LoginCmd)
	NanoboxCmd.AddCommand(LogoutCmd)
	NanoboxCmd.AddCommand(BuildCmd)
	NanoboxCmd.AddCommand(DevCmd)
}

// validCheck ...
func validCheck(checks ...string) func(ccmd *cobra.Command, args []string) {
	return func(ccmd *cobra.Command, args []string) {
		if err := validate.Check(checks...); err != nil {
			fmt.Printf("Validation Failed:\n%s\n", err.Error())
			os.Exit(1)
		}
	}
}

// handleError ...
func handleError(err error) {
	if err != nil {
		fmt.Printf("It appears we have ran into some error:\n%s\n", err.Error())
	}
}
