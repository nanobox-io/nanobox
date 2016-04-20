package code

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/ip_control"
)

type codePublish struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("code_publish", codePublishFunc)
}

func codePublishFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return &codePublish{config: config}, nil
}

func (self codePublish) Results() processor.ProcessConfig {
	return self.config
}

func (self *codePublish) Process() error {
	box := boxfile.NewFromPath(util.BoxfileLocation())
	image := box.Node("build").StringValue("image")

	if image == "" {
		image = "nanobox/build"
	}

	_, err := docker.ImagePull(image)
	if err != nil {
		return err
	}

	// create build container
	localIp, err := ip_control.ReserveLocal()
	if err != nil {
		return err
	}

	// return ip
	defer ip_control.ReturnIP(localIp)

	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("%s-build", util.AppName()),
		Image:   "nanobox/build", // this will need to be configurable some time
		Network: "virt",
		IP:      localIp.String(),
		Binds: []string{
			fmt.Sprintf("/share/%s/code:/share/code", util.AppName()),
			fmt.Sprintf("/share/%s/engine:/share/engine", util.AppName()),
			fmt.Sprintf("/mnt/%s/build:/mnt/build", util.AppName()),
			fmt.Sprintf("/mnt/%s/deploy:/mnt/deploy", util.AppName()),
			fmt.Sprintf("/mnt/%s/cache:/mnt/cache", util.AppName()),
		},
	}

	// start container
	container, err := docker.CreateContainer(config)
	if err != nil {
		lumber.Error("container: ", err)
		return err
	}

	// shutdown container
	defer docker.ContainerRemove(container.ID)

	// we need to run the boxfile hook so the system
	// can recognize get the new services
	output, err := util.Exec(container.ID, "boxfile", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}
	self.config.Meta["boxfile"] = output

	hoarder := models.Service{}
	data.Get(util.AppName(), "hoarder", &hoarder)

	// make a random build id string
	// TODO: make these parts either send local or send remote
	// if remote generate a better build id
	self.config.Meta["buildID"] = "1234"
	// create payload
	payload := fmt.Sprintf(`{"build":"%s","warehouse":"%s","warehouse_token":"123","boxfile":"%s"}`, self.config.Meta["buildID"], hoarder.InternalIP, output)

	// run build hooks
	output, err = util.Exec(container.ID, "publish", payload)
	if err != nil {
		fmt.Println(output)
		return err
	}

	return nil

}	
