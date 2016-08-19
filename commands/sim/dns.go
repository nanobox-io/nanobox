package sim

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app/dns"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// DNSCmd ...
	DNSCmd = &cobra.Command{
		Use:   "dns",
		Short: "Manages hostname mappings for your sim app.",
		Long:  ``,
	}

	// DNSListCmd ...
	DNSListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists the dns entries for this app",
		Long: `
Lists a hostnames maped to your sim app. The domain provided is added
to your local hosts file pointing the the IP of your sim app.
		`,
		Run: dnsListFn,
	}

	// DNSAddCmd ...
	DNSAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a hostname map to your sim app.",
		Long: `
Adds a hostname map to your sim app. The domain provided is added
to your local hosts file pointing the the IP of your sim app.
		`,
		Run: dnsAddFn,
	}

	// DNSRemoveCmd ...
	DNSRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "Removes a hostname map from your sim app.",
		Long: `
Removes a hostname map from your sim app. The domain must perfectly
match an DNS entry in your to your local hosts file.
		`,
		Run: dnsRmFn,
	}
	// DNSRemoveCmd ...
	DNSRemoveAllCmd = &cobra.Command{
		Use:    "rm-all",
		Short:  "Removes all hostname mappings associated with your sim app.",
		Long:   ``,
		Run:    dnsRmAllFn,
		Hidden: true,
	}
)

//
func init() {
	DNSCmd.AddCommand(DNSListCmd)
	DNSCmd.AddCommand(DNSAddCmd)
	DNSCmd.AddCommand(DNSRemoveCmd)
	DNSCmd.AddCommand(DNSRemoveAllCmd)
}

// dnsListFn will run the DNS processor for adding DNS entires to the "hosts"
// file
func dnsListFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	dnsList := dns.List{App: app}
	display.CommandErr(dnsList.Run())
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

nanobox sim dns add <name>
`, len(args))

		return
	}

	// set the meta arguments to be used in the processor and run the processor
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	dnsAdd := dns.Add{App: app, Name: args[0]}
	display.CommandErr(dnsAdd.Run())
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

ex: nanobox sim dns rm <name>

`, len(args))

		return
	}

	// set the meta arguments to be used in the processor and run the processor
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	dnsRemove := dns.Remove{App: app, Name: args[0]}
	display.CommandErr(dnsRemove.Run())
}

// dnsRmAllFn will run the DNS processor for removing DNS entries from the "hosts"
// file
func dnsRmAllFn(ccmd *cobra.Command, args []string) {

	// set the meta arguments to be used in the processor and run the processor
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	dnsRemoveAll := dns.RemoveAll{App: app}
	display.CommandErr(dnsRemoveAll.Run())
}
