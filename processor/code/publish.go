package code

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/ip_control"
	"github.com/nanobox-io/nanobox/util/print"
)

type codePublish struct {
	control processor.ProcessControl
	service models.Service
	image string
}

func init() {
	processor.Register("code_publish", codePublishFunc)
}

func codePublishFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// {"build":"%s","warehouse":"%s","warehouse_token":"123","boxfile":"%s"}
	if control.Meta["build_id"] == "" ||
		control.Meta["warehouse_url"] == "" ||
		control.Meta["warehouse_token"] == "" {
		return nil, errors.New("missing build_id || warehouse_url || warehouse_token")
	}
	return &codePublish{control: control}, nil
}

func (self codePublish) Results() processor.ProcessControl {
	return self.control
}

func (self *codePublish) Process() error {

	// pull the image needed to publish the code
	if err := self.pullImage(); err != nil {
		return err
	}

	// create build container
	localIp, err := ip_control.ReserveLocal()
	if err != nil {
		return err
	}
	// return ip
	defer ip_control.ReturnIP(localIp)

	self.service.InternalIP = localIp.String()

	if err := self.createContainer(); err != nil {
		return err		
	}
	// shutdown container
	defer self.destroyContainer()

	if err := self.runBoxfileHook(); err != nil {
		return self.runDebugSession(err)
	}

	if err := self.runPublishHook(); err != nil {
		
	}
	if err != nil {
		return self.runDebugSession(err)
	}

	return nil
}

func (self *codePublish) pullImage() error {
	box := boxfile.NewFromPath(util.BoxfileLocation())
	self.image = box.Node("build").StringValue("image")

	if self.image == "" {
		self.image = "nanobox/build:v1"
	}

	if !docker.ImageExists(self.image) {
		prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(self.control.DisplayLevel+1), self.image)
		_, err := docker.ImagePull(self.image, &print.DockerPercentDisplay{Prefix: prefix})
		if err != nil {
			return err
		}

	}
	return nil	
}

func (self *codePublish) createContainer() error {
	appName := util.AppName()
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox-%s-build", util.AppName()),
		Image:   self.image, // this will need to be configurable some time
		Network: "virt",
		IP:      self.service.InternalIP,
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
	self.service.ID = container.ID
	self.service.Name = "build"
	return nil	
}

func (self *codePublish) destroyContainer() error {
	return docker.ContainerRemove(self.service.ID)
}

// runBoxfileHook runs the boxfile hook in the build container
func (self *codePublish) runBoxfileHook() error {
	output, err := util.Exec(self.service.ID, "boxfile", "{}", processor.ExecWriter())
	// set the boxfile in the meta
	self.control.Meta["boxfile"] = output
	self.control.Trace(stylish.Bullet("published boxfile:\n%s", output))
	return err
}

func (self *codePublish) runPublishHook() error {
	// run build hooks
	pload := map[string]interface{}{}
	pload["build"] = self.control.Meta["build_id"]
	pload["warehouse"] = self.control.Meta["warehouse_url"]
	pload["warehouse_token"] = self.control.Meta["warehouse_token"]
	pload["boxfile"] = self.control.Meta["boxfile"]
	b, _ := json.Marshal(pload)
	_, err := util.Exec(self.service.ID, "publish", string(b), processor.ExecWriter())
	return err
}

// runDebugSession drops the user in the build container to debug
func (self *codePublish) runDebugSession(err error) error {
	fmt.Println("there has been a failure during the publish")
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

