package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
)

var (

	// DNSCmd ...
	DNSCmd = &cobra.Command{
		Use:   "dns",
		Short: "",
		Long:  ``,
	}

	// DNSAddCmd ...
	DNSAddCmd = &cobra.Command{
		Use:   "add",
		Short: "",
		Long:  ``,

		// nanobox dev dns add <name>
		Run: func(ccmd *cobra.Command, args []string) {
			processor.DefaultConfig.Meta["name"] = args[0]
			processor.Run("dev_dns_add", processor.DefaultConfig)
		},
	}

	// // DNSListCmd ...
	// DNSListCmd = &cobra.Command{
	// 	Use:   "ls",
	// 	Short: "",
	// 	Long:  ``,
	//
	// 	Run: func(ccmd *cobra.Command, args []string) {
	// 		processor.Run("dev_dns_list", processor.DefaultConfig)
	// 	},
	// }

	// DNSRemoveCmd ...
	DNSRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "",
		Long:  ``,

		// nanobox dev dns remove <name>
		Run: func(ccmd *cobra.Command, args []string) {
			processor.DefaultConfig.Meta["name"] = args[0]
			processor.Run("dev_dns_remove", processor.DefaultConfig)
		},
	}
)

//
func init() {
	DNSCmd.AddCommand(DNSAddCmd)
	// DNSCmd.AddCommand(DNSListCmd)
	DNSCmd.AddCommand(DNSRemoveCmd)
}
