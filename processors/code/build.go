package code

import (
	"fmt"
	"net"
	"os"

	dockType "github.com/docker/engine-api/types"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
)

// Build ...
type Build struct {
	Env       models.Env
	boxfile   boxfile.Boxfile
	localIP   net.IP
	image     string
	container dockType.ContainerJSON
}

//
func (codeBuild *Build) Run() error {

	// remove any leftover build containers that may exist
	docker.ContainerRemove(fmt.Sprintf("nanobox_%s_build", config.EnvID()))

	//
	if err := codeBuild.downloadImage(); err != nil {
		return err
	}

	//
	if err := codeBuild.reserveIP(); err != nil {
		return err
	}
	defer codeBuild.releaseIP()

	//
	if err := codeBuild.startContainer(); err != nil {
		return err
	}
	defer codeBuild.stopContainer()

	// set the prefix so the utilExec lumber logging has context
	lumber.Prefix("code:Build")
	defer lumber.Prefix("")

	display.StartTask("running Hooks")

	// run the user hook in the build container
	display.Debug("running user hook")
	if _, err := util.Exec(codeBuild.container.ID, "user", config.UserPayload(), display.NewStreamer("debug")); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the configure hook in the build container
	display.Info("running configure hook")
	if _, err := util.Exec(codeBuild.container.ID, "configure", "{}", display.NewStreamer("info")); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the fetch hook in the build container
	display.Info("running fetch hook")
	if _, err := util.Exec(codeBuild.container.ID, "fetch", "{}", display.NewStreamer("info")); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the setup hook in the build container
	display.Info("running setup hook")
	if _, err := util.Exec(codeBuild.container.ID, "setup", "{}", display.NewStreamer("info")); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the boxfile hook in the build container also sets the boxfile in my meta
	// for later use
	if err := codeBuild.runBoxfileHook(); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the prepare hook in the build container
	display.Info("running prepare hook")
	if _, err := util.Exec(codeBuild.container.ID, "prepare", "{}", display.NewStreamer("info")); err != nil {
		return codeBuild.runDebugSession(err)
	}

	if !registry.GetBool("no-compile") {
		// run the compile hook in the build container
		display.Info("running compile hook")
		if _, err := util.Exec(codeBuild.container.ID, "compile", "{}", display.NewStreamer("info")); err != nil {
			return codeBuild.runDebugSession(err)
		}

		// run the pack-app hook in the build container
		display.Debug("running pack-app hook")
		if _, err := util.Exec(codeBuild.container.ID, "pack-app", "{}", display.NewStreamer("debug")); err != nil {
			return codeBuild.runDebugSession(err)
		}

	}

	// run the pack-build hook in the build container
	display.Debug("running pack-build hook")
	if _, err := util.Exec(codeBuild.container.ID, "pack-build", "{}", display.NewStreamer("info")); err != nil {
		return codeBuild.runDebugSession(err)
	}

	if !registry.GetBool("no-compile") {
		// run the clean hook in the build container
		display.Debug("running clean hook")
		if _, err := util.Exec(codeBuild.container.ID, "clean", "{}", display.NewStreamer("info")); err != nil {
			return codeBuild.runDebugSession(err)
		}

		// run the pack-deploy hook in the build container
		display.Debug("running pack-deploy hook")
		if _, err := util.Exec(codeBuild.container.ID, "pack-deploy", "{}", display.NewStreamer("debug")); err != nil {
			return codeBuild.runDebugSession(err)
		}
	}

	lumber.Debug("build:end:env: %+v", codeBuild.Env)
	display.StopTask()

	return nil
}

// downloadImage downloads a build image
func (codeBuild *Build) downloadImage() error {
	display.StartTask("downloading image")
	// load the boxfile from the users file
	// because it is the only one we have
	codeBuild.boxfile = boxfile.NewFromPath(config.Boxfile())

	codeBuild.image = codeBuild.boxfile.Node("build").StringValue("image")

	if codeBuild.image == "" {
		codeBuild.image = "nanobox/build:v1"
	}

	streamer := display.NewStreamer("info")	
	dockerPercent := &display.DockerPercentDisplay{Output: streamer, Prefix: codeBuild.image}
	if _, err := docker.ImagePull(codeBuild.image, dockerPercent); err != nil {
		lumber.Error("code:Build:downloadImage:docker.ImagePull(%s, nil): %s", codeBuild.image, err.Error())
		display.ErrorTask()
		return err
	}

	display.StopTask()
	return nil
}

// reserveIP reserves a local IP for the build container
func (codeBuild *Build) reserveIP() error {
	IP, err := dhcp.ReserveLocal()
	if err != nil {
		lumber.Error("code:Build:reserveIP:dhcp.ReserveLocal(): %s", err.Error())
		return err
	}

	codeBuild.localIP = IP

	return nil
}

// releaseIP releases a local IP back into the pool
func (codeBuild *Build) releaseIP() error {
	return dhcp.ReturnIP(codeBuild.localIP)
}

// startContainer starts a build container
func (codeBuild *Build) startContainer() error {
	display.StartTask("starting build container")

	appName := config.EnvID()
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox_%s_build", config.EnvID()),
		Image:   codeBuild.image, // this will need to be controlurable some time
		Network: "virt",
		IP:      codeBuild.localIP.String(),
		Binds: []string{
			fmt.Sprintf("%s%s/code:/share/code", provider.HostShareDir(), appName),
			fmt.Sprintf("%s%s/engine:/share/engine", provider.HostShareDir(), appName),
			fmt.Sprintf("%s%s/build:/mnt/build", provider.HostMntDir(), appName),
			fmt.Sprintf("%s%s/deploy:/mnt/deploy", provider.HostMntDir(), appName),
			fmt.Sprintf("%s%s/app:/mnt/app", provider.HostMntDir(), appName),
			fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), appName),
		},
	}

	// start container
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("code:Build:startContainer:docker.CreateContainer(%+v): %s", config, err.Error())
		display.ErrorTask()
		return err
	}

	codeBuild.container = container
	display.StopTask()
	return nil
}

// stopContainer stops the docker build container
func (codeBuild *Build) stopContainer() error {
	if err := docker.ContainerRemove(codeBuild.container.ID); err != nil {
		lumber.Error("code:Build:stopContainer:docker.ContainerRemove(%s): %s", codeBuild.container.ID, err.Error())
		return err
	}
	return nil
}

// runBoxfileHook runs the boxfile hook in the build container
func (codeBuild *Build) runBoxfileHook() error {
	output, err := util.Exec(codeBuild.container.ID, "boxfile", "{}", display.NewStreamer("info"))
	if err != nil {
		return err
	}

	codeBuild.Env.BuiltBoxfile = output
	lumber.Debug("build:boxfilehook:env: %+v", codeBuild.Env)

	return codeBuild.Env.Save()
}

// runDebugSession drops the user in the build container to debug
func (codeBuild *Build) runDebugSession(err error) error {
	fmt.Println("there has been a failure")
	if registry.GetBool("debug") {
		fmt.Println(err)
		fmt.Println("we will be dropping you into the failed build container")
		fmt.Println("GOOD LUCK!")
		component := models.Component{
			ID: codeBuild.container.ID,
		}
		envConsole := env.Console{
			Component: component,
		}
		err := envConsole.Run()
		if err != nil {
			fmt.Println("unable to enter console", err)
			os.Exit(1)
		}
	} else {
		return err
	}

	return nil
}
