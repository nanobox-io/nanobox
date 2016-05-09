package nanopack

import (
	"fmt"
	"github.com/nanobox-io/nanobox/processor"
	"os"
)

type nanopackSetup struct {
	config processor.ProcessConfig
}

func init() {
	processor.Register("nanopack_setup", nanopackSetupFunc)
}

func nanopackSetupFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return nanopackSetup{config}, nil
}

func (self nanopackSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self nanopackSetup) Process() error {
	fmt.Println("-> Setup Nanopack services")

	// setup Portal
	portal := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		Meta: map[string]string{
			"name":  "portal",
			"image": "nanobox/portal",
		},
	}
	err := processor.Run("service_setup", portal)
	if err != nil {
		fmt.Println("portal_setup:", err)
		os.Exit(1)
	}

	// setup Mist
	mist := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		Meta: map[string]string{
			"name":  "mist",
			"image": "nanobox/mist",
		},
	}
	err = processor.Run("service_setup", mist)
	if err != nil {
		fmt.Println("mist_setup:", err)
		os.Exit(1)
	}

	// setup Logvac
	logvac := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		Meta: map[string]string{
			"name":  "logvac",
			"image": "nanobox/logvac",
		},
	}
	err = processor.Run("service_setup", logvac)
	if err != nil {
		fmt.Println("logvac_setup:", err)
		os.Exit(1)
	}

	// setup Warehouse
	hoarder := processor.ProcessConfig{
		DevMode: self.config.DevMode,
		Verbose: self.config.Verbose,
		Meta: map[string]string{
			"name":  "hoarder",
			"image": "nanobox/hoarder",
		},
	}
	err = processor.Run("service_setup", hoarder)
	if err != nil {
		fmt.Println("hoarder_setup:", err)
		os.Exit(1)
	}

	// start Portal
	err = processor.Run("service_configure", portal)
	if err != nil {
		fmt.Println("portal_start:", err)
		os.Exit(1)
	}

	// start Mist
	err = processor.Run("service_configure", mist)
	if err != nil {
		fmt.Println("mist_start:", err)
		os.Exit(1)
	}

	// start Logvac
	err = processor.Run("service_configure", logvac)
	if err != nil {
		fmt.Println("logvac_start:", err)
		os.Exit(1)
	}

	// start Warehouse
	err = processor.Run("service_configure", hoarder)
	if err != nil {
		fmt.Println("hoarder_start:", err)
		os.Exit(1)
	}

	return nil
}
