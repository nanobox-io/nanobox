package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// Up ...
func Up() error {

	// run a nanobox start
	display.OpenContext("(nanobox start)")
	if err := processors.Start(); err != nil {
		return err
	}
	display.CloseContext()

	envModel, _ := models.FindEnvByID(config.EnvID())
	appModel, _ := models.FindAppBySlug(config.EnvID(), "dev")

	// run a nanobox build
	display.OpenContext("(nanobox build)")
	if err := processors.Build(envModel); err != nil {
		return err
	}
	display.CloseContext()

	// run a dev start
	display.OpenContext("(nanobox dev start)")
	if err := Start(envModel, appModel); err != nil {
		return err
	}
	display.CloseContext()

	// run a dev deploy
	display.OpenContext("(nanobox dev deploy)")
	if err := Deploy(envModel, appModel); err != nil {
		return err
	}
	display.CloseContext()

	return nil
}
