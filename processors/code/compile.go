package code

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	hook_generator "github.com/nanobox-io/nanobox/generators/hooks/build"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
)

// Compile builds the codebase that can then be deployed
func Compile(envModel *models.Env) error {
	display.OpenContext("Compiling application")
	defer display.CloseContext()

	// pull the latest build image
	buildImage, err := pullBuildImage()
	if err != nil {
		return fmt.Errorf("failed to pull the compile image: %s", err.Error())
	}

	// if a compile container was leftover from a previous compile, let's remove it
	docker.ContainerRemove(container_generator.CompileName())

	display.StartTask("Starting docker container")

	// start the container
	config := container_generator.CompileConfig(buildImage)
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Compile:docker.CreateContainer(%+v): %s", config, err.Error())
		return fmt.Errorf("failed to start docker container: %s", err.Error())
	}

	display.StopTask()

	// ensure we stop the container when we're done
	defer docker.ContainerRemove(container_generator.CompileName())

	if err := prepareCompileEnvironment(container.ID); err != nil {
		return err
	}

	if err := compileCode(container.ID); err != nil {
		return err
	}

	return nil
}

// prepareCompileEnvironment runs hooks to prepare the build environment
func prepareCompileEnvironment(containerID string) error {
	display.StartTask("Preparing environment for compile")
	defer display.StopTask()

	// run the user hook
	if _, err := hookit.RunUserHook(containerID, hook_generator.UserPayload()); err != nil {
		display.ErrorTask()
		return runDebugSession(containerID, err)
	}
	
	// run the configure hook
	if _, err := hookit.RunConfigureHook(containerID, hook_generator.ConfigurePayload()); err != nil {
		display.ErrorTask()
		return runDebugSession(containerID, err)
	}
	
	// run the boxfile hook
	if _, err := hookit.RunBoxfileHook(containerID, hook_generator.BoxfilePayload()); err != nil {
		display.ErrorTask()
		return runDebugSession(containerID, err)
	}

	return nil
}

// compileCode runs the hooks to compile the codebase
func compileCode(containerID string) error {

	display.StartTask("Compiling code")
	defer display.StopTask()

	// run the compile hook
	if _, err := hookit.RunCompileHook(containerID, hook_generator.CompilePayload()); err != nil {
		display.ErrorTask()
		return runDebugSession(containerID, err)
	}

	// run the pack-app hook
	if _, err := hookit.RunPackAppHook(containerID, hook_generator.PackAppPayload()); err != nil {
		display.ErrorTask()
		return runDebugSession(containerID, err)
	}

	return nil
}
