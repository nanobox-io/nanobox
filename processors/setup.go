package processors

import (
	"fmt"
	"strconv"

	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/config"
)

func Setup() error {
	if config.ConfigExists() {
		return nil
	}

	setupConf := &config.SetupConf{
		Provider: "docker-machine",
		Mount:    "native",
		CPUs:     1,
		RAM:      1,
	}

	var err error
	// ask about provider
	setupConf.Provider, err = display.Ask("provider")
	if err != nil {
		return err
	}

	if setupConf.Provider != "native" && setupConf.Provider != "docker-machine" {
		return fmt.Errorf("we only support native and docker-machine providers currently")
	}


	// if provider == docker-machine ask more questions
	if setupConf.Provider == "native" {
		config.ConfigFile(setupConf)
		return nil
	}
	// ask about mount types
	setupConf.Mount, err = display.Ask("mount-type")
	if err != nil {
		return err
	}

	if setupConf.Mount != "native" && setupConf.Mount != "netfs" {
		return fmt.Errorf("we only support native and netfs networks currently")
	}

	// ask about cpus
	cpuString, err := display.Ask("CPUS")
	if err != nil {
		return err
	}

	if setupConf.CPUs, err = strconv.Atoi(cpuString); err != nil {
		return fmt.Errorf("failed to convert input to an int")
	}

	// ask about rams
	ramString, err := display.Ask("RAM")
	if err != nil {
		return err
	}

	if setupConf.RAM, err = strconv.Atoi(ramString); err != nil {
		return fmt.Errorf("failed to convert input to int")
	}

	config.ConfigFile(setupConf)
	return nil

}