//
package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	//
	LinkCmd = &cobra.Command{
		Use:   "link",
		Short: "link a dev application with a production one",
		Long:  `can be linked to more then one app`,

		Run: func(ccmd *cobra.Command, args []string) {
			processor.DefaultConfig.Meta["name"] = appName
			processor.DefaultConfig.Meta["alias"] = alias
			processor.Run("link", processor.DefaultConfig)
		},
	}
)

func init() {
	LoginCmd.Flags().StringVarP(&app, "app_name", "n", "", "app name")
	LoginCmd.Flags().StringVarP(&alias, "alias", "a", "", "alias")
}
