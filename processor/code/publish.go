package code

import (
	"encoding/json"
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/env"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/display"
)

// Publish ...
type Publish struct {
	Env            models.Env
	Image          string
	BuildID        string
	WarehouseURL   string
	WarehouseToken string
	PreviousBuild  string
	component      models.Component
}

//
func (publish *Publish) Run() error {

	// pull the image needed to publish the code
	if err := publish.pullImage(); err != nil {
		return err
	}

	// create build container
	localIP, err := dhcp.ReserveLocal()
	if err != nil {
		lumber.Error("code:Publish:dhcp.ReserveLocal(): %s", err.Error())
		return err
	}
	defer dhcp.ReturnIP(localIP)
	publish.component.InternalIP = localIP.String()

	// create a container
	if err := publish.createContainer(); err != nil {
		return err
	}
	defer publish.destroyContainer()

	lumber.Prefix("code:Publish")
	defer lumber.Prefix("")

	// run user hook
	// TODO: display output from hooks
	if _, err := util.Exec(publish.component.ID, "user", config.UserPayload(), nil); err != nil {
		return publish.runDebugSession(err)
	}

	if err := publish.runPublishHook(); err != nil {
		return publish.runDebugSession(err)
	}

	return nil
}

// pullImage ...
func (publish *Publish) pullImage() error {
	box := boxfile.NewFromPath(config.Boxfile())
	publish.Image = box.Node("build").StringValue("image")

	if publish.Image == "" {
		publish.Image = "nanobox/build:v1"
	}

	if !docker.ImageExists(publish.Image) {
		streamer := display.NewStreamer("info")
		dockerPercent := &display.DockerPercentDisplay{Output: streamer, Prefix: publish.Image}
		_, err := docker.ImagePull(publish.Image, dockerPercent)

		if err != nil {
			lumber.Error("code:Publish:pullImage:docker.ImagePull(%s, nil): %s", publish.Image, err.Error())
			return err
		}

	}
	return nil
}

// createContainer ...
func (publish *Publish) createContainer() error {
	appName := config.EnvID()
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox_%s_build", config.EnvID()),
		Image:   publish.Image, // this will need to be configurable some time
		Network: "virt",
		IP:      publish.component.InternalIP,
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
		lumber.Error("code:Publish:createContainer:docker.CreateContainer(%+v): %s", config, err.Error())
		return err
	}

	publish.component.ID = container.ID
	publish.component.Name = "build"

	return nil
}

// destroyContainer ...
func (publish *Publish) destroyContainer() error {
	if err := docker.ContainerRemove(publish.component.ID); err != nil {
		lumber.Error("code:Publish:destroyContainer:docker.ContainerRemove(%s): %s", publish.component.ID, err.Error())
		return err
	}
	return nil
}

// runPublishHook ...
func (publish *Publish) runPublishHook() error {
	// run build hooks
	pload := map[string]interface{}{}
	if publish.PreviousBuild != "" {
		pload["previous_build"] = publish.PreviousBuild
	}
	pload["build"] = publish.BuildID
	pload["warehouse"] = publish.WarehouseURL
	pload["warehouse_token"] = publish.WarehouseToken
	pload["boxfile"] = publish.Env.BuiltBoxfile
	b, _ := json.Marshal(pload)

	_, err := util.Exec(publish.component.ID, "publish", string(b), nil)

	return err
}

// runDebugSession drops the user in the build container to debug
func (publish *Publish) runDebugSession(err error) error {
	fmt.Println("there has been a failure during the publish")
	if registry.GetBool("debug") {
		fmt.Println(err)
		fmt.Println("we will be dropping you into the failed build container")
		fmt.Println("GOOD LUCK!")

		envConsole := env.Console{
			Component: publish.component,
		}
		err := envConsole.Run()
		if err != nil {
			fmt.Println("unable to enter console", err)
		}
	} else {
		return err
	}

	return nil
}
