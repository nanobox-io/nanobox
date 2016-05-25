package commands

import(
  "github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (
  DevResetCmd = &cobra.Command{
    Use: "reset",
    Short: "reset the dev usage counters",
    Long: ``,

    PreRun: validCheck("provider"),
    Run: func(ccmd *cobra.Command, args []string) {
      // TODO: Take an extra arguement and decide what we want to reset
      processor.Run("dev_reset", processor.DefaultConfig)
    },
  }
)
