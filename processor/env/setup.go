package env

import (
	"github.com/nanobox-io/nanobox/processor"
)

// processEnvSetup ...
type processEnvSetup struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("env_setup", envSetupFn)
}

//
func envSetupFn(control processor.ProcessControl) (processor.Processor, error) {
	// control.Meta["processEnvSetup-control"]

	// do some control validation check on the meta for the flags and make sure they
	// work

	return &processEnvSetup{control: control}, nil
}

//
func (envSetup processEnvSetup) Results() processor.ProcessControl {
	return envSetup.control
}

//
func (envSetup *processEnvSetup) Process() error {

	if err := envSetup.setupProvider(); err != nil {
		return err
	}

	if err := envSetup.setupMounts(); err != nil {
		return err
	}

	// if there is an environment then we should set up app
	// if not (in the case of a build) no app setup is necessary
	if envSetup.control.Env != "" {
		if err := envSetup.setupApp(); err != nil {
			return err
		}
	}

	return nil
}

// setupProvider sets up the provider
func (envSetup *processEnvSetup) setupProvider() error {
	return processor.Run("provider_setup", envSetup.control)
}

// setupMounts will add the envs and mounts for this app
func (envSetup *processEnvSetup) setupMounts() error {
	return processor.Run("env_mount", envSetup.control)
}

// setupApp sets up the app plaftorm and data services
func (envSetup *processEnvSetup) setupApp() error {

	// setup the app
	if err := processor.Run("app_setup", envSetup.control); err != nil {
		return err
	}

	// clean up after any possible failures in a previous deploy
	if err := processor.Run("service_clean", envSetup.control); err != nil {
		return err
	}

	// setup the platform services
	return processor.Run("platform_setup", envSetup.control)
}
