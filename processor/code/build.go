package code

import (
	"fmt"
	"net"
	"os"
	// "io"

	"github.com/nanobox-io/nanobox-boxfile"
	dockType "github.com/docker/engine-api/types"
	"github.com/nanobox-io/golang-docker-client"
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
	config 				processor.ProcessConfig
	boxfile 			boxfile.Boxfile
	buildBoxfile 	models.Boxfile
	localIP 			net.IP
	image					string
	container 		dockType.ContainerJSON
}

func init() {
	processor.Register("code_build", codeBuildFunc)
}

func codeBuildFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return &codeBuild{config: config}, nil
}

func (self codeBuild) Results() processor.ProcessConfig {
	return self.config
}

func (self *codeBuild) Process() error {

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

	if err := self.runUserHook(); err != nil {
		return self.runDebugSession(err)
	}

	if err := self.runConfigureHook(); err != nil {
		return self.runDebugSession(err)
	}

	if err := self.runFetchHook(); err != nil {
		return self.runDebugSession(err)
	}

	if err := self.runSetupHook(); err != nil {
		return self.runDebugSession(err)
	}

	if err := self.runBoxfileHook(); err != nil {
		return self.runDebugSession(err)
	}

	if err := self.runPrepareHook(); err != nil {
		return self.runDebugSession(err)
	}

	if self.config.Meta["build"] == "true" {

		if err := self.runCompileHook(); err != nil {
			return self.runDebugSession(err)
		}

		if err := self.runPackAppHook(); err != nil {
			return self.runDebugSession(err)
		}

	}

	if err := self.runPackBuildHook(); err != nil {
		return self.runDebugSession(err)
	}

	if self.config.Meta["build"] == "true" {

		if err := self.runCleanHook(); err != nil {
			return self.runDebugSession(err)
		}

		if err := self.runPackDeployHook(); err != nil {
			return self.runDebugSession(err)
		}

	}

	return nil
}

// loadBoxfile loads the boxfile into the config state
func (self *codeBuild) loadBoxfile() error {
	// can this error?
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
	prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(self.config.DisplayLevel), self.image)
	if _, err := docker.ImagePull(self.image, &print.DockerPercentDisplay{Prefix: prefix}); err != nil {
		return err
	}
  fmt.Print(stylish.ProcessEnd())

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
		Name:    fmt.Sprintf("%s-build", util.AppName()),
		Image:   self.image, // this will need to be configurable some time
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

// runUserHook runs the user hook in the build container
func (self *codeBuild) runUserHook() error {
	_, err := util.Exec(self.container.ID, "user", util.UserPayload(), os.Stdout)
	return err
}

// runConfigureHook runs the configure hook in the build container
func (self *codeBuild) runConfigureHook() error {
	_, err := util.Exec(self.container.ID, "configure", "{}", os.Stdout)
	return err
}

// runFetchHook runs the fetch hook in the build container
func (self *codeBuild) runFetchHook() error {
	_, err := util.Exec(self.container.ID, "fetch", "{}", os.Stdout)
	return err
}

// runSetupHook runs the setup hook in the build container
func (self *codeBuild) runSetupHook() error {
	_, err := util.Exec(self.container.ID, "setup", "{}", os.Stdout)
	return err
}

// runBoxfileHook runs the boxfile hook in the build container
func (self *codeBuild) runBoxfileHook() error {
	output, err := util.Exec(self.container.ID, "boxfile", "{}", os.Stdout)
	if err != nil {
		return err
	}

	// set the boxfile in the meta
	self.config.Meta["boxfile"] = output

	// store it in the database as well
	self.buildBoxfile.Data = []byte(output)

	if err := data.Put(util.AppName()+"_meta", "build_boxfile", self.buildBoxfile); err != nil {
		return err
	}

	return nil
}

// runPrepareHook runs the prepare hook in the build container
func (self *codeBuild) runPrepareHook() error {
	_, err := util.Exec(self.container.ID, "prepare", "{}", os.Stdout)
	return err
}

// runCompileHook runs the compile hook in the build container
func (self *codeBuild) runCompileHook() error {
	_, err := util.Exec(self.container.ID, "compile", "{}", os.Stdout)
	return err
}

// runPackAppHook runs the pack-app hook in the build container
func (self *codeBuild) runPackAppHook() error {
	_, err := util.Exec(self.container.ID, "pack-app", "{}", os.Stdout)
	return err
}

// runPackBuildHook runs the pack-build hook in the build container
func (self *codeBuild) runPackBuildHook() error {
	_, err := util.Exec(self.container.ID, "pack-build", "{}", os.Stdout)
	return err
}

// runCleanHook runs the clean hook in the build container
func (self *codeBuild) runCleanHook() error {
	_, err := util.Exec(self.container.ID, "clean", "{}", os.Stdout)
	return err
}

// runPackDeployHook runs the pack-deploy hook in the build container
func (self *codeBuild) runPackDeployHook() error {
	_, err := util.Exec(self.container.ID, "pack-deploy", "{}", os.Stdout)
	return err
}

// runDebugSession drops the user in the build container to debug
func (self *codeBuild) runDebugSession(err error) error {
	fmt.Println("there has been a failure")
	if self.config.Verbose {
		fmt.Println(err)
		fmt.Println("we will be dropping you into the failed build container")
		fmt.Println("GOOD LUCK!")
		self.config.Meta["name"] = "build"
		err := processor.Run("dev_console", self.config)
		if err != nil {
			fmt.Println("unable to enter console", err)
		}
	} else {
		return err
	}

	return nil
}
