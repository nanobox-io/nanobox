package processor

import (
	"os"
	"fmt"
	"regexp"
	
	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

)

type run struct {
	config ProcessConfig
}

func init() {
	Register("run", runFunc)
}

func runFunc(config ProcessConfig) (Processor, error) {
	// config.Meta["run-config"]
	// do some config validation
	// check on the meta for the flags and make sure they work

	return run{config}, nil
}

func (self run) Results() ProcessConfig {
	return self.config
}

func (self run) Process() error {
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

	// build code
	buildProcessor, err := Build("code_build", self.config)
	if err != nil {
		fmt.Println("code_build:", err)
		os.Exit(1)
	}
	err = buildProcessor.Process()
	if err != nil {
		fmt.Println("code_build:", err)
		os.Exit(1)
	}

	// combine the boxfiles 
	buildResult := buildProcessor.Results()
	if buildResult.Meta["boxfile"] == "" {
		fmt.Println("boxfile is empty!")
		os.Exit(1)
	}

	boxfile := boxfile.New([]byte(buildResult.Meta["boxfile"]))

	// start services
	for _, serviceName := range boxfile.Nodes("service") {
		serviceType := regexp.MustCompile(`\d+`).ReplaceAllString(serviceName, "")
		image := boxfile.Node(serviceName).StringValue("image")
		if image == "" {
			image = "nanobox/"+serviceType
		}
		service := ProcessConfig{
			DevMode: self.config.DevMode,
			Verbose: self.config.Verbose,
			Meta: map[string]string{
				"name":  serviceName,
				"image": image,
			},
		}
		err := processor.Run("service_setup", service)
		if err != nil {
			fmt.Printf("service_setup (%s): %s\n", serviceName, err.Error())
			os.Exit(1)
		}

		err := processor.Run("service_start", service)
		if err != nil {
			fmt.Printf("service_setup (%s): %s\n", serviceName, err.Error())
			os.Exit(1)
		}

	}

	// start code
	for _, codeName := range buildBox.Nodes("code") {
		image := buildBox.Node(codeName).StringValue("image")
		if image == "" {
			image = "nanobox/code"
		}
		code := ProcessConfig{
			DevMode: self.config.DevMode,
			Verbose: self.config.Verbose,
			Meta: map[string]string{
				"name":  codeName,
				"image": image,
			},
		}
		err := processor.Run("code_setup", code)
		if err != nil {
			fmt.Printf("code_setup (%s): %s\n", codeName, err.Error())
			os.Exit(1)
		}

		err := processor.Run("code_start", code)
		if err != nil {
			fmt.Printf("service_setup (%s): %s\n", codeName, err.Error())
			os.Exit(1)
		}

	}

	// update nanoagent portal
	err = Run("update_portal", self.config)
	if err != nil {
		fmt.Println("update_portal:", err)
		os.Exit(1)
	}

	return nil
}