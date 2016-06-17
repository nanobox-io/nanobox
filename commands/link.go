package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
)

var (

	// LinkCmd ...
	LinkCmd = &cobra.Command{
		Use:   "link",
		Short: "Manages links between dev & production apps.",
		Long:  ``,
	}

	// LinkAddCmd ...
	LinkAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a new link between a dev & production app.",
		Long: `
Adds a new link between a dev and production app. A dev app can
be linked to multiple production apps. Each link needs an alias.
If no alias is provided, 'default' is assumed.
		`,
		Run: linkAddFn,
	}

	// LinkListCmd ...
	LinkListCmd = &cobra.Command{
		Use:   "ls",
		Short: "Lists all links for the current dev app.",
		Long:  ``,
		Run:   linkListFn,
	}

	// LinkRemoveCmd ...
	LinkRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "Removes a link between a dev & production app.",
		Long:  ``,
		Run:   linkRemoveFn,
	}

	// linkCmdFlags ...
	linkCmdFlags = struct {
		app   string
		alias string
	}{}
)

//
func init() {
	LinkAddCmd.Flags().StringVarP(&linkCmdFlags.app, "app", "n", "", "app name")
	LinkCmd.PersistentFlags().StringVarP(&linkCmdFlags.alias, "alias", "a", "", "alias")

	LinkCmd.AddCommand(LinkAddCmd)
	LinkCmd.AddCommand(LinkListCmd)
	LinkCmd.AddCommand(LinkRemoveCmd)
}

// linkAddFn ...
func linkAddFn(ccmd *cobra.Command, args []string) {

	// set the meta arguments to be used in the processor and run the processor
	processor.DefaultConfig.Meta["app"] = linkCmdFlags.app
	processor.DefaultConfig.Meta["alias"] = linkCmdFlags.alias
	print.OutputCommandErr(processor.Run("link_add", processor.DefaultConfig))
}

// linkListFn ...
func linkListFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("link_list", processor.DefaultConfig))
}

// linkRemoveFn ...
func linkRemoveFn(ccmd *cobra.Command, args []string) {

	// set the meta arguments to be used in the processor and run the processor
	processor.DefaultConfig.Meta["alias"] = linkCmdFlags.alias
	print.OutputCommandErr(processor.Run("link_remove", processor.DefaultConfig))
}
