package code

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/print"
)

// processCodePublish ...
type processCodePublish struct {
	control processor.ProcessControl
	service models.Service
	image   string
}

//
func init() {
	processor.Register("code_publish", codePublishFn)
}

//
func codePublishFn(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// {BUILD:"%s","warehouse":"%s","warehouse_token":"123","boxfile":"%s"}
	if control.Meta["build_id"] == "" ||
		control.Meta["warehouse_url"] == "" ||
		control.Meta["warehouse_token"] == "" {
		return nil, errors.New("missing build_id || warehouse_url || warehouse_token")
	}

	return &processCodePublish{control: control}, nil
}

//
func (codePublish processCodePublish) Results() processor.ProcessControl {
	return codePublish.control
}

//
func (codePublish *processCodePublish) Process() error {

	// pull the image needed to publish the code
	if err := codePublish.pullImage(); err != nil {
		return err
	}

	// create build container
	localIP, err := dhcp.ReserveLocal()
	if err != nil {
		return err
	}
	defer dhcp.ReturnIP(localIP)

	codePublish.service.InternalIP = localIP.String()

	if err := codePublish.createContainer(); err != nil {
		return err
	}
	// shutdown container
	defer codePublish.destroyContainer()

	if err := codePublish.runBoxfileHook(); err != nil {
		return codePublish.runDebugSession(err)
	}

	if err := codePublish.runPublishHook(); err != nil {

	}
	if err != nil {
		return codePublish.runDebugSession(err)
	}

	return nil
}

// pullImage ...
func (codePublish *processCodePublish) pullImage() error {
	box := boxfile.NewFromPath(config.Boxfile())
	codePublish.image = box.Node(BUILD).StringValue("image")

	if codePublish.image == "" {
		codePublish.image = "nanobox/build:v1"
	}

	if !docker.ImageExists(codePublish.image) {
		prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(codePublish.control.DisplayLevel+1), codePublish.image)
		_, err := docker.ImagePull(codePublish.image, &print.DockerPercentDisplay{Prefix: prefix})
		if err != nil {
			return err
		}

	}
	return nil
}

// createContainer ...
func (codePublish *processCodePublish) createContainer() error {
	appName := config.AppName()
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox_%s_build", config.AppName()),
		Image:   codePublish.image, // this will need to be configurable some time
		Network: "virt",
		IP:      codePublish.service.InternalIP,
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
		lumber.Error("container: ", err)
		return err
	}
	codePublish.service.ID = container.ID
	codePublish.service.Name = BUILD

	return nil
}

// destroyContainer ...
func (codePublish *processCodePublish) destroyContainer() error {
	return docker.ContainerRemove(codePublish.service.ID)
}

// runBoxfileHook runs the boxfile hook in the build container
func (codePublish *processCodePublish) runBoxfileHook() error {
	output, err := util.Exec(codePublish.service.ID, "boxfile", "{}", processor.ExecWriter())
	// set the boxfile in the meta
	codePublish.control.Meta["boxfile"] = output
	codePublish.control.Trace(stylish.Bullet("published boxfile:\n%s", output))
	return err
}

// runPublishHook ...
func (codePublish *processCodePublish) runPublishHook() error {
	// run build hooks
	pload := map[string]interface{}{}
	pload[BUILD] = codePublish.control.Meta["build_id"]
	pload["warehouse"] = codePublish.control.Meta["warehouse_url"]
	pload["warehouse_token"] = codePublish.control.Meta["warehouse_token"]
	pload["boxfile"] = codePublish.control.Meta["boxfile"]
	b, _ := json.Marshal(pload)
	_, err := util.Exec(codePublish.service.ID, "publish", string(b), processor.ExecWriter())

	return err
}

// runDebugSession drops the user in the build container to debug
func (codePublish *processCodePublish) runDebugSession(err error) error {
	fmt.Println("there has been a failure during the publish")
	if codePublish.control.Verbose {
		fmt.Println(err)
		fmt.Println("we will be dropping you into the failed build container")
		fmt.Println("GOOD LUCK!")
		codePublish.control.Meta["name"] = BUILD
		err := processor.Run("dev_console", codePublish.control)
		if err != nil {
			fmt.Println("unable to enter console", err)
		}
	} else {
		return err
	}

	return nil
}
