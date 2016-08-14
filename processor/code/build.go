package code

import (
	"fmt"
	"net"
	"os"

	dockType "github.com/docker/engine-api/types"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/env"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/dhcp"
)

// Build ...
type Build struct {
	Env models.Env
	boxfile      boxfile.Boxfile
	localIP      net.IP
	image        string
	container    dockType.ContainerJSON
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

	// run the user hook in the build container
	if _, err := util.Exec(codeBuild.container.ID, "user", config.UserPayload(), nil); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the configure hook in the build container
	if _, err := util.Exec(codeBuild.container.ID, "configure", "{}", nil); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the fetch hook in the build container
	if _, err := util.Exec(codeBuild.container.ID, "fetch", "{}", nil); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the setup hook in the build container
	if _, err := util.Exec(codeBuild.container.ID, "setup", "{}", nil); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the boxfile hook in the build container also sets the boxfile in my meta
	// for later use
	if err := codeBuild.runBoxfileHook(); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the prepare hook in the build container
	if _, err := util.Exec(codeBuild.container.ID, "prepare", "{}", nil); err != nil {
		return codeBuild.runDebugSession(err)
	}

	if !registry.GetBool("no-compile") {
		// run the compile hook in the build container
		if _, err := util.Exec(codeBuild.container.ID, "compile", "{}", nil); err != nil {
			return codeBuild.runDebugSession(err)
		}

		// run the pack-app hook in the build container
		if _, err := util.Exec(codeBuild.container.ID, "pack-app", "{}", nil); err != nil {
			return codeBuild.runDebugSession(err)
		}

	}

	// run the pack-build hook in the build container
	if _, err := util.Exec(codeBuild.container.ID, "pack-build", "{}", nil); err != nil {
		return codeBuild.runDebugSession(err)
	}

	if !registry.GetBool("no-compile") {
		// run the clean hook in the build container
		if _, err := util.Exec(codeBuild.container.ID, "clean", "{}", nil); err != nil {
			return codeBuild.runDebugSession(err)
		}

		// run the pack-deploy hook in the build container
		if _, err := util.Exec(codeBuild.container.ID, "pack-deploy", "{}", nil); err != nil {
			return codeBuild.runDebugSession(err)
		}
	}

	return nil
}

// downloadImage downloads a build image
func (codeBuild *Build) downloadImage() error {
	// load the boxfile from the users file
	// because it is the only one we have
	codeBuild.boxfile = boxfile.NewFromPath(config.Boxfile())

	codeBuild.image = codeBuild.boxfile.Node("build").StringValue("image")

	if codeBuild.image == "" {
		codeBuild.image = "nanobox/build:v1"
	}

	// TODO: replace with displays tuff
	// prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(codeBuild.control.DisplayLevel+1), codeBuild.image)
	// if _, err := docker.ImagePull(codeBuild.image, &print.DockerPercentDisplay{Prefix: prefix}); err != nil {
	if _, err := docker.ImagePull(codeBuild.image, nil); err != nil {
		return err
	}

	return nil
}

// reserveIP reserves a local IP for the build container
func (codeBuild *Build) reserveIP() error {
	IP, err := dhcp.ReserveLocal()
	if err != nil {
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
		return err
	}

	codeBuild.container = container

	return nil
}

// stopContainer stops the docker build container
func (codeBuild *Build) stopContainer() error {
	return docker.ContainerRemove(codeBuild.container.ID)
}

// runBoxfileHook runs the boxfile hook in the build container
func (codeBuild *Build) runBoxfileHook() error {
	output, err := util.Exec(codeBuild.container.ID, "boxfile", "{}", nil)
	if err != nil {
		return err
	}

	codeBuild.Env.BuiltBoxfile = output

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
