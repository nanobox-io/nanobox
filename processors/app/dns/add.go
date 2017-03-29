package dns

import (
	// "fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/dns"
	"github.com/nanobox-io/nanobox/processors/server"
)

var AppSetup func(envModel *models.Env, appModel *models.App, name string) error

// Add adds a dns entry to the local hosts file
func Add(envModel *models.Env, appModel *models.App, name string) error {

	if err := AppSetup(envModel, appModel, appModel.Name); err != nil {
		return util.ErrorAppend(err, "failed to setup app")
	}

	// fetch the IP
	// env in dev is used in the dev container
	// env in sim is used for portal
	envIP := appModel.LocalIPs["env"]

	// generate the dns entry
	entry := dns.Entry(envIP, name, appModel.ID)

	// short-circuit if this entry already exists
	if dns.Exists(entry) {
		return nil
	}

	// make sure the server is running since it will do the dns addition
	if err := server.Setup(); err != nil {
		return util.ErrorAppend(err, "failed to setup server")
	}

	// add the entry
	if err := dns.Add(entry); err != nil {
		lumber.Error("dns:Add:dns.Add(%s): %s", entry, err.Error())
		return util.ErrorAppend(err, "unable to add dns entry")
	}

	display.Info("\n%s %s added\n", display.TaskComplete, name)

	return nil
}