package provider

import (
	"fmt"

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
	control processor.ProcessControl
}

func providerSetupFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.

	return providerSetup{control}, nil
}

func (self providerSetup) Results() processor.ProcessControl {
	return self.control
}

func (self providerSetup) Process() error {
	// set the provider display level
	provider.Display(!processor.DefaultConfig.Quiet)

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

	if err := self.setupEvars(); err != nil {
		return err
	}

	return self.mountFolders()
}

// add the app name to the environment variables
func (self *providerSetup) setupEvars() error {
	// setup my data in the database
	app := models.App{}
	data.Get("apps", util.AppName(), &app)
	data.Put("apps", util.AppName(), app)
	evars := models.EnvVars{}
	err := data.Get(util.AppName()+"_meta", "env", &evars)
	if evars["APP_NAME"] == "" {
		evars["APP_NAME"] = util.AppName()
		return data.Put(util.AppName()+"_meta", "env", evars)
	}
	return err
}

// mount the folders for the app and any engine that is 
// local
func (self *providerSetup) mountFolders() error {
	// mount my folder
	if util.EngineDir() != "" {
		fmt.Println("  -> mount engine")
		err := provider.AddMount(util.EngineDir(), provider.HostShareDir()+util.AppName()+"/engine")
		if err != nil {
			lumber.Error("AddMount", err)
			return err
		}
	}
	return provider.AddMount(util.LocalDir(), provider.HostShareDir()+util.AppName()+"/code")
	
}