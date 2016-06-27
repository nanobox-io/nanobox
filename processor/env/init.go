package env

import (

  "github.com/nanobox-io/golang-docker-client"
  
  "github.com/nanobox-io/nanobox/processor"
  "github.com/nanobox-io/nanobox/provider"
)

// processEnvInit ...
type processEnvInit struct {
  control processor.ProcessControl
}

// assumes things were setup by the start
// but we still need docker credentials
// this process sets those up
func init() {
  processor.Register("env_init", envInitFn)
}

//
func envInitFn(control processor.ProcessControl) (processor.Processor, error) {
  // control.Meta["processEnvInit-control"]

  // do some control validation check on the meta for the flags and make sure they
  // work

  return &processEnvInit{control: control}, nil
}

//
func (envSetup processEnvInit) Results() processor.ProcessControl {
  return envSetup.control
}

//
func (envSetup *processEnvInit) Process() error {
  if err := provider.DockerEnv(); err != nil {
    return err
  }

  if err := docker.Initialize("env"); err != nil {
    return err
  }
  return nil
}
