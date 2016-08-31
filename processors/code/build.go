package code

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox/commands/registry"
	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	hook_generator "github.com/nanobox-io/nanobox/generators/hooks/build"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
)

// Build builds the codebase that can then be deployed
func Build(env *models.Env) error {
	display.OpenContext("building code")
	defer display.CloseContext()

	// pull the latest build image
	buildImage, err := pullBuildImage()
	if err != nil {
		return fmt.Errorf("failed to pull the build image: %s", err.Error())
	}

	// if a build container was leftover from a previous build, let's remove it
	docker.ContainerRemove(container_generator.BuildName())

	// reserve an IP for the build container
	ip, err := dhcp.ReserveLocal()
	if err != nil {
		lumber.Error("code:Build:dhcp.ReserveLocal(): %s", err.Error())
		return fmt.Errorf("failed to reserve an ip for the build container: %s", err.Error())
	}

	// ensure we release the IP when we're done
	defer dhcp.ReturnIP(ip)

	// start the container
	display.StartTask("starting build container")
	config := container_generator.BuildConfig(buildImage, ip.String())
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Build:docker.CreateContainer(%+v): %s", config, err.Error())
		display.ErrorTask()
		return fmt.Errorf("failed to start docker container: %s", err.Error())
	}
	display.StopTask()

	// ensure we stop the container when we're done
	defer docker.ContainerRemove(container_generator.BuildName())

	display.StartTask("running build hooks")
	// run the user hook
	if _, err := hookit.RunUserHook(container.ID, hook_generator.UserPayload()); err != nil {
		display.ErrorTask()
		err = fmt.Errorf("failed to run user hook: %s", err.Error())
		return runDebugSession(container.ID, err)
	}

	// run the configure hook
	if _, err := hookit.RunConfigureHook(container.ID, hook_generator.ConfigurePayload()); err != nil {
		display.ErrorTask()
		err = fmt.Errorf("failed to run configure hook: %s", err.Error())
		return runDebugSession(container.ID, err)
	}

	// run the fetch hook
	if _, err := hookit.RunFetchHook(container.ID, hook_generator.FetchPayload()); err != nil {
		display.ErrorTask()
		err = fmt.Errorf("failed to run fetch hook: %s", err.Error())
		return runDebugSession(container.ID, err)
	}

	// run the setup hook
	if _, err := hookit.RunSetupHook(container.ID, hook_generator.SetupPayload()); err != nil {
		display.ErrorTask()
		err = fmt.Errorf("failed to run setup hook: %s", err.Error())
		return runDebugSession(container.ID, err)
	}

	// run the boxfile hook
	boxOutput, err := hookit.RunBoxfileHook(container.ID, hook_generator.BoxfilePayload())
	if err != nil {
		display.ErrorTask()
		err = fmt.Errorf("failed to run boxfile hook: %s", err.Error())
		return runDebugSession(container.ID, err)
	}

	// persist the boxfile output to the env model
	env.BuiltBoxfile = boxOutput
	if err := env.Save(); err != nil {
		lumber.Error("code:Build:models:Env:Save(): %s", err.Error())
		return fmt.Errorf("failed to persist build boxfile to db: %s", err.Error())
	}

	// run the prepare hook
	if _, err := hookit.RunPrepareHook(container.ID, hook_generator.PreparePayload()); err != nil {
		display.ErrorTask()
		err = fmt.Errorf("failed to run prepare hook: %s", err.Error())
		return runDebugSession(container.ID, err)
	}

	if !registry.GetBool("skip-compile") {
		// run the compile hook
		if _, err := hookit.RunCompileHook(container.ID, hook_generator.CompilePayload()); err != nil {
			display.ErrorTask()
			err = fmt.Errorf("failed to run compile hook: %s", err.Error())
			return runDebugSession(container.ID, err)
		}

		// run the pack-app hook
		if _, err := hookit.RunPackAppHook(container.ID, hook_generator.PackAppPayload()); err != nil {
			display.ErrorTask()
			err = fmt.Errorf("failed to run pack-app hook: %s", err.Error())
			return runDebugSession(container.ID, err)
		}
	}

	// run the pack-build hook
	if _, err := hookit.RunPackBuildHook(container.ID, hook_generator.PackBuildPayload()); err != nil {
		display.ErrorTask()
		err = fmt.Errorf("failed to run pack-build hook: %s", err.Error())
		return runDebugSession(container.ID, err)
	}

	if !registry.GetBool("skip-compile") {
		// run the clean hook
		if _, err := hookit.RunCleanHook(container.ID, hook_generator.CleanPayload()); err != nil {
			display.ErrorTask()
			err = fmt.Errorf("failed to run clean hook: %s", err.Error())
			return runDebugSession(container.ID, err)
		}

		// run the pack-deploy hook
		if _, err := hookit.RunPackDeployHook(container.ID, hook_generator.PackDeployPayload()); err != nil {
			display.ErrorTask()
			err = fmt.Errorf("failed to run pack-deploy hook: %s", err.Error())
			return runDebugSession(container.ID, err)
		}
	}

	display.StopTask()

	return nil
}
