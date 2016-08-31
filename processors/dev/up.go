package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/config"
)

// Up ...
func Up() error {

	// run a nanobox start
	if err := processors.Start(); err != nil {
		return err
	}

	envModel, _ := models.FindEnvByID(config.EnvID())
	appModel, _ := models.FindAppBySlug(config.EnvID(), "dev")

	// run a nanobox build
	if err := processors.Build(envModel); err != nil {
		return err
	}

	// run a dev start
	if err := Start(envModel, appModel); err != nil {
		return err
	}

	// run a dev deploy
	if err := Deploy(envModel, appModel); err != nil {
		return err
	}

	return nil
}
