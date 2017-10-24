package env

import (
	"fmt"
	"path/filepath"

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

	// BUG(glinton) if the `build` isn't successful and the engine changes in the boxfile,
	// this mount sticks around and must be cleaned up manually.
	//
	// mount the engine if it's a local directory
	engineDir, _ := config.EngineDir()
	if engineDir != "" {
		src := engineDir                                                // local directory
		dst := filepath.Join(provider.HostShareDir(), env.ID, "engine") // b2d "global zone"

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
