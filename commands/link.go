//
package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (
	alias string

	LinkCmd = &cobra.Command{
		Use:   "link",
		Short: "link a dev application with a production one",
		Long:  `can be linked to more then one app`,
	}

	//
	LinkAddCmd = &cobra.Command{
		Use:   "add",
		Short: "link a dev application with a production one",
		Long:  `can be linked to more then one app`,

		Run: func(ccmd *cobra.Command, args []string) {
			processor.DefaultConfig.Meta["name"] = app
			processor.DefaultConfig.Meta["alias"] = alias
			handleError(processor.Run("link_add", processor.DefaultConfig))
		},
	}

	LinkListCmd = &cobra.Command{
		Use:   "List",
		Short: "list",
		Long:  `list`,

		Run: func(ccmd *cobra.Command, args []string) {
			handleError(processor.Run("link_list", processor.DefaultConfig))
		},
	}

	LinkRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "remove",
		Long:  `remove`,

		Run: func(ccmd *cobra.Command, args []string) {
			processor.DefaultConfig.Meta["alias"] = alias
			handleError(processor.Run("link_remove", processor.DefaultConfig))
		},
	}
)

func init() {
	LinkAddCmd.Flags().StringVarP(&app, "app_name", "n", "", "app name")
	LinkCmd.PersistentFlags().StringVarP(&alias, "alias", "a", "", "alias")

	LinkCmd.AddCommand(LinkAddCmd)
	LinkCmd.AddCommand(LinkListCmd)
	LinkCmd.AddCommand(LinkRemoveCmd)

}
