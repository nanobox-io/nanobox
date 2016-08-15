package dev

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/component"
	"github.com/nanobox-io/nanobox/processor/env"
)

// Deploy ...
type Deploy struct {
	Env models.Env
	App models.App
}

//
func (deploy Deploy) Run() error {
	// run the share init which gives access to docker
	envInit := env.Init{}
	if err := envInit.Run(); err != nil {
		return err
	}

	// syncronize the services as per the new boxfile
	componentSync := component.Sync{
		Env: deploy.Env,
		App: deploy.App,
	}
	if err := componentSync.Run(); err != nil {
		return err
	}

	return nil
}
