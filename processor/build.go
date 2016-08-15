package processor

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/code"
	"github.com/nanobox-io/nanobox/processor/env"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Build ...
type Build struct {
	Env models.Env
}

//
func (build Build) Run() error {

	// by aquiring a local lock we are only allowing
	// one build to happen at a time
	locker.LocalLock()
	defer locker.LocalUnlock()

	setup := env.Setup{}
	// get the vm and app up.
	if err := setup.Run(); err != nil {
		return err
	}

	build.Env = setup.Env

	// build code
	codeBuild := code.Build{Env: build.Env}

	return codeBuild.Run()
}
