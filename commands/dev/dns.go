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
		Run:   dnsAddFunc,
	}

	// DNSRemoveCmd ...
	DNSRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "",
		Long:  ``,
		Run:   dnsRmFunc,
	}
)

//
func init() {
	DNSCmd.AddCommand(DNSAddCmd)
	DNSCmd.AddCommand(DNSRemoveCmd)
}

// dnsAddFunc will run the DNS processor for adding DNS entires to the "hosts"
// file
func dnsAddFunc(ccmd *cobra.Command, args []string) {
	processor.DefaultConfig.Meta["name"] = args[0]
	processor.Run("dev_dns_add", processor.DefaultConfig)
}

// dnsRmFunc will run the DNS processor for removing DNS entries from the "hosts"
// file
func dnsRmFunc(ccmd *cobra.Command, args []string) {
	processor.DefaultConfig.Meta["name"] = args[0]
	processor.Run("dev_dns_remove", processor.DefaultConfig)
}
