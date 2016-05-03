package code

import (
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/print"
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
	box := boxfile.NewFromPath(util.BoxfileLocation())
	image := box.Node("build").StringValue("image")

	if image == "" {
		image = "nanobox/build:v1"
	}

	fmt.Println("pulling image ("+image+")")
	_, err := docker.ImagePull(image, &print.DockerImageDisplaySimple{})
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

	// run user hook
	output, err := util.Exec(container.ID, "user", util.UserPayload())
	if err != nil {
		fmt.Println("user", output)
		goto FAILURE
	}

	// run build hooks
	output, err = util.Exec(container.ID, "configure", "{}")
	if err != nil {
		fmt.Println("configure", output)
		goto FAILURE
	}

	output, err = util.Exec(container.ID, "fetch", "{}")
	if err != nil {
		fmt.Println("fetch", output)
		goto FAILURE
	}

	output, err = util.Exec(container.ID, "setup", "{}")
	if err != nil {
		fmt.Println("setup", output)
		goto FAILURE
	}

	output, err = util.Exec(container.ID, "boxfile", "{}")
	if err != nil {
		fmt.Println("boxfile", output)
		goto FAILURE
	}
	fmt.Println("boxfile", output)
	self.config.Meta["boxfile"] = output

	output, err = util.Exec(container.ID, "prepare", "{}")
	if err != nil {
		fmt.Println("prepare", output)
		goto FAILURE
	}

	// conditionally build
	if self.config.Meta["build"] == "true" {
		output, err = util.Exec(container.ID, "compile", "{}")
		if err != nil {
			fmt.Println("compile", output)
			goto FAILURE
		}

		output, err = util.Exec(container.ID, "pack-app", "{}")
		if err != nil {
			fmt.Println("pack", output)
			goto FAILURE
		}

	}

	output, err = util.Exec(container.ID, "pack-deploy", "{}")
	if err != nil {
		fmt.Println("pack-live", output)
		goto FAILURE
	}

	// conditionally build
	if self.config.Meta["build"] == "true" {
		output, err = util.Exec(container.ID, "clean", "{}")
		if err != nil {
			fmt.Println("clean", output)
			goto FAILURE
		}

		output, err = util.Exec(container.ID, "pack-deploy", "{}")
		if err != nil {
			fmt.Println("pack-deploy", output)
			goto FAILURE
		}
	}

	return nil

FAILURE:
	// a failure has happend and we are going to jump into the console
	fmt.Println("there has been a failure")
	if self.config.Verbose {
		fmt.Println("we will be dropping you into the failed build container")
		fmt.Println("GOOD LUCK!")
		self.config.Meta["name"] = "build"
		err := processor.Run("console", self.config)
		if err != nil {
			fmt.Println("unable to enter console", err)
		}
	}
	return err
}
