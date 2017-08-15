package env

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Mount sets up the env mounts
func Mount(env *models.Env) error {

	if !provider.RequiresMount() {
		return nil
	}

	display.StartTask("Mounting codebase")
	defer display.StopTask()

	// mount the engine if it's a local directory
	if config.EngineDir() != "" {
		src := config.EngineDir()
		dst := fmt.Sprintf("%s%s/engine", provider.HostShareDir(), env.ID)

		// first, export the env on the workstation
		if err := provider.AddMount(src, dst); err != nil {
			display.ErrorTask()
			return util.ErrorAppend(err, "failed to mount the engine share on the provider")
		}
	}

	// mount the app src
	src := env.Directory
	dst := fmt.Sprintf("%s%s/code", provider.HostShareDir(), env.ID)

	// first export the env on the workstation
	if err := provider.AddMount(src, dst); err != nil {
		display.ErrorTask()
		return util.ErrorAppend(err, "failed to mount the code share on the provider")
	}

	// // setup mount directories
	// provider.Run([]string{"mkdir", "-p", fmt.Sprintf("%s%s/build", provider.HostMntDir(), env.ID)})
	// provider.Run([]string{"mkdir", "-p", fmt.Sprintf("%s%s/deploy", provider.HostMntDir(), env.ID)})
	// provider.Run([]string{"mkdir", "-p", fmt.Sprintf("%s%s/cache", provider.HostMntDir(), env.ID)})

	return nil
}
