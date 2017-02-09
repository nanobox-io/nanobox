package app

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app/dns"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
)

func init() {
	dns.AppSetup = Setup
}

// Setup sets up the app on the provider and in the database
func Setup(envModel *models.Env, appModel *models.App, name string) error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if this app is already active
	if appModel.State == "active" {
		goto RESERVE
	}

	// generate the app data
	if err := appModel.Generate(envModel, name); err != nil {
		lumber.Error("app:Setup:models.App:Generate(): %s", err.Error())
		return util.ErrorAppend(err, "failed to generate app data")
	}

RESERVE:
	// reserve IPs
	if err := reserveIPs(appModel); err != nil {
		return util.ErrorAppend(err, "failed to reserve app IPs")
	}

	// set app state to active
	appModel.State = "active"
	if err := appModel.Save(); err != nil {
		lumber.Error("app:Setup:models:App:Save()")
		return util.ErrorAppend(err, "failed to persist app state")
	}

	return nil
}

// reserIPs reserves app-level ip addresses
func reserveIPs(appModel *models.App) error {

	if appModel.LocalIPs["env"] != "" {
		return nil
	}

	display.StartTask("Reserving IPs")
	defer display.StopTask()

	if appModel.LocalIPs["env"] == "" {
		// reserve a dev ip
		envIP, err := dhcp.ReserveLocal()
		if err != nil {
			display.ErrorTask()
			lumber.Error("app:reserveIPs:dhcp.ReserveLocal()")
			return util.ErrorAppend(err, "failed to reserve an env IP")
		}

		// now assign the IPs onto the app model
		appModel.LocalIPs["env"] = envIP.String()

	}

	if appModel.Name == "sim" {
		if appModel.LocalIPs["logvac"] == "" {
			// reserve a logvac ip
			logvacIP, err := dhcp.ReserveLocal()
			if err != nil {
				display.ErrorTask()
				lumber.Error("app:reserveIPs:dhcp.ReserveLocal()")
				return util.ErrorAppend(err, "failed to reserve a logvac IP")
			}
			appModel.LocalIPs["logvac"] = logvacIP.String()

		}

		if appModel.LocalIPs["mist"] == "" {
			// reserve a mist ip
			mistIP, err := dhcp.ReserveLocal()
			if err != nil {
				display.ErrorTask()
				lumber.Error("app:reserveIPs:dhcp.ReserveLocal(): %s", err.Error())
				return util.ErrorAppend(err, "failed to reserve a mist IP: %s", err.Error())
			}

			appModel.LocalIPs["mist"] = mistIP.String()

		}

	}

	// save the app
	if err := appModel.Save(); err != nil {
		display.ErrorTask()
		lumber.Error("app:reserveIPs:models:App:Save(): %s", err.Error())
		return util.ErrorAppend(err, "failed to persist IPs to the db: %s", err.Error())
	}

	return nil
}
