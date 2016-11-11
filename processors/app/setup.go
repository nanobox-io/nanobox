package app

import (
	"fmt"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Setup sets up the app on the provider and in the database
func Setup(envModel *models.Env, appModel *models.App, name string) error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if this app is already active
	if appModel.State == "active" {
		return nil
	}

	// generate the app data
	if err := appModel.Generate(envModel, name); err != nil {
		lumber.Error("app:Setup:models.App:Generate(): %s", err.Error())
		return fmt.Errorf("failed to generate app data: %s", err.Error())
	}
	
	// reserve IPs
	if err := reserveIPs(appModel); err != nil {
		return fmt.Errorf("failed to reserve app IPs: %s", err.Error())
	}

	// set app state to active
	appModel.State = "active"
	if err := appModel.Save(); err != nil {
		lumber.Error("app:Setup:models:App:Save(): %s", err.Error())
		return fmt.Errorf("failed to persist app state: %s", err.Error())
	}

	return nil
}

// reserIPs reserves app-level ip addresses
func reserveIPs(appModel *models.App) error {
	display.StartTask("Reserving IPs")
	defer display.StopTask()

	// reserve a dev ip
	envIP, err := dhcp.ReserveGlobal()
	if err != nil {
		display.ErrorTask()
		lumber.Error("app:reserveIPs:dhcp.ReserveGlobal(): %s", err.Error())
		return fmt.Errorf("failed to reserve an env IP: %s", err.Error())
	}

	// now assign the IPs onto the app model
	appModel.GlobalIPs["env"] = envIP.String()

	if appModel.Name == "sim" {
		// reserve a logvac ip
		logvacIP, err := dhcp.ReserveLocal()
		if err != nil {
			display.ErrorTask()
			lumber.Error("app:reserveIPs:dhcp.ReserveLocal(): %s", err.Error())
			return fmt.Errorf("failed to reserve a logvac IP: %s", err.Error())
		}

		// reserve a mist ip
		mistIP, err := dhcp.ReserveLocal()
		if err != nil {
			display.ErrorTask()
			lumber.Error("app:reserveIPs:dhcp.ReserveLocal(): %s", err.Error())
			return fmt.Errorf("failed to reserve a mist IP: %s", err.Error())
		}

		appModel.LocalIPs["logvac"] = logvacIP.String()
		appModel.LocalIPs["mist"] = mistIP.String()

	}

	// save the app
	if err := appModel.Save(); err != nil {
		display.ErrorTask()
		lumber.Error("app:reserveIPs:models:App:Save(): %s", err.Error())
		return fmt.Errorf("failed to persist IPs to the db: %s", err.Error())
	}

	return nil
}
