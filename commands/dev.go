//
package commands

import (
  "github.com/spf13/cobra"

  "github.com/nanobox-io/nanobox/processor"
  "github.com/nanobox-io/nanobox/util/nanofile"
)

var (

  //
  DevCmd = &cobra.Command{
    Use:   "dev",
    Short: "Starts the nanobox, provisions app, & opens an interactive terminal",
    Long:  ``,

    PreRun:  validCheck("vagrant", "virtualbox"),
    Run:     func(ccmd *cobra.Command, args []string) {
      processor.DefaultConfig.Meta["dev-config"] = devconfig
      if rebuild {
        processor.DefaultConfig.Meta["build-option"] = "rebuild"
      }
      processor.Run("dev", processor.DefaultConfig)
    },
    // PostRun: halt,
  }

  //
  devconfig string // sets the type of environment to be configured on the guest vm
  // nobuild   bool   // force skip a deploy
  rebuild   bool   // force a deploy
)

//
func init() {
  DevCmd.Flags().StringVarP(&devconfig, "dev-config", "", nanofile.Viper().GetString("dev-config"), "The environment to configure on the guest vm")
  // DevCmd.Flags().BoolVarP(&nobuild, "no-build", "", false, "Force skip a rebuild")
  DevCmd.Flags().BoolVarP(&rebuild, "rebuild", "", false, "Force a rebuild")

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

