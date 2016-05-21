package code

import (
	"fmt"
	"net"
	"os"
	"io"

	"github.com/nanobox-io/nanobox-boxfile"
	dockType "github.com/docker/engine-api/types"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/util/dockerexec"
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
		self.runDebugSession(err)
		return err
	}

	if err := self.runConfigureHook(); err != nil {
		self.runDebugSession(err)
		return err
	}

	if err := self.runFetchHook(); err != nil {
		self.runDebugSession(err)
		return err
	}

	if err := self.runSetupHook(); err != nil {
		self.runDebugSession(err)
		return err
	}

	if err := self.runBoxfileHook(); err != nil {
		self.runDebugSession(err)
		return err
	}

	if err := self.runPrepareHook(); err != nil {
		self.runDebugSession(err)
		return err
	}

	if self.config.Meta["build"] == "true" {

		if err := self.runCompileHook(); err != nil {
			self.runDebugSession(err)
			return err
		}

		if err := self.runPackAppHook(); err != nil {
			self.runDebugSession(err)
			return err
		}

	}

	if err := self.runPackBuildHook(); err != nil {
		self.runDebugSession(err)
		return err
	}

	if self.config.Meta["build"] == "true" {

		if err := self.runCleanHook(); err != nil {
			self.runDebugSession(err)
			return err
		}

		if err := self.runPackDeployHook(); err != nil {
			self.runDebugSession(err)
			return err
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

	label := "Downloading docker image " + self.image
	fmt.Print(stylish.NestedProcessStart(label, self.config.DisplayLevel))

	// Create a pipe to for the JSONMessagesStream to read from
	pr, pw := io.Pipe()
	prefix := stylish.GenerateNestedPrefix(self.config.DisplayLevel + 1)
  go print.DisplayJSONMessagesStream(pr, os.Stdout, os.Stdout.Fd(), true, prefix, nil)
	if _, err := docker.ImagePull(self.image, pw); err != nil {
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
			fmt.Sprintf("/share/%s/code:/share/code", appName),
			fmt.Sprintf("/share/%s/engine:/share/engine", appName),
			fmt.Sprintf("/mnt/sda1/%s/build:/mnt/build", appName),
			fmt.Sprintf("/mnt/sda1/%s/deploy:/mnt/deploy", appName),
			fmt.Sprintf("/mnt/sda1/%s/app:/mnt/app", appName),
			fmt.Sprintf("/mnt/sda1/%s/cache:/mnt/cache", appName),
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
	cmd := dockerexec.Command(self.container.ID, "user", util.UserPayload())
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// runConfigureHook runs the configure hook in the build container
func (self *codeBuild) runConfigureHook() error {
	cmd := dockerexec.Command(self.container.ID, "configure", "{}")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// runFetchHook runs the fetch hook in the build container
func (self *codeBuild) runFetchHook() error {
	cmd := dockerexec.Command(self.container.ID, "fetch", "{}")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// runSetupHook runs the setup hook in the build container
func (self *codeBuild) runSetupHook() error {
	cmd := dockerexec.Command(self.container.ID, "setup", "{}")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// runBoxfileHook runs the boxfile hook in the build container
func (self *codeBuild) runBoxfileHook() error {
	cmd := dockerexec.Command(self.container.ID, "boxfile", "{}")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	output := cmd.Output()

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
	cmd := dockerexec.Command(self.container.ID, "prepare", "{}")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// runCompileHook runs the compile hook in the build container
func (self *codeBuild) runCompileHook() error {
	cmd := dockerexec.Command(self.container.ID, "compile", "{}")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// runPackAppHook runs the pack-app hook in the build container
func (self *codeBuild) runPackAppHook() error {
	cmd := dockerexec.Command(self.container.ID, "pack-app", "{}")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// runPackBuildHook runs the pack-build hook in the build container
func (self *codeBuild) runPackBuildHook() error {
	cmd := dockerexec.Command(self.container.ID, "pack-build", "{}")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// runCleanHook runs the clean hook in the build container
func (self *codeBuild) runCleanHook() error {
	cmd := dockerexec.Command(self.container.ID, "clean", "{}")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// runPackDeployHook runs the pack-deploy hook in the build container
func (self *codeBuild) runPackDeployHook() error {
	cmd := dockerexec.Command(self.container.ID, "pack-deploy", "{}")
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
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
	}

	return nil
}
