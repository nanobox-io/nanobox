package share

import (

  "github.com/nanobox-io/golang-docker-client"
  
  "github.com/nanobox-io/nanobox/processor"
  "github.com/nanobox-io/nanobox/provider"
)

// processShareInit ...
type processShareInit struct {
  control processor.ProcessControl
}

// assumes things were setup by the start
// but we still need docker credentials
// this process sets those up
func init() {
  processor.Register("share_init", shareInitFn)
}

//
func shareInitFn(control processor.ProcessControl) (processor.Processor, error) {
  // control.Meta["processShareInit-control"]

  // do some control validation check on the meta for the flags and make sure they
  // work

  return &processShareInit{control: control}, nil
}

//
func (shareSetup processShareInit) Results() processor.ProcessControl {
  return shareSetup.control
}

//
func (shareSetup *processShareInit) Process() error {
  if err := provider.DockerEnv(); err != nil {
    return err
  }

  if err := docker.Initialize("env"); err != nil {
    return err
  }
  return nil
}
