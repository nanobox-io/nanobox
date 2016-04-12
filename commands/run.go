//
package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	//
	RunCmd = &cobra.Command{
		Use:   "run",
		Short: "start a nanobox application as if it is in production",
		Long:  ``,

		PreRun: validCheck("provider"),
		Run: func(ccmd *cobra.Command, args []string) {
			processor.Run("run", processor.DefaultConfig)
		},
		// PostRun: halt,
	}
)

//
func init() {
	// DevCmd.Flags().BoolVarP(&nobuild, "no-build", "", false, "Force skip a rebuild")

	// all gone for now.. will bring back the necessary ones.
	// // 'hidden' commands
	// DevCmd.AddCommand(buildCmd)
	// DevCmd.AddCommand(createCmd)
	// DevCmd.AddCommand(deployCmd)
	// DevCmd.AddCommand(execCmd)
	// DevCmd.AddCommand(initCmd)
	// DevCmd.AddCommand(logCmd)
	// DevCmd.AddCommand(resumeCmd)
	// DevCmd.AddCommand(sshCmd)
	// DevCmd.AddCommand(watchCmd)

	// // 'nanobox dev' commands
	// DevCmd.AddCommand(bootstrapCmd)
	// DevCmd.AddCommand(runCmd)
	// DevCmd.AddCommand(reloadCmd)
	// DevCmd.AddCommand(stopCmd)
	// DevCmd.AddCommand(destroyCmd)
	// DevCmd.AddCommand(infoCmd)
	// DevCmd.AddCommand(consoleCmd)
	// DevCmd.AddCommand(updateImagesCmd)
}
