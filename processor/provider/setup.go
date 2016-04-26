package provider

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/locker"
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
	locker.GlobalLock()
	defer locker.GlobalUnlock()
	err := provider.Create()
	if err != nil {
		lumber.Error("Create()", err)
		return err
	}

	err = provider.Start()
	if err != nil {
		lumber.Error("Start()", err)
		return err
	}

	err = provider.DockerEnv()
	if err != nil {
		lumber.Error("DockerEnv()", err)
		return err
	}

	err = docker.Initialize("env")
	if err != nil {
		lumber.Error("docker.Initialize", err)
		return err
	}

	// setup my data in the database
	data.Put("apps", util.AppName(), models.App{})
	evars := models.EnvVars{}
	data.Get(util.AppName()+"_meta", "env", &evars)
	evars["APP_NAME"] = util.AppName()
	data.Put(util.AppName()+"_meta", "env", evars)

	// mount my folder
	if util.EngineDir() != "" {
		err = provider.AddMount(util.EngineDir(), "/share/"+util.AppName()+"/engine")
		if err != nil {
			lumber.Error("AddMount", err)
			return err
		}
	}
	return provider.AddMount(util.LocalDir(), "/share/"+util.AppName()+"/code")
}
