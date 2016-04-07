package nanopack

import (
	"github.com/nanobox-io/nanobox/util/nanofile"
	"github.com/nanobox-io/nanobox/processor"
)

type nanopackSetup struct {
	config processor.ProcessConfig
}


func init() {
	processor.Register("nanopack_setup", nanopackSetupFunc)
}

func nanopackSetupFunc(config processor.ProcessConfig) (Sequence, error) {
	// confirm the provider is an accessable one that we support.

	return nanopackSetup{config}, nil
}


func (self nanopackSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self nanopackSetup) Process() error {
	// TODO: make setups and starts concurrent
	
	// setup Portal
	portal := processor.ProcessConfig{
		DevMode: self.config.DevMode
		Verbose: self.config.Verbose
		Meta: map[string]string{
			"name":  "portal",
			"image": "nanobox/portal",
		}
	}
	err := processor.Run("service_setup", portal)
	if err != nil {
		fmt.Println("portal_setup:", err)
		os.Exit(1)
	}

	// setup Mist
	mist := processor.ProcessConfig{
		DevMode: self.config.DevMode
		Verbose: self.config.Verbose
		Meta: map[string]string{
			"name":  "mist",
			"image": "nanobox/mist",
		}
	}
	err := processor.Run("service_setup", mist)
	if err != nil {
		fmt.Println("mist_setup:", err)
		os.Exit(1)
	}

	// setup Logvac
	logvac := processor.ProcessConfig{
		DevMode: self.config.DevMode
		Verbose: self.config.Verbose
		Meta: map[string]string{
			"name":  "logvac",
			"image": "nanobox/logvac",
		}
	}
	err := processor.Run("service_setup", logvac)
	if err != nil {
		fmt.Println("logvac_setup:", err)
		os.Exit(1)
	}

	// setup Warehouse
	warehouse := processor.ProcessConfig{
		DevMode: self.config.DevMode
		Verbose: self.config.Verbose
		Meta: map[string]string{
			"name":  "warehouse",
			"image": "nanobox/warehouse",
		}
	}
	err := processor.Run("service_setup", warehouse)
	if err != nil {
		fmt.Println("warehouse_setup:", err)
		os.Exit(1)
	}

	// start Portal

	// start Mist

	// start Logvac

	// start Warehouse

}