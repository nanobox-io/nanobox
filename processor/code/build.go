package code

import (
	"fmt"
	"net"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/ip_control"
)

type codeBuild struct {
	config processor.ProcessConfig
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
	// clean up old build containers
	container, err := docker.GetContainer(fmt.Sprintf("%s-build", util.AppName()))
	if err == nil {
		docker.ContainerRemove(container.ID)
		ipString := docker.GetIP(container)
		ip_control.ReturnIP(net.ParseIP(ipString))
	}

	box := boxfile.NewFromPath(util.BoxfileLocation())
	image := box.Node("build").StringValue("image")

	if image == "" {
		image = "nanobox/build"
	}

	_, err = docker.ImagePull(image)
	if err != nil {
		return err
	}

	// create build container
	local_ip, err := ip_control.ReserveLocal()
	if err != nil {
		return err
	}

	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("%s-build", util.AppName()),
		Image:   "nanobox/build", // this will need to be configurable some time
		Network: "virt",
		IP:      local_ip.String(),
		Binds: []string{
			fmt.Sprintf("/share/%s/code:/share/code", util.AppName()),
			fmt.Sprintf("/share/%s/engine:/share/engine", util.AppName()),
			fmt.Sprintf("/mnt/%s/build:/mnt/build", util.AppName()),
			fmt.Sprintf("/mnt/%s/deploy:/mnt/deploy", util.AppName()),
			fmt.Sprintf("/mnt/%s/cache:/mnt/cache", util.AppName()),
		},
	}

	container, err = docker.CreateContainer(config)
	if err != nil {
		lumber.Error("container: ", err)
		return err
	}

	// run build hooks
	output, err := util.Exec(container.ID, "configure", "{}")
	if err != nil {
		fmt.Println("configure", output)
		return err
	}

	output, err = util.Exec(container.ID, "fetch", "{}")
	if err != nil {
		fmt.Println("build", output)
		return err
	}

	output, err = util.Exec(container.ID, "sniff", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	output, err = util.Exec(container.ID, "setup", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	output, err = util.Exec(container.ID, "boxfile", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}
	self.config.Meta["boxfile"] = output

	output, err = util.Exec(container.ID, "prepare", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	output, err = util.Exec(container.ID, "build", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	output, err = util.Exec(container.ID, "pack", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	output, err = util.Exec(container.ID, "publish", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	// shutdown build container
	docker.ContainerRemove(container.ID)
	ip_control.ReturnIP(local_ip)
	return nil
}
