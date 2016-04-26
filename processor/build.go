package processor

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/util/locker"
)

type build struct {
	config ProcessConfig
}

func init() {
	Register("build", buildFunc)
}

func buildFunc(config ProcessConfig) (Processor, error) {
	return build{config}, nil
}

func (self build) Results() ProcessConfig {
	return self.config
}

func (self build) Process() error {
	locker.LocalLock()
	defer locker.LocalUnlock()
	self.config.Meta["build"] = "true"

	// setup the environment (boot vm)
	err := Run("provider_setup", self.config)
	if err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		os.Exit(1)
	}

	// build code
	err = Run("code_build", self.config)
	if err != nil {
		fmt.Println("code_build:", err)
		os.Exit(1)
	}

	return nil
}
