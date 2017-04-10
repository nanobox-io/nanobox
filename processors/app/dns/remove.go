package dns

import (
	// "fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/server"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/dns"
)

// Remove removes a dns entry from the local hosts file
func Remove(a *models.App, name string) error {
	// fetch the IP
	// env in dev is used in the dev container
	// env in sim is used for portal
	envIP := a.LocalIPs["env"]

	// generate the dns entry
	entry := dns.Entry(envIP, name, a.ID)

	// short-circuit if this entry doesn't exist
	if !dns.Exists(entry) {
		return nil
	}

	// make sure the server is running since it will do the dns work
	if err := server.Setup(); err != nil {
		return util.ErrorAppend(err, "failed to setup server")
	}

	// remove the entry
	if err := dns.Remove(entry); err != nil {
		lumber.Error("dns:Remove:dns.Remove(%s): %s", entry, err.Error())
		return util.ErrorAppend(err, "unable to add dns entry: %s")
	}

	display.Info("\n%s %s removed\n", display.TaskComplete, name)

	return nil
}
