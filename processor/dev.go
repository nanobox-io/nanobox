package processor

import (
	"os"
)

type dev struct {
	config ProcessConfig
}

func init() {
	Register("dev", devFunc)
}

func devFunc(config ProcessConfig) (Sequence, error) {
	config.Meta["dev-config"]
	// do some config validation
	// check on the meta for the flags and make sure they work

	return dev{config}, nil
}

func (self dev) Results() ProcessConfig {
	return self.config
}

func (self dev) Process() error {
	// setup the environment (boot vm)
	err := Run("provider_setup", self.config)
	if err != nil {
		fmt.Println("provider_setup:", err)
		os.Exit(1)
	}

	// start nanoagent service
	err := Run("nanoagent_setup", self.config)
	if err != nil {
		fmt.Println("nanoagent_setup:", err)
		os.Exit(1)
	}

	// build code
	err := Run("code_build", self.config)
	if err != nil {
		fmt.Println("code_build:", err)
		os.Exit(1)
	}

	// start services
	err := Run("service_setup", self.config)
	if err != nil {
		fmt.Println("service_setup:", err)
		os.Exit(1)
	}

	// start code
	err := Run("code_setup", self.config)
	if err != nil {
		fmt.Println("code_setup:", err)
		os.Exit(1)
	}

	// update nanoagent portal
	err := Run("update_portal", self.config)
	if err != nil {
		fmt.Println("update_portal:", err)
		os.Exit(1)
	}

	return nil
}