package code

import (
	"github.com/nanobox-io/nanobox/processor"
)


type codeBuild struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("service_setup", codeBuildFunc)
}

func codeBuildFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return &codeBuild{config: config}, nil
}


func (self serviceSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self *serviceSetup) Process() error {
	// clean up old build containers
	container, err := docker.GetContainer(fmt.Sprintf("%s-build", util.AppName(), ))
	if err == nil || container.ID != "" {
		docker.ContainerRemove(container.ID)
		ip_control.ReturnIP(docker.GetIP(container))
	}

	// create build container
	local_ip, err := ip_control.ReserveLocal()
	if err != nil {
		return err
	}

	config := docker.ContainerConfig{
		Name: fmt.Sprintf("%s-build", util.AppName()),
		Image: "nanobox/build", // this will need to be configurable some time
 		Network: "virt",
 		IP: local_ip.String(),
 		Binds: []string{
			fmt.Sprintf("/share/%s/code:/share/code", util.AppName()),
			fmt.Sprintf("/share/%s/engine:/share/engine", util.AppName()),
			fmt.Sprintf("/mnt/%s/build:/mnt/build", util.AppName()),
			fmt.Sprintf("/mnt/%s/deploy:/mnt/deploy", util.AppName()),
			fmt.Sprintf("/mnt/%s/cache:/mnt/cache", util.AppName()),
 		}
	}

	// run build hooks
	output, err := util.Exec(service.ID, "configure", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}	

	output, err = util.Exec(service.ID, "fetch", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	output, err = util.Exec(service.ID, "sniff", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	output, err = util.Exec(service.ID, "setup", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	output, err = util.Exec(service.ID, "boxfile", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}
	self.config.Meta["boxfile"] = output

	output, err = util.Exec(service.ID, "prepare", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	output, err = util.Exec(service.ID, "build", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	output, err = util.Exec(service.ID, "pack", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}

	output, err = util.Exec(service.ID, "publish", "{}")
	if err != nil {
		fmt.Println(output)
		return err
	}


	// shutdown build container
	docker.ContainerRemove(container.ID)
	ip_control.ReturnIP(local_ip)
}