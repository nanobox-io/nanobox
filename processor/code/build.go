package code

import (
	"fmt"
	"net"
	// "io"

	dockType "github.com/docker/engine-api/types"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
	"github.com/nanobox-io/nanobox/util/print"
)

type codeBuild struct {
	control       processor.ProcessControl
	boxfile      boxfile.Boxfile
	buildBoxfile models.Boxfile
	localIP      net.IP
	image        string
	container    dockType.ContainerJSON
}

func init() {
	processor.Register("code_build", codeBuildFunc)
}

func codeBuildFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return &codeBuild{control: control}, nil
}

func (self codeBuild) Results() processor.ProcessControl {
	return self.control
}

func (self *codeBuild) Process() error {
	self.control.Display(stylish.Bullet("Building Code"))

	if err := self.loadBoxfile(); err != nil {
		return err
	}

	if err := self.downloadImage(); err != nil {
		return err
	}

	if err := self.reserveIP(); err != nil {
		return err
	}

	defer self.releaseIP()

	if err := self.startContainer(); err != nil {
		return err
	}

	defer self.stopContainer()

	// run the user hook in the build container
	if _, err := util.Exec(self.container.ID, "user", util.UserPayload(), processor.ExecWriter()); err != nil {
		return self.runDebugSession(err)
	}

	// run the configure hook in the build container
	if _, err := util.Exec(self.container.ID, "configure", "{}", processor.ExecWriter()); err != nil {
		return self.runDebugSession(err)
	}

	// run the fetch hook in the build container
	if _, err := util.Exec(self.container.ID, "fetch", "{}", processor.ExecWriter()); err != nil {
		return self.runDebugSession(err)
	}

	// run the setup hook in the build container
	if _, err := util.Exec(self.container.ID, "setup", "{}", processor.ExecWriter()); err != nil {
		return self.runDebugSession(err)
	}

	// run the boxfile hook in the build container
	// also sets the boxfile in my meta for later use
	if err := self.runBoxfileHook(); err != nil {
		return self.runDebugSession(err)
	}

	// run the prepare hook in the build container
	if _, err := util.Exec(self.container.ID, "prepare", "{}", processor.ExecWriter()); err != nil {
		return self.runDebugSession(err)
	}

	if self.control.Meta["build"] == "true" {
		// run the compile hook in the build container
		if _, err := util.Exec(self.container.ID, "compile", "{}", processor.ExecWriter()); err != nil {
			return self.runDebugSession(err)
		}

		// run the pack-app hook in the build container
		if _, err := util.Exec(self.container.ID, "pack-app", "{}", processor.ExecWriter()); err != nil {
			return self.runDebugSession(err)
		}

	}

	// run the pack-build hook in the build container
	if _, err := util.Exec(self.container.ID, "pack-build", "{}", processor.ExecWriter()); err != nil {
		return self.runDebugSession(err)
	}

	if self.control.Meta["build"] == "true" {
		// run the clean hook in the build container
		if _, err := util.Exec(self.container.ID, "clean", "{}", processor.ExecWriter()); err != nil {
			return self.runDebugSession(err)
		}

		// run the pack-deploy hook in the build container
		if _, err := util.Exec(self.container.ID, "pack-deploy", "{}", processor.ExecWriter()); err != nil {
			return self.runDebugSession(err)
		}
	}

	return nil
}

// loadBoxfile loads the boxfile into the control state
func (self *codeBuild) loadBoxfile() error {
	self.boxfile = boxfile.NewFromPath(util.BoxfileLocation())

	return nil
}

// downloadImage downloads a build image
func (self *codeBuild) downloadImage() error {
	self.image = self.boxfile.Node("build").StringValue("image")

	if self.image == "" {
		self.image = "nanobox/build:v1"
	}

	// create the prefix for the image message
	prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(self.control.DisplayLevel+1), self.image)
	if _, err := docker.ImagePull(self.image, &print.DockerPercentDisplay{Prefix: prefix}); err != nil {
		return err
	}

	return nil
}

// reserveIP reserves a local IP for the build container
func (self *codeBuild) reserveIP() error {
	IP, err := ip_control.ReserveLocal()
	if err != nil {
		return err
	}

	self.localIP = IP

	return nil
}

// releaseIP releases a local IP back into the pool
func (self *codeBuild) releaseIP() error {
	return ip_control.ReturnIP(self.localIP)
}

// startContainer starts a build container
func (self *codeBuild) startContainer() error {

	appName := util.AppName()
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox-%s-build", util.AppName()),
		Image:   self.image, // this will need to be controlurable some time
		Network: "virt",
		IP:      self.localIP.String(),
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

	self.container = container

	return nil
}

// stopContainer stops the docker build container
func (self *codeBuild) stopContainer() error {
	return docker.ContainerRemove(self.container.ID)
}

// runBoxfileHook runs the boxfile hook in the build container
func (self *codeBuild) runBoxfileHook() error {
	output, err := util.Exec(self.container.ID, "boxfile", "{}", processor.ExecWriter())
	if err != nil {
		return err
	}

	// set the boxfile in the meta
	self.control.Meta["boxfile"] = output

	// store it in the database as well
	self.buildBoxfile.Data = []byte(output)

	return data.Put(util.AppName()+"_meta", "build_boxfile", self.buildBoxfile)
}

// runDebugSession drops the user in the build container to debug
func (self *codeBuild) runDebugSession(err error) error {
	fmt.Println("there has been a failure")
	if self.control.Verbose {
		fmt.Println(err)
		fmt.Println("we will be dropping you into the failed build container")
		fmt.Println("GOOD LUCK!")
		self.control.Meta["name"] = "build"
		err := processor.Run("dev_console", self.control)
		if err != nil {
			fmt.Println("unable to enter console", err)
		}
	} else {
		return err
	}

	return nil
}
