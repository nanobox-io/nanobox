package processors

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/code"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
)

// Build sets up the environment and runs a code build
func Build(envModel *models.Env) error {
	// by aquiring a local lock we are only allowing
	// one build to happen at a time
	locker.LocalLock()
	defer locker.LocalUnlock()

	// init docker client and env mounts
	if err := env.Setup(envModel); err != nil {
		return util.ErrorAppend(err, "failed to init docker client")
	}

	// print a warning if this is the first build
	if envModel.BuiltBoxfile == "" {
		display.FirstBuild()
	}

	// build code
	if err := code.Build(envModel); err != nil {
		return util.ErrorAppend(err, "failed to build the code")
	}

	return nil
}
