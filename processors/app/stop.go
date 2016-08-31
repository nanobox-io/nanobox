package app

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/component"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Stop will stop all services associated with an app
func Stop(appModel *models.App) error {
	locker.LocalLock()
	defer locker.LocalUnlock()

	// short-circuit if the app is already down
	// TODO: also check if any containers are running
	if appModel.Status != "up" {
		return nil
	}

	// load the env for the display context
	envModel, err := appModel.Env()
	if err != nil {
		lumber.Error("app:Stop:models.App.Env(): %s", err.Error())
		return fmt.Errorf("failed to load app env: %s", err.Error())
	}

	display.OpenContext("%s (%s)", envModel.Name, appModel.Name)
	defer display.CloseContext()

	// initialize docker for the provider
	if err := provider.Init(); err != nil {
		return fmt.Errorf("failed to initialize docker environment: %s", err.Error())
	}

	// stop all app components
	if err := component.StopAll(appModel); err != nil {
		return fmt.Errorf("failed to stop all app components: %s", err.Error())
	}

	// stop any dev containers
	docker.ContainerRemove(fmt.Sprintf("nanobox_%s", appModel.ID))

	// set the status to down
	appModel.Status = "down"
	if err := appModel.Save(); err != nil {
		lumber.Error("app:Stop:models.App.Save(): %s", err.Error())
		return fmt.Errorf("failed to persist app status: %s", err.Error())
	}

	return nil
}
