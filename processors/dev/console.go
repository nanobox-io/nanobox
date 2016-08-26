package dev

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	dockType "github.com/docker/engine-api/types"
	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/locker"
	"github.com/nanobox-io/nanobox/util/provider"
)

// Console ...
type Console struct {
	App    models.App
	DevRun bool

	boxfile   boxfile.Boxfile
	localIP   net.IP
	image     string
	container dockType.ContainerJSON
}

//
func (console *Console) Run() error {

	// this is bad... we should probably print a pretty message explaining that the
	// app was left in a bad state and needs to be reset
	defer func() {
		if err := console.teardown(); err != nil {
			return
		}
	}()

	// run the share init which gives access to docker
	envInit := env.Init{}
	if err := envInit.Run(); err != nil {
		return err
	}

	//
	if err := console.setup(); err != nil {
		return err
	}

	// if run then start the run commands
	// and log but do not continue to the regular console
	if console.DevRun {
		devRun := Run{App: console.App}
		return devRun.Run()
	}

	//
	if err := console.printMOTD(); err != nil {
		return err
	}

	//
	if err := console.runConsole(); err != nil {
		return err
	}

	return nil
}

// setup ...
func (console *Console) setup() error {

	// establish a local lock to ensure we're the only ones bringing up the
	// dev container. Also, we need to ensure the lock is released even in we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	// let anyone else know we're using the dev container
	counter.Increment(console.App.ID)

	//
	if err := console.loadBoxfile(); err != nil {
		return err
	}

	//
	if !isDevRunning() {

		//
		if err := console.setImage(); err != nil {
			return err
		}

		//
		if err := console.downloadImage(); err != nil {
			return err
		}

		//
		if err := console.reserveIP(); err != nil {
			return err
		}

		//
		if err := console.launchContainer(); err != nil {
			return err
		}

		//
		if err := console.attachNetwork(); err != nil {
			return err
		}

		lumber.Prefix("dev:Console")
		defer lumber.Prefix("")

		// TODO: The nil in this call needs to be replaced with something from the new display
		if _, err := util.Exec(console.container.ID, "user", config.UserPayload(), nil); err != nil {
			return err
		}

		// TODO: The nil in this call needs to be replaced with something from the new display
		if _, err := util.Exec(console.container.ID, "dev", console.devPayload(), nil); err != nil {
			return err
		}
	}

	return nil
}

// teardown ...
func (console *Console) teardown() error {

	// establish a local app lock to ensure we're the only ones bringing
	// down the app platform. Also ensure that we release it even if we error
	locker.LocalLock()
	defer locker.LocalUnlock()

	counter.Decrement(console.App.ID)

	//
	if devIsUnused() {

		//
		if err := console.removeContainer(); err != nil {
			lumber.Error("An error occurred durring dev console teadown:%s", err.Error())
		}

		//
		if err := console.detachNetwork(); err != nil {
			lumber.Error("An error occurred durring dev console teadown:%s", err.Error())
		}

		//
		if err := console.releaseIP(); err != nil {
			lumber.Error("An error occurred durring dev console teadown:%s", err.Error())
		}
	}

	return nil
}

// loadBoxfile loads the build boxfile from the database
func (console *Console) loadBoxfile() error {

	env, _ := models.FindEnvByID(console.App.EnvID)
	console.boxfile = boxfile.New([]byte(env.BuiltBoxfile))
	if !console.boxfile.Valid {
		return fmt.Errorf("the boxfile from the build is invalid")
	}
	return nil
}

// setImage sets the image to use for the dev container
func (console *Console) setImage() error {

	console.image = console.boxfile.Node("build").StringValue("image")

	if console.image == "" {
		console.image = "nanobox/dev"
	}

	return nil
}

// downloadImage downloads the dev docker image
func (console *Console) downloadImage() error {
	if !docker.ImageExists(console.image) {
		streamer := display.NewStreamer("info")
		dockerPercent := &display.DockerPercentDisplay{Output: streamer, Prefix: console.image}
		if _, err := docker.ImagePull(console.image, dockerPercent); err != nil {
			return err
		}
	}

	return nil
}

// reserveIP reserves a local IP for the build container
func (console *Console) reserveIP() error {
	IP, err := dhcp.ReserveLocal()

	console.localIP = IP

	return err
}

// releaseIP releases a local IP back into the pool
func (console *Console) releaseIP() error {
	return dhcp.ReturnIP(console.localIP)
}

// launchContainer starts the dev container
func (console *Console) launchContainer() error {
	// parse the boxfile data

	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox_%s", console.App.ID),
		Image:   console.image, // this will need to be configurable some time
		Network: "virt",
		IP:      console.localIP.String(),
		Binds: []string{
			fmt.Sprintf("%s%s/code:/app", provider.HostShareDir(), console.App.EnvID),
			fmt.Sprintf("%s%s/build:/data", provider.HostMntDir(), console.App.EnvID),
			fmt.Sprintf("%s%s/cache:/mnt/cache", provider.HostMntDir(), console.App.EnvID),
		},
	}

	//
	for _, libDir := range console.boxfile.Node("code.build").StringSliceValue("lib_dirs") {
		path := fmt.Sprintf("/mnt/sda1/%s/cache/lib_dirs/%s:/app/%s", console.App.EnvID, libDir, libDir)
		config.Binds = append(config.Binds, path)
	}

	// start container
	container, err := docker.CreateContainer(config)

	if err != nil {
		return err
	}

	console.container = container

	return nil
}

// removeContainer will lookup the dev container and remove it
func (console *Console) removeContainer() error {

	name := fmt.Sprintf("nanobox_%s", console.App.ID)

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
func (console *Console) devPayload() string {
	rtn := map[string]interface{}{}
	rtn["env"] = console.App.Evars
	bytes, _ := json.Marshal(rtn)
	return string(bytes)
}

// runConsole will establish a console within the dev container
func (console *Console) runConsole() error {

	// create a dummy component using the appname
	component := models.Component{
		ID: "nanobox_" + console.App.ID,
	}
	// for tyler: I dont like forcing someone into zsh..
	// your chosen shell is a very personal
	envConsole := env.Console{
		Component: component,
		Cwd:       console.cwd(),
		Shell:     "zsh",
	}

	return envConsole.Run()
}

// cwd sets the cwd from the boxfile or provides a sensible default
func (console *Console) cwd() string {
	// parse the boxfile data

	if console.boxfile.Node("dev").StringValue("cwd") != "" {
		return console.boxfile.Node("dev").StringValue("cwd")
	}

	return "/app"
}

// printMOTD prints the motd with information for the user to connect
func (console *Console) printMOTD() error {

	// fetch the dev ip
	devIP := console.App.GlobalIPs["env"]

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
func (console *Console) attachNetwork() error {

	// fetch the devIP
	devIP := console.App.GlobalIPs["env"]

	//
	if err := provider.AddIP(devIP); err != nil {
		return fmt.Errorf("provider:add_ip: %s", err.Error())
	}

	//
	if err := provider.AddNat(devIP, console.localIP.String()); err != nil {
		return fmt.Errorf("provider:add_nat: %s", err.Error())
	}

	return nil
}

// detachNetwork detaches the container from the host network
func (console *Console) detachNetwork() error {

	// fetch the devIP
	devIP := console.App.GlobalIPs["env"]

	//
	if err := provider.RemoveNat(devIP, console.localIP.String()); err != nil {
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
	count, err := counter.Get(config.EnvID() + "_dev")
	return count == 0 && err == nil
}

// isDevRunning returns true if a service is already running
func isDevRunning() bool {
	name := fmt.Sprintf("nanobox_%s_dev", config.EnvID())

	container, err := docker.GetContainer(name)

	// if the container doesn't exist then just return false
	return err == nil && container.State.Status == "running"
}
