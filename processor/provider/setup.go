package provider

import (

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
)

type providerSetup struct {
	config processor.ProcessConfig
}


func providerSetupFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return providerSetup{config}, nil
}

func (self providerSetup) Results() processor.ProcessConfig {
	return self.config
}

func (self providerSetup) Process() error {
	err := provider.Create()
	if err != nil {
		return err
	}

	err = provider.Start()
	if err != nil {
		return err
	}

	err = provider.DockerEnv()
	if err != nil {
		return err
	}
	
	err = docker.Initialize("env")
	if err != nil {
		return err
	}

	// mount my folder
	if util.EngineDir() != "" {
		err = provider.AddMount(util.EngineDir(), "/share/"+util.AppName()+"/engine")
		if err != nil {
			return err
		}
	}
	return errprovider.AddMount(util.LocalDir(), "/share/"+util.AppName()+"/code")	
}