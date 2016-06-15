package code

import (
	"fmt"
	"net"

	dockType "github.com/docker/engine-api/types"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/print"
)

// processCodeBuild ...
type processCodeBuild struct {
	control      processor.ProcessControl
	boxfile      boxfile.Boxfile
	buildBoxfile models.Boxfile
	localIP      net.IP
	image        string
	container    dockType.ContainerJSON
}

//
func init() {
	processor.Register("code_build", codeBuildFn)
}

//
func codeBuildFn(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	return &processCodeBuild{control: control}, nil
}

//
func (codeBuild processCodeBuild) Results() processor.ProcessControl {
	return codeBuild.control
}

//
func (codeBuild *processCodeBuild) Process() error {
	codeBuild.control.Display(stylish.Bullet("Building Code"))

	//
	if err := codeBuild.loadBoxfile(); err != nil {
		return err
	}

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
	if _, err := util.Exec(codeBuild.container.ID, "user", config.UserPayload(), processor.ExecWriter()); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the configure hook in the build container
	if _, err := util.Exec(codeBuild.container.ID, "configure", "{}", processor.ExecWriter()); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the fetch hook in the build container
	if _, err := util.Exec(codeBuild.container.ID, "fetch", "{}", processor.ExecWriter()); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the setup hook in the build container
	if _, err := util.Exec(codeBuild.container.ID, "setup", "{}", processor.ExecWriter()); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the boxfile hook in the build container also sets the boxfile in my meta
	// for later use
	if err := codeBuild.runBoxfileHook(); err != nil {
		return codeBuild.runDebugSession(err)
	}

	// run the prepare hook in the build container
	if _, err := util.Exec(codeBuild.container.ID, "prepare", "{}", processor.ExecWriter()); err != nil {
		return codeBuild.runDebugSession(err)
	}

	//
	if codeBuild.control.Meta["build"] == "true" {

		// run the compile hook in the build container
		if _, err := util.Exec(codeBuild.container.ID, "compile", "{}", processor.ExecWriter()); err != nil {
			return codeBuild.runDebugSession(err)
		}

		// run the pack-app hook in the build container
		if _, err := util.Exec(codeBuild.container.ID, "pack-app", "{}", processor.ExecWriter()); err != nil {
			return codeBuild.runDebugSession(err)
		}
	}

	// run the pack-build hook in the build container
	if _, err := util.Exec(codeBuild.container.ID, "pack-build", "{}", processor.ExecWriter()); err != nil {
		return codeBuild.runDebugSession(err)
	}

	//
	if codeBuild.control.Meta["build"] == "true" {

		// run the clean hook in the build container
		if _, err := util.Exec(codeBuild.container.ID, "clean", "{}", processor.ExecWriter()); err != nil {
			return codeBuild.runDebugSession(err)
		}

		// run the pack-deploy hook in the build container
		if _, err := util.Exec(codeBuild.container.ID, "pack-deploy", "{}", processor.ExecWriter()); err != nil {
			return codeBuild.runDebugSession(err)
		}
	}

	return nil
}

// loadBoxfile loads the boxfile into the control state
func (codeBuild *processCodeBuild) loadBoxfile() error {
	codeBuild.boxfile = boxfile.NewFromPath(config.Boxfile())

	return nil
}

// downloadImage downloads a build image
func (codeBuild *processCodeBuild) downloadImage() error {
	codeBuild.image = codeBuild.boxfile.Node("build").StringValue("image")

	if codeBuild.image == "" {
		codeBuild.image = "nanobox/build:v1"
	}

	// create the prefix for the image message
	prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(codeBuild.control.DisplayLevel+1), codeBuild.image)
	if _, err := docker.ImagePull(codeBuild.image, &print.DockerPercentDisplay{Prefix: prefix}); err != nil {
		return err
	}

	return nil
}

// reserveIP reserves a local IP for the build container
func (codeBuild *processCodeBuild) reserveIP() error {
	IP, err := dhcp.ReserveLocal()
	if err != nil {
		return err
	}

	codeBuild.localIP = IP

	return nil
}

// releaseIP releases a local IP back into the pool
func (codeBuild *processCodeBuild) releaseIP() error {
	return dhcp.ReturnIP(codeBuild.localIP)
}

// startContainer starts a build container
func (codeBuild *processCodeBuild) startContainer() error {

	appName := config.AppName()
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox-%s-build", config.AppName()),
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
func (codeBuild *processCodeBuild) stopContainer() error {
	return docker.ContainerRemove(codeBuild.container.ID)
}

// runBoxfileHook runs the boxfile hook in the build container
func (codeBuild *processCodeBuild) runBoxfileHook() error {
	output, err := util.Exec(codeBuild.container.ID, "boxfile", "{}", processor.ExecWriter())
	if err != nil {
		return err
	}

	// set the boxfile in the meta
	codeBuild.control.Meta["boxfile"] = output

	// store it in the database as well
	codeBuild.buildBoxfile.Data = []byte(output)

	return data.Put(config.AppName()+"_meta", "build_boxfile", codeBuild.buildBoxfile)
}

// runDebugSession drops the user in the build container to debug
func (codeBuild *processCodeBuild) runDebugSession(err error) error {
	fmt.Println("there has been a failure")
	if codeBuild.control.Verbose {
		fmt.Println(err)
		fmt.Println("we will be dropping you into the failed build container")
		fmt.Println("GOOD LUCK!")
		codeBuild.control.Meta["name"] = "build"
		err := processor.Run("dev_console", codeBuild.control)
		if err != nil {
			fmt.Println("unable to enter console", err)
		}
	} else {
		return err
	}

	return nil
}
