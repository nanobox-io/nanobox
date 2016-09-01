package processors

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/code"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/provider"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Build sets up the environment and runs a code build
func Build(envModel *models.Env) error {
	// by aquiring a local lock we are only allowing
	// one build to happen at a time
	locker.LocalLock()
	defer locker.LocalUnlock()

	// setup the provider
	if err := provider.Setup(); err != nil {
		return fmt.Errorf("failed to setup the provider: %s", err.Error())
	}

	// setup the env
	if err := env.Setup(envModel); err != nil {
		return fmt.Errorf("failed to setup the env: %s", err.Error())
	}

	// build code
	if err := code.Build(envModel); err != nil {
		return fmt.Errorf("failed to build the code: %s", err.Error())
	}

	return nil
}
