package code

import (
	"strings"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"

	"github.com/nanobox-io/nanobox-boxfile"
	container_generator "github.com/nanobox-io/nanobox/generators/containers"
	hook_generator "github.com/nanobox-io/nanobox/generators/hooks/build"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/hookit"
)

// Build builds the codebase that can then be deployed
func Build(envModel *models.Env) error {
	display.OpenContext("Building runtime")
	defer display.CloseContext()

	// pull the latest build image
	buildImage, err := pullBuildImage()
	if err != nil {
		return util.ErrorAppend(err, "failed to pull the build image")
	}

	// if a build container was leftover from a previous build, let's remove it
	docker.ContainerRemove(container_generator.BuildName())

	display.StartTask("Starting docker container")

	// start the container
	contConfig := container_generator.BuildConfig(buildImage)
	container, err := docker.CreateContainer(contConfig)
	if err != nil {
		lumber.Error("code:Build:docker.CreateContainer(%+v): %s", contConfig, err.Error())
		return util.ErrorAppend(err, "failed to start docker container")
	}

	display.StopTask()

	if err := prepareBuildEnvironment(container.ID); err != nil {
		return err
	}

	if err := gatherRequirements(envModel, container.ID); err != nil {
		return err
	}

	populateBuildTriggers(envModel)

	if err := setupBuildMounts(container.ID); err != nil {
		return err
	}

	if err := installRuntimes(container.ID); err != nil {
		return err
	}

	if err := packageBuild(container.ID); err != nil {
		return err
	}

	envModel.LastBuild = time.Now()
	envModel.Save()

	// ensure we stop the container when we're done
	if err := docker.ContainerRemove(container_generator.BuildName()); err != nil {
		return util.ErrorAppend(err, "unable to remove docker contianer")
	}

	return nil
}

// prepareBuildEnvironment runs hooks to prepare the build environment
func prepareBuildEnvironment(containerID string) error {
	display.StartTask("Preparing environment for build")
	defer display.StopTask()

	// run the user hook
	if out, err := hookit.DebugExec(containerID, "user", hook_generator.UserPayload(), "info"); err != nil {
		// handle 'exec failed: argument list too long' error
		if strings.Contains(out, "argument list too long") {
			if err2, ok := err.(util.Err); ok {
				err2.Suggest = "You may have too many ssh keys, please specify the one you need with `nanobox config set ssh-key ~/.ssh/id_rsa`"
				err2.Output = out
				err2.Code = "1001"
				return util.ErrorAppend(err2, "failed to run the user hook")
			}
		}
		return util.ErrorAppend(err, "failed to run the user hook")
	}

	// run the configure hook
	if out, err := hookit.DebugExec(containerID, "configure", hook_generator.ConfigurePayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the configure hook")
		}
		return util.ErrorAppend(err, "failed to run configure hook")
	}

	// run the fetch hook
	if out, err := hookit.DebugExec(containerID, "fetch", hook_generator.FetchPayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the fetch hook")
		}
		return util.ErrorAppend(err, "failed to run the fetch hook")
	}

	// run the setup hook
	if out, err := hookit.DebugExec(containerID, "setup", hook_generator.SetupPayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the setup hook")
		}
		return util.ErrorAppend(err, "failed to run setup hook")
	}

	return nil
}

// gatherRequirements runs hooks to gather requirements
func gatherRequirements(envModel *models.Env, containerID string) error {
	display.StartTask("Gathering requirements")
	defer display.StopTask()

	// run the boxfile hook
	boxOutput, err := hookit.DebugExec(containerID, "boxfile", hook_generator.BoxfilePayload(), "info")
	if err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = boxOutput
			return util.ErrorAppend(err2, "failed to run the boxfile hook")
		}
		return util.ErrorAppend(err, "failed to run boxfile hook")
	}

	box := boxfile.NewFromPath(config.Boxfile())

	// set the boxfile data but do not save
	// if something else here fails we want to only save at the end
	envModel.UserBoxfile = box.String()
	envModel.BuiltBoxfile = boxOutput
	envModel.BuiltID = util.RandomString(30)

	return nil
}

// populate the build triggers so we can know next time if a change has happened
func populateBuildTriggers(envModel *models.Env) {
	if envModel.BuildTriggers == nil {
		envModel.BuildTriggers = map[string]string{}
	}
	box := boxfile.New([]byte(envModel.BuiltBoxfile))
	for _, trigger := range box.Node("run.config").StringSliceValue("build_triggers") {
		envModel.BuildTriggers[trigger] = util.FileMD5(trigger)
	}
}

// setupBuildMounts prepares the environment for the build
func setupBuildMounts(containerID string) error {
	display.StartTask("Mounting cache_dirs")
	defer display.StopTask()

	// run the mount hook
	if out, err := hookit.DebugExec(containerID, "mount", hook_generator.MountPayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the mount hook")
		}
		return util.ErrorAppend(err, "failed to run mount hook")
	}

	return nil
}

// installRuntimes runs the hooks to install binaries and runtimes
func installRuntimes(containerID string) error {
	display.StartTask("Installing binaries and runtimes")
	defer display.StopTask()

	// run the build hook
	if out, err := hookit.DebugExec(containerID, "build", hook_generator.BuildPayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the build hook")
		}
		return util.ErrorAppend(err, "failed to run build hook")
	}

	return nil
}

// packageBuild runs the hooks to package the build
func packageBuild(containerID string) error {
	display.StartTask("Packaging build")
	defer display.StopTask()

	// run the pack-build hook
	if out, err := hookit.DebugExec(containerID, "pack-build", hook_generator.PackBuildPayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the pack-build hook")
		}
		return util.ErrorAppend(err, "failed to run pack-build hook")
	}

	// run the clean hook
	if out, err := hookit.DebugExec(containerID, "clean", hook_generator.CleanPayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the clean hook")
		}
		return util.ErrorAppend(err, "failed to run clean hook")
	}

	// run the pack-deploy hook
	if out, err := hookit.DebugExec(containerID, "pack-deploy", hook_generator.PackDeployPayload(), "info"); err != nil {
		if err2, ok := err.(util.Err); ok {
			err2.Output = out
			return util.ErrorAppend(err2, "failed to run the pack-deploy hook")
		}
		return util.ErrorAppend(err, "failed to run pack-deploy hook")
	}

	return nil
}
