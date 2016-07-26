package dev

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/jcelliott/lumber"
	dockType "github.com/docker/engine-api/types"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/print"
)

// processDevConsole ...
type processDevConsole struct {
	control   processor.ProcessControl
	app       models.App
	boxfile   models.Boxfile
	localIP   net.IP
	image     string
	container dockType.ContainerJSON
}

//
func init() {
	processor.Register("dev_console", devConsoleFn)
}

//
func devConsoleFn(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	return &processDevConsole{control: control}, nil
}

//
func (devConsole processDevConsole) Results() processor.ProcessControl {
	return devConsole.control
}

//
func (devConsole *processDevConsole) Process() error {

	// this is bad... we should probably print a pretty message explaining that the
	// app was left in a bad state and needs to be reset
	defer func() {
		if err := devConsole.teardown(); err != nil {
			return
		}
	}()

	// run the share init which gives access to docker
	if err := processor.Run("env_init", devConsole.control); err != nil {
		return err
	}

	//
	if err := devConsole.loadApp(); err != nil {
		return err
	}

	//
	if err := devConsole.setup(); err != nil {
		return err
	}

	//
	if err := devConsole.printMOTD(); err != nil {
		return err
	}

	//
	if err := devConsole.runConsole(); err != nil {
		return err
	}

	return nil
}

// loadApp loads the app from the db
func (devConsole *processDevConsole) loadApp() error {

	key := fmt.Sprintf("%s_%s", config.AppID(), devConsole.control.Env)
	if err := data.Get("apps", key, &devConsole.app); err != nil {
		return err
	}

	return nil
}

// setup ...
func (devConsole *processDevConsole) setup() error {

	// establish a local lock to ensure we're the only ones bringing up the
	// dev container. Also, we need to ensure the lock is released even in we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// let anyone else know we're using the dev container
	counter.Increment(config.AppID() + "_dev")

	//
	if err := devConsole.loadBoxfile(); err != nil {
		devConsole.control.Display(stylish.Error("Build Boxfile", "I was unable to load a build boxfile, Did you run a build before console?"))
		return err
	}

	//
	if !isDevRunning() {

		//
		if err := devConsole.setImage(); err != nil {
			return err
		}

		//
		if err := devConsole.downloadImage(); err != nil {
			return err
		}

		//
		if err := devConsole.reserveIP(); err != nil {
			return err
		}

		//
		if err := devConsole.launchContainer(); err != nil {
			return err
		}

		//
		if err := devConsole.attachNetwork(); err != nil {
			return err
		}

		//
		if _, err := util.Exec(devConsole.container.ID, "user", config.UserPayload(), processor.ExecWriter()); err != nil {
			return err
		}

		//
		if _, err := util.Exec(devConsole.container.ID, "dev", devConsole.devPayload(), processor.ExecWriter()); err != nil {
			return err
		}
	}

	return nil
}

// teardown ...
func (devConsole *processDevConsole) teardown() error {

	// if you want to debug a problem do not teardown the container
	if devConsole.control.Debug {
		return nil
	}

	// establish a local app lock to ensure we're the only ones bringing
	// down the app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	counter.Decrement(config.AppID() + "_dev")

	//
	if devIsUnused() {

		//
		if err := devConsole.removeContainer(); err != nil {
			lumber.Error("An error occurred durring dev console teadown:%s", err.Error())
		}

		//
		if err := devConsole.detachNetwork(); err != nil {
			lumber.Error("An error occurred durring dev console teadown:%s", err.Error())
		}

		//
		if err := devConsole.releaseIP(); err != nil {
			lumber.Error("An error occurred durring dev console teadown:%s", err.Error())
		}
	}

	return nil
}

// loadBoxfile loads the build boxfile from the database
func (devConsole *processDevConsole) loadBoxfile() error {

	if err := data.Get(config.AppID()+"_meta", "build_boxfile", &devConsole.boxfile); err != nil {
		return err
	}

	return nil
}

// setImage sets the image to use for the dev container
func (devConsole *processDevConsole) setImage() error {
	boxfile := boxfile.New(devConsole.boxfile.Data)

	devConsole.image = boxfile.Node("build").StringValue("image")

	if devConsole.image == "" {
		devConsole.image = "nanobox/dev"
	}

	return nil
}

// downloadImage downloads the dev docker image
func (devConsole *processDevConsole) downloadImage() error {
	if !docker.ImageExists(devConsole.image) {
		prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(devConsole.control.DisplayLevel+1), devConsole.image)
		if _, err := docker.ImagePull(devConsole.image, &print.DockerPercentDisplay{Prefix: prefix}); err != nil {
			return err
		}
	}

	return nil
}

// reserveIP reserves a local IP for the build container
func (devConsole *processDevConsole) reserveIP() error {
	IP, err := dhcp.ReserveLocal()
	if err != nil {
		return err
	}

	devConsole.localIP = IP

	return nil
}

// releaseIP releases a local IP back into the pool
func (devConsole *processDevConsole) releaseIP() error {
	return dhcp.ReturnIP(devConsole.localIP)
}

// launchContainer starts the dev container
func (devConsole *processDevConsole) launchContainer() error {
	// parse the boxfile data
	boxfile := boxfile.New(devConsole.boxfile.Data)
	appName := config.AppID()

	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox_%s_dev", appName),
		Image:   devConsole.image, // this will need to be configurable some time
		Network: "virt",
		IP:      devConsole.localIP.String(),
		Binds: []string{
			fmt.Sprintf("%s%s/code:/app", provider.HostShareDir(), appName),
			fmt.Sprintf("%s%s/build:/data", provider.HostMntDir(), appName),
			fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), appName),
		},
	}

	//
	for _, libDir := range boxfile.Node("code.build").StringSliceValue("lib_dirs") {
		path := fmt.Sprintf("/mnt/sda1/%s/cache/lib_dirs/%s:/app/%s", appName, libDir, libDir)
		config.Binds = append(config.Binds, path)
	}

	// start container
	container, err := docker.CreateContainer(config)

	if err != nil {
		return err
	}

	devConsole.container = container

	return nil
}

// removeContainer will lookup the dev container and remove it
func (devConsole *processDevConsole) removeContainer() error {

	name := fmt.Sprintf("nanobox_%s_dev", config.AppID())

	// grab the container info
	container, err := docker.GetContainer(name)
	if err != nil {
		// if we cant get the container it may have been removed by someone else
		// just return here
		return nil
	}

	if err := docker.ContainerRemove(container.ID); err != nil {
		// but if the container exists and we cant remove it for some other reason
		// we need to report that error
		return err
	}

	return nil
}

// runUserHook runs the user hook in the dev container
func (devConsole *processDevConsole) devPayload() string {
	rtn := map[string]interface{}{}
	envVars := models.Evars{}
	data.Get(config.AppID()+"_meta", "dev_env", &envVars)
	rtn["env"] = envVars
	bytes, _ := json.Marshal(rtn)
	return string(bytes)
}

// runConsole will establish a console within the dev container
func (devConsole *processDevConsole) runConsole() error {

	config := processor.ProcessControl{
		Env:     devConsole.control.Env,
		Verbose: devConsole.control.Verbose,
		Meta: map[string]string{
			"container": fmt.Sprintf("nanobox_%s_dev", config.AppID()),
			"cwd":       devConsole.cwd(),
			"shell":     "zsh",
		},
	}

	if err := processor.Run("env_console", config); err != nil {
		return err
	}

	return nil
}

// cwd sets the cwd from the boxfile or provides a sensible default
func (devConsole *processDevConsole) cwd() string {
	// parse the boxfile data
	boxfile := boxfile.New(devConsole.boxfile.Data)

	if boxfile.Node("dev").StringValue("cwd") != "" {
		return boxfile.Node("dev").StringValue("cwd")
	}

	return "/app"
}

// printMOTD prints the motd with information for the user to connect
func (devConsole *processDevConsole) printMOTD() error {

	// fetch the dev ip
	devIP := devConsole.app.GlobalIPs["env"]

	os.Stderr.WriteString(fmt.Sprintf(`
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

--------------------------------------------------------------------------------
+ You are in a virtual machine (vm)
+ Your local source code has been mounted into the vm
+ Changes to your code in either the vm or workstation will be mirrored
+ If you run a server, access it at >> %s
--------------------------------------------------------------------------------

`, devIP))

	return nil
}

// attachNetwork attaches the container to the host network
func (devConsole *processDevConsole) attachNetwork() error {

	// fetch the devIP
	devIP := devConsole.app.GlobalIPs["env"]

	//
	if err := provider.AddIP(devIP); err != nil {
		return fmt.Errorf("provider:add_ip: %s", err.Error())
	}

	//
	if err := provider.AddNat(devIP, devConsole.localIP.String()); err != nil {
		return fmt.Errorf("provider:add_nat: %s", err.Error())
	}

	return nil
}

// detachNetwork detaches the container from the host network
func (devConsole *processDevConsole) detachNetwork() error {

	// fetch the devIP
	devIP := devConsole.app.GlobalIPs["env"]

	//
	if err := provider.RemoveNat(devIP, devConsole.localIP.String()); err != nil {
		return err
	}

	//
	if err := provider.RemoveIP(devIP); err != nil {
		return err
	}

	return nil
}

// devIsUnused returns true if the dev isn't being used by any other session
func devIsUnused() bool {
	count, err := counter.Get(config.AppID() + "_dev")
	return count == 0 && err == nil
}

// isDevRunning returns true if a service is already running
func isDevRunning() bool {
	name := fmt.Sprintf("nanobox_%s_dev", config.AppID())

	container, err := docker.GetContainer(name)

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}
