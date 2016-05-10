package code

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/ip_control"
	"github.com/nanobox-io/nanobox/util/print"
)

type codePublish struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("code_publish", codePublishFunc)
}

func codePublishFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// {"build":"%s","warehouse":"%s","warehouse_token":"123","boxfile":"%s"}
	if config.Meta["build_id"] == "" ||
		config.Meta["warehouse_url"] == "" ||
		config.Meta["warehouse_token"] == "" {
		return nil, errors.New("missing build_id || warehouse_url || warehouse_token")
	}
	return &codePublish{config: config}, nil
}

func (self codePublish) Results() processor.ProcessConfig {
	return self.config
}

func (self *codePublish) Process() error {
	box := boxfile.NewFromPath(util.BoxfileLocation())
	image := box.Node("build").StringValue("image")

	if image == "" {
		image = "nanobox/build:v1"
	}

	if !docker.ImageExists(image) {
		_, err := docker.ImagePull(image, &print.DockerImageDisplaySimple{Prefix: "downloading "+image})
		if err != nil {
			return err
		}

	}

	// create build container
	localIp, err := ip_control.ReserveLocal()
	if err != nil {
		return err
	}

	// return ip
	defer ip_control.ReturnIP(localIp)
	appName := util.AppName()
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("%s-build", util.AppName()),
		Image:   image, // this will need to be configurable some time
		Network: "virt",
		IP:      localIp.String(),
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
		lumber.Error("container: ", err)
		return err
	}

	// shutdown container
	defer docker.ContainerRemove(container.ID)

	hoarder := models.Service{}
	data.Get(util.AppName(), "hoarder", &hoarder)
	pload := map[string]interface{}{}
	// we need to run the boxfile hook so the system
	// can recognize get the new services

	var b []byte
	output, err := util.Exec(container.ID, "boxfile", "{}")
	if err != nil {
		fmt.Println("output:", output)
		goto FAILURE
	}
	fmt.Println("PUBLISHED boxfile:", output)
	self.config.Meta["boxfile"] = output

	// run build hooks
	pload["build"] = self.config.Meta["build_id"]
	pload["warehouse"] = self.config.Meta["warehouse_url"]
	pload["warehouse_token"] = self.config.Meta["warehouse_token"]
	pload["boxfile"] = output
	b, err = json.Marshal(pload)
	output, err = util.Exec(container.ID, "publish", string(b))
	if err != nil {
		fmt.Println("output:", output)
		goto FAILURE
	}

	return nil
FAILURE:
	// a failure has happend and we are going to jump into the console
	fmt.Println("there has been a failure")
	fmt.Println("err:", err)
	if self.config.Verbose {
		fmt.Println("we will be dropping you into the failed build container")
		fmt.Println("GOOD LUCK!")
		self.config.Meta["name"] = "build"
		err := processor.Run("dev_console", self.config)
		if err != nil {
			fmt.Println("unable to enter console", err)
		}
	}
	return err
}
