package code

import (
	"fmt"
	"net"
	"io"
	"os"

	dockType "github.com/docker/engine-api/types"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/util/dockerexec"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/locker"
)

type codeDev struct {
	config		processor.ProcessConfig
	boxfile 	models.Boxfile
	localIP		net.IP
	image 		string
	container	dockType.ContainerJSON
}

func init() {
	processor.Register("code_dev", codeDevFunc)
}

func codeDevFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return &codeDev{config: config}, nil
}

func (self codeDev) Results() processor.ProcessConfig {
	return self.config
}

func (self *codeDev) Process() error {

	if err := self.setup(); err != nil {
		// todo: how to display this?
		goto CLEANUP
	}

	if err := self.runConsole(); err != nil {
		// todo: how to display this?
		goto CLEANUP
	}

CLEANUP:

	if err := self.teardown(); err != nil {
		// this is bad...
		// we should probably print a pretty message explaining that the app
		// was left in a bad state and needs to be reset
		return err
	}

	return nil
}

func (self *codeDev) setup() error {

	// let anyone else know we're using the provider
	counter.Increment(util.AppName() + "_dev")

	// establish a local lock to ensure we're the only ones bringing up the
	// dev container. Also, we need to ensure the lock is released even in we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	if running := isDevRunning(); running == false {

		if err := self.loadBoxfile(); err != nil {
			return err
		}

		if err := self.setImage(); err != nil {
			return err
		}

		if err := self.downloadImage(); err != nil {
			return err
		}

		if err := self.reserveIP(); err != nil {
			return err
		}

		if err := self.launchContainer(); err != nil {
			return err
		}

		if err := self.runUserHook(); err != nil {
			return err
		}

	}

	return nil
}

func (self *codeDev) teardown() error {

	counter.Decrement(util.AppName() + "_dev")

	// establish a local app lock to ensure we're the only ones bringing
	// down the app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	if unused := devIsUnused(); unused == true {

		if err := self.removeContainer(); err != nil {
			return err
		}

		if err := self.releaseIP(); err != nil {
			return err
		}

	}

	return nil
}

// loadBoxfile loads the build boxfile from the database
func (self *codeDev) loadBoxfile() error {

	if err := data.Get(util.AppName() + "_meta", "build_boxfile", &self.boxfile); err != nil {
		return err
	}

	return nil
}

// setImage sets the image to use for the dev container
func (self *codeDev) setImage() error {
	boxfile := boxfile.New(self.boxfile.Data)

	self.image = boxfile.Node("build").StringValue("image")

	if self.image == "" {
		self.image = "nanobox/build:v1"
	}

	return nil
}

// downloadImage downloads the dev docker image
func (self *codeDev) downloadImage() error {
	label := "Pulling latest image " + self.image
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
func (self *codeDev) reserveIP() error {
	IP, err := ip_control.ReserveLocal()
	if err != nil {
		return err
	}

	self.localIP = IP

	return nil
}

// releaseIP releases a local IP back into the pool
func (self *codeDev) releaseIP() error {
	return ip_control.ReturnIP(self.localIP)
}

// launchContainer starts the dev container
func (self *codeDev) launchContainer() error {
	// parse the boxfile data
	boxfile := boxfile.New(self.boxfile.Data)
	appName := util.AppName()

	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("%s-dev", appName),
		Image:   self.image, // this will need to be configurable some time
		Network: "virt",
		IP:      self.localIP.String(),
		Binds: []string{
			fmt.Sprintf("/share/%s/code:/app", appName),
			fmt.Sprintf("/mnt/sda1/%s/build:/data", appName),
			fmt.Sprintf("/mnt/sda1/%s/cache:/mnt/cache", appName),
		},
	}

	for _, lib_dir := range boxfile.Node("code.build").StringSliceValue("lib_dirs") {
		path := "/mnt/sda1/%s/cache/lib_dirs/%s:/app/%s"
		config.Binds = append(config.Binds, fmt.Sprintf(path, appName, lib_dir, lib_dir))
	}

	// start container
	container, err := docker.CreateContainer(config)

	if err != nil {
		return err
	}

	self.container = container

	return nil
}

// removeContainer will lookup the dev container and remove it
func (self *codeDev) removeContainer() error {

	name := fmt.Sprintf("%s-dev", util.AppName())

	// grab the container info
	container, err := docker.GetContainer(name)

	if err != nil {
		return err
	}

	if err := docker.ContainerRemove(container.ID); err != nil {
		return err
	}

	return nil
}

// runUserHook runs the user hook in the dev container
func (self *codeDev) runUserHook() error {
	cmd := dockerexec.Command(self.container.ID, "user", util.UserPayload())
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// runConsole will establish a console within the dev container
func (self *codeDev) runConsole() error {

	config := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		Meta: map[string]string{
			"name": "dev",
			"working_dir": self.cwd(),
		},
	}

	err := processor.Run("dev_console", config)
	if err != nil {
		fmt.Println("dev_console:", err)
		return err
	}

	return nil
}

// cwd sets the cwd from the boxfile or provides a sensible default
func (self *codeDev) cwd() string {
	// parse the boxfile data
	boxfile := boxfile.New(self.boxfile.Data)

	if boxfile.Node("dev").StringValue("cwd") != "" {
		return boxfile.Node("dev").StringValue("cwd")
	}

	return "/app"
}

// devIsUnused returns true if the dev isn't being used by any other session
func devIsUnused() bool {
	count, err := counter.Get(util.AppName() + "_dev")

	if count == 0 && err == nil {
		return true
	}

	return false
}

// isDevRunning returns true if a service is already running
func isDevRunning() bool {
	name := fmt.Sprintf("%s-%s", util.AppName(), "dev")

	container, err := docker.GetContainer(name)

	// if the container doesn't exist then just return false
	if err != nil {
		return false
	}

	// return true if the container is running
	if container.State.Status == "running" {
		return true
	}

	return false
}
