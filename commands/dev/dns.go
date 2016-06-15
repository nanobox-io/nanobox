package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
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
		Run:   dnsAddFn,
	}

	// DNSRemoveCmd ...
	DNSRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "",
		Long:  ``,
		Run:   dnsRmFn,
	}
)

//
func init() {
	DNSCmd.AddCommand(DNSAddCmd)
	DNSCmd.AddCommand(DNSRemoveCmd)
}

// dnsAddFn will run the DNS processor for adding DNS entires to the "hosts"
// file
func dnsAddFn(ccmd *cobra.Command, args []string) {
	processor.DefaultConfig.Meta["name"] = args[0]
	print.OutputCommandErr(processor.Run("dev_dns_add", processor.DefaultConfig))
}

// dnsRmFn will run the DNS processor for removing DNS entries from the "hosts"
// file
func dnsRmFn(ccmd *cobra.Command, args []string) {
	processor.DefaultConfig.Meta["name"] = args[0]
	print.OutputCommandErr(processor.Run("dev_dns_remove", processor.DefaultConfig))
}
