package processor

import (
	"github.com/jcelliott/lumber"
	
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/code"
	"github.com/nanobox-io/nanobox/processor/env"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/display"
)

// Build ...
type Build struct {
	Env models.Env
}

//
func (build *Build) Run() error {
	display.OpenContext("running build")
	defer display.CloseContext()

	// by aquiring a local lock we are only allowing
	// one build to happen at a time
	locker.LocalLock()
	defer locker.LocalUnlock()

	setup := &env.Setup{}
	// get the vm and app up.
	if err := setup.Run(); err != nil {
		return err
	}

	// build code
	codeBuild := &code.Build{Env: setup.Env}
	err := codeBuild.Run()
	if err != nil {
		return err
	}

	build.Env = codeBuild.Env

	lumber.Debug("processor:build:Env: %+v", build.Env)
	
	return nil
}
