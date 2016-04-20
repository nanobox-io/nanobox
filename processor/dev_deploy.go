package processor

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"
)

type devDeploy struct {
	config ProcessConfig
}

func init() {
	Register("dev_deploy", devDeployFunc)
}

func devDeployFunc(config ProcessConfig) (Processor, error) {
	// config.Meta["devDeploy-config"]
	// do some config validation
	// check on the meta for the flags and make sure they work

	return devDeploy{config}, nil
}

func (self devDeploy) Results() ProcessConfig {
	return self.config
}

func (self devDeploy) Process() error {
	// setup the environment (boot vm)
	err := Run("provider_setup", self.config)
	if err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		os.Exit(1)
	}

	// start nanopack service
	err = Run("nanopack_setup", self.config)
	if err != nil {
		fmt.Println("nanoagent_setup:", err)
		os.Exit(1)
	}

	// publish code
	publishProcessor, err := Build("code_publish", self.config)
	if err != nil {
		fmt.Println("code_publish:", err)
		os.Exit(1)
	}
	err = publishProcessor.Process()
	if err != nil {
		fmt.Println("code_publish:", err)
		os.Exit(1)
	}
	publishResult := publishProcessor.Results()
	if publishResult.Meta["boxfile"] == "" {
		fmt.Println("boxfile is empty!")
		os.Exit(1)
	}
	boxfile := boxfile.New([]byte(publishResult.Meta["boxfile"]))
	self.config.Meta["boxfile"] = publishResult.Meta["boxfile"]

	// syncronize the services as per the new boxfile
	err = Run("service_sync", self.config)
	if err != nil {
		fmt.Println("service_sync:", err)
		lumber.Close()
		os.Exit(1)
	}	

	// start code
	for _, codeName := range boxfile.Nodes("code") {
		image := boxfile.Node(codeName).StringValue("image")
		if image == "" {
			image = "nanobox/code"
		}
		code := ProcessConfig{
			DevMode: self.config.DevMode,
			Verbose: self.config.Verbose,
			Meta: map[string]string{
				"name":  codeName,
				"image": image,
				"boxfile": self.config.Meta["boxfile"],
			},
		}
		err := Run("code_setup", code)
		if err != nil {
			fmt.Printf("code_setup (%s): %s\n", codeName, err.Error())
			os.Exit(1)
		}

		err = Run("code_start", code)
		if err != nil {
			fmt.Printf("code_start (%s): %s\n", codeName, err.Error())
			os.Exit(1)
		}

	}

	// update nanoagent portal
	err = Run("update_portal", self.config)
	if err != nil {
		fmt.Println("update_portal:", err)
		os.Exit(1)
	}

	// hang and do some logging until they are ready
	// shut down all services 
	// possibly shut down provider

	return nil
}
