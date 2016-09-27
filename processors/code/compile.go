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

	// if a build container was leftover from a previous build, let's remove it
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

	if err := prepareSimpleEnvironment(container.ID); err != nil {
		return err
	}

	if err := runCompile(container.ID); err != nil {
		return err
	}

	return nil
}

// prepareSimpleEnvironment runs hooks to prepare the build environment
func prepareSimpleEnvironment(containerID string) error {
	display.StartTask("Preparing environment for compile")
	defer display.StopTask()

	// run the user hook
	if _, err := hookit.DebugExec(containerID, "user", hook_generator.UserPayload(), "info"); err != nil {
		return err
	}

	// run the fetch hook
	if _, err := hookit.DebugExec(containerID, "fetch", hook_generator.FetchPayload(), "info"); err != nil {
		return err
	}

	return nil
}

// runCompile runs the hooks to compile the codebase
func runCompile(containerID string) error {

	display.StartTask("Compiling code")
	defer display.StopTask()

	// run the compile hook
	if _, err := hookit.DebugExec(containerID, "compile", hook_generator.CompilePayload(), "info"); err != nil {
		return err
	}

	// run the pack-app hook
	if _, err := hookit.DebugExec(containerID, "pack-app", hook_generator.PackAppPayload(), "info"); err != nil {
		return err
	}

	return nil
}
