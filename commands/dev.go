//
package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	//
	DevCmd = &cobra.Command{
		Use:   "dev",
		Short: "Starts the nanobox, provisions app, & opens an interactive terminal",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			// processor.DefaultConfig.Meta["dev-config"] = devconfig
			processor.Run("dev", processor.DefaultConfig)
		},
		// PostRun: halt,
	}

	//
	devconfig string // sets the type of environment to be configured on the guest vm
	// nobuild   bool   // force skip a deploy
	rebuild bool // force a deploy
)

//
func init() {
	// DevCmd.Flags().StringVarP(&devconfig, "dev-config", "", nanofile.Viper().GetString("dev-config"), "The environment to configure on the guest vm")
	// DevCmd.Flags().BoolVarP(&nobuild, "no-build", "", false, "Force skip a rebuild")
	// DevCmd.Flags().BoolVarP(&rebuild, "rebuild", "", false, "Force a rebuild")

	// all gone for now.. will bring back the necessary ones.
	// // 'hidden' commands
	// DevCmd.AddCommand(buildCmd)
	// DevCmd.AddCommand(createCmd)
	// DevCmd.AddCommand(deployCmd)
	// DevCmd.AddCommand(initCmd)
	// DevCmd.AddCommand(logCmd)
	// DevCmd.AddCommand(resumeCmd)
	// DevCmd.AddCommand(sshCmd)
	// DevCmd.AddCommand(watchCmd)

	// // 'nanobox dev' commands
	// DevCmd.AddCommand(bootstrapCmd)
	DevCmd.AddCommand(DevDeployCmd)

	// DevCmd.AddCommand(reloadCmd)
	DevCmd.AddCommand(DevDestroyCmd)
	DevCmd.AddCommand(InfoCmd)
	DevCmd.AddCommand(ExecCmd)
	DevCmd.AddCommand(ConsoleCmd)
	DevCmd.AddCommand(DevEnvCmd)
	// DevCmd.AddCommand(updateImagesCmd)
}
