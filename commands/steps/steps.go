package steps

import "github.com/spf13/cobra"

var (
	stepList = map[string]step{}
)

type (
	CompleteCheckFunc func() bool
	CmdFunc           func(ccmd *cobra.Command, args []string)

	step struct {
		complete CompleteCheckFunc
		cmd      CmdFunc
	}
)
