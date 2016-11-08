package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/dns"
)

var (

	// DnsCmd ...
	DnsCmd = &cobra.Command{
		Use:   "dns",
		Short: "Manages dns",
		Long:  ``,
	}
)

//
func init() {
	DnsCmd.AddCommand(dns.AddCmd)
	DnsCmd.AddCommand(dns.RemoveCmd)
	DnsCmd.AddCommand(dns.ListCmd)
}
