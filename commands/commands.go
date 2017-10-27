// Package commands defines the commands that nanobox can run
package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/commands/server"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/update"
)

var (
	// debug mode
	debugMode bool

	// display level debug
	displayDebugMode bool

	// display level trace
	displayTraceMode bool

	internalCommand bool
	showVersion     bool
	endpoint        string

	// NanoboxCmd ...
	NanoboxCmd = &cobra.Command{
		Use:   "nanobox",
		Short: "",
		Long:  ``,
		PersistentPreRun: func(ccmd *cobra.Command, args []string) {
			// report the command to nanobox
			processors.SubmitLog(strings.Replace(ccmd.CommandPath(), "nanobox ", "", 1))
			// mixpanel.Report(strings.Replace(ccmd.CommandPath(), "nanobox ", "", 1))

			registry.Set("debug", debugMode)

			// setup the display output
			if displayDebugMode {
				lumber.Level(lumber.DEBUG)
				display.Summary = false
				display.Level = "debug"
			}

			if displayTraceMode {
				lumber.Level(lumber.TRACE)
				display.Summary = false
				display.Level = "trace"
			}

			// alert the user if an update is needed
			update.Check()

			configModel, _ := models.LoadConfig()

			// TODO: look into global messaging
			if internalCommand {
				registry.Set("internal", internalCommand)
				// setup a file logger, this will be replaced in verbose mode.
				fileLogger, _ := lumber.NewAppendLogger(filepath.ToSlash(filepath.Join(config.GlobalDir(), "nanobox.log")))
				lumber.SetLogger(fileLogger)

			} else {
				// We should only allow admin in 3 cases
				// 1 cimode
				// 2 server is running
				// 3 configuring
				fullCmd := strings.Join(os.Args, " ")
				if util.IsPrivileged() &&
					!configModel.CIMode &&
					!strings.Contains(fullCmd, "set ci") &&
					!strings.Contains(ccmd.CommandPath(), "server") {
					// if it is not an internal command (starting the server requires privilages)
					// we wont run nanobox as privilage
					display.UnexpectedPrivilage()
					os.Exit(1)
				}
			}

			if endpoint != "" {
				registry.Set("endpoint", endpoint)
			}

			if configModel.CIMode {
				lumber.Level(lumber.INFO)
				display.Summary = false
				display.Level = "info"
			}
		},

		Run: func(ccmd *cobra.Command, args []string) {
			if displayDebugMode || showVersion {
				fmt.Println(models.VersionString())
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
	NanoboxCmd.PersistentFlags().StringVarP(&endpoint, "endpoint", "e", "", "production endpoint")
	NanoboxCmd.PersistentFlags().MarkHidden("endpoint")
	NanoboxCmd.PersistentFlags().BoolVarP(&internalCommand, "internal", "", false, "Skip pre-requisite checks")
	NanoboxCmd.PersistentFlags().MarkHidden("internal")
	NanoboxCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "", false, "In the event of a failure, drop into debug context")
	NanoboxCmd.PersistentFlags().BoolVarP(&displayDebugMode, "verbose", "v", false, "Increases display output and sets level to debug")
	NanoboxCmd.PersistentFlags().BoolVarP(&showVersion, "version", "", false, "Print version information and exit")
	NanoboxCmd.PersistentFlags().BoolVarP(&displayTraceMode, "trace", "t", false, "Increases display output and sets level to trace")

	// log specific flags
	LogCmd.Flags().BoolVarP(&logRaw, "raw", "r", false, "Print raw log timestamps instead")
	LogCmd.Flags().BoolVarP(&logFollow, "follow", "f", false, "Follow logs (live feed)")
	LogCmd.Flags().IntVarP(&logNumber, "number", "n", 0, "Number of historic logs to print")
	// todo:
	// LogCmd.Flags().StringVarP(&logStart, "start", "", "", "Timestamp of oldest historic log to print")
	// LogCmd.Flags().StringVarP(&logEnd, "end", "", "", "Timestamp of newest historic log to print")
	// LogCmd.Flags().StringVarP(&logLimit, "limit", "", "", "Time to limit amount of historic logs to print")

	// subcommands
	NanoboxCmd.AddCommand(ConfigureCmd)
	NanoboxCmd.AddCommand(RunCmd)
	NanoboxCmd.AddCommand(BuildCmd)
	NanoboxCmd.AddCommand(CompileCmd)
	NanoboxCmd.AddCommand(DeployCmd)
	NanoboxCmd.AddCommand(ConsoleCmd)
	NanoboxCmd.AddCommand(RemoteCmd)
	NanoboxCmd.AddCommand(StatusCmd)
	NanoboxCmd.AddCommand(LoginCmd)
	NanoboxCmd.AddCommand(LogoutCmd)
	NanoboxCmd.AddCommand(CleanCmd)
	NanoboxCmd.AddCommand(InfoCmd)
	NanoboxCmd.AddCommand(TunnelCmd)
	NanoboxCmd.AddCommand(ImplodeCmd)
	NanoboxCmd.AddCommand(DestroyCmd)
	NanoboxCmd.AddCommand(StartCmd)
	NanoboxCmd.AddCommand(StopCmd)
	NanoboxCmd.AddCommand(UpdateCmd)
	NanoboxCmd.AddCommand(EvarCmd)
	NanoboxCmd.AddCommand(DnsCmd)
	NanoboxCmd.AddCommand(LogCmd)
	NanoboxCmd.AddCommand(VersionCmd)
	NanoboxCmd.AddCommand(server.ServerCmd)

	// hidden subcommands
	NanoboxCmd.AddCommand(EnvCmd)
	NanoboxCmd.AddCommand(InspectCmd)
}
