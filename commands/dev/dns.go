package dev

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
)

var (

	// DNSCmd ...
	DNSCmd = &cobra.Command{
		Use:   "dns",
		Short: "Manages hostname mappings for your dev app.",
		Long:  ``,
	}

	// DNSAddCmd ...
	DNSAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a hostname map to your dev app.",
		Long:  `
Adds a hostname map to your dev app. The domain provided is added
to your local hosts file pointing the the IP of your dev app.
		`,
		Run:   dnsAddFn,
	}

	// DNSRemoveCmd ...
	DNSRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "Removes a hostname map from your dev app.",
		Long:  `
Removes a hostname map from your dev app. The domain must perfectly
match an DNS entry in your to your local hosts file.
		`,
		Run:   dnsRmFn,
	}

	// DNSRemoveCmd ...
	DNSRemoveAllCmd = &cobra.Command{
		Use:   "rm-all",
		Short: "Removes all hostname mappings associated with your dev app.",
		Long:  ``,
		Run:   dnsRmAllFn,
		Hidden: true,
	}
)

//
func init() {
	DNSCmd.AddCommand(DNSAddCmd)
	DNSCmd.AddCommand(DNSRemoveCmd)
	DNSCmd.AddCommand(DNSRemoveAllCmd)
}

// dnsAddFn will run the DNS processor for adding DNS entires to the "hosts"
// file
func dnsAddFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will return with instructions
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
name of the DNS entry you would like to add:

nanobox dev dns add <name>
`, len(args))

		return
	}

	// set the meta arguments to be used in the processor and run the processor
	processor.DefaultControl.Meta["name"] = args[0]
	processor.DefaultControl.Env = "dev"
	print.OutputCommandErr(processor.Run("env_dns_add", processor.DefaultControl))
}

// dnsRmFn will run the DNS processor for removing a DNS from the "hosts"
// file
func dnsRmFn(ccmd *cobra.Command, args []string) {

	// validate we have args required to set the meta we'll need; if we don't have
	// the required args this will return with instructions
	if len(args) != 1 {
		fmt.Printf(`
Wrong number of arguments (expecting 1 got %v). Run the command again with the
name of the DNS entry you would like to remove:

ex: nanobox dev dns rm <name>

`, len(args))

		return
	}

	// set the meta arguments to be used in the processor and run the processor
	processor.DefaultControl.Meta["name"] = args[0]
	processor.DefaultControl.Env = "dev"
	print.OutputCommandErr(processor.Run("env_dns_remove", processor.DefaultControl))
}

// dnsRmAllFn will run the DNS processor for removing DNS entries from the "hosts"
// file
func dnsRmAllFn(ccmd *cobra.Command, args []string) {

	// set the meta arguments to be used in the processor and run the processor
	processor.DefaultControl.Env = "dev"
	print.OutputCommandErr(processor.Run("env_dns_remove_all", processor.DefaultControl))
}
