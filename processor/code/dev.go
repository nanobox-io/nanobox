package code

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	// "io"

	dockType "github.com/docker/engine-api/types"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/print"
)

type codeDev struct {
	control   processor.ProcessControl
	boxfile   models.Boxfile
	localIP   net.IP
	image     string
	container dockType.ContainerJSON
}

func init() {
	processor.Register("code_dev", codeDevFunc)
}

func codeDevFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return &codeDev{control: control}, nil
}

func (self codeDev) Results() processor.ProcessControl {
	return self.control
}

func (self *codeDev) Process() error {

	defer func() {
		if err := self.teardown(); err != nil {
			// this is bad...
			// we should probably print a pretty message explaining that the app
			// was left in a bad state and needs to be reset
			return
		}
	}()

	if err := self.setup(); err != nil {
		// todo: how to display this?
		return err
	}

	if err := self.printMOTD(); err != nil {
		return err
	}

	if err := self.runConsole(); err != nil {
		// todo: how to display this?
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

	if err := self.loadBoxfile(); err != nil {
		return err
	}

	if !isDevRunning() {

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

		if _, err := util.Exec(self.container.ID, "user", util.UserPayload(), processor.ExecWriter()); err != nil {
			return err
		}

		if _, err := util.Exec(self.container.ID, "dev", self.devPayload(), processor.ExecWriter()); err != nil {
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

	if devIsUnused() {

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

	if err := data.Get(util.AppName()+"_meta", "build_boxfile", &self.boxfile); err != nil {
		return err
	}

	return nil
}

// setImage sets the image to use for the dev container
func (self *codeDev) setImage() error {
	boxfile := boxfile.New(self.boxfile.Data)

	self.image = boxfile.Node("build").StringValue("image")

	if self.image == "" {
		self.image = "nanobox/dev"
	}

	return nil
}

// downloadImage downloads the dev docker image
func (self *codeDev) downloadImage() error {
	if !docker.ImageExists(self.image) {
		prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(self.control.DisplayLevel+1), self.image)
		_, err := docker.ImagePull(self.image, &print.DockerPercentDisplay{Prefix: prefix})
		if err != nil {
			return err
		}

	}
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
		Name:    fmt.Sprintf("nanobox-%s-dev", appName),
		Image:   self.image, // this will need to be configurable some time
		Network: "virt",
		IP:      self.localIP.String(),
		Binds: []string{
			fmt.Sprintf("%s%s/code:/app", provider.HostShareDir(), appName),
			fmt.Sprintf("%s%s/build:/data", provider.HostMntDir(), appName),
			fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), appName),
		},
	}

	for _, lib_dir := range boxfile.Node("code.build").StringSliceValue("lib_dirs") {
		path := fmt.Sprintf("/mnt/sda1/%s/cache/lib_dirs/%s:/app/%s", appName, lib_dir, lib_dir)
		config.Binds = append(config.Binds, path)
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

	name := fmt.Sprintf("nanobox-%s-dev", util.AppName())

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
func (self *codeDev) devPayload() string {
	rtn := map[string]interface{}{}
	envVars := models.EnvVars{}
	data.Get(util.AppName()+"_meta", "env", &envVars)
	rtn["env"] = envVars
	bytes, _ := json.Marshal(rtn)
	return string(bytes)
}

// runConsole will establish a console within the dev container
func (self *codeDev) runConsole() error {

	config := processor.ProcessControl{
		DevMode: self.control.DevMode,
		Verbose: self.control.Verbose,
		Meta: map[string]string{
			"name":        "dev",
			"working_dir": self.cwd(),
			"shell":       "zsh",
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

// printMOTD prints the motd with information for the user to connect
func (self *codeDev) printMOTD() error {
	os.Stderr.WriteString(`
                                   **
                                ********
                             ***************
                          *********************
                            *****************
                          ::    *********    ::
                             ::    ***    ::
                           ++   :::   :::   ++
                              ++   :::   ++
                                 ++   ++
                                    +
                    _  _ ____ _  _ ____ ___  ____ _  _
                    |\ | |__| |\ | |  | |__) |  |  \/
                    | \| |  | | \| |__| |__) |__| _/\_
`)

	return nil
}

// devIsUnused returns true if the dev isn't being used by any other session
func devIsUnused() bool {
	count, err := counter.Get(util.AppName() + "_dev")
	return count == 0 && err == nil
}

// isDevRunning returns true if a service is already running
func isDevRunning() bool {
	name := fmt.Sprintf("nanobox-%s-%s", util.AppName(), "dev")

	container, err := docker.GetContainer(name)

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}
