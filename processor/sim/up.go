package sim

import (
	"github.com/nanobox-io/nanobox/processor"
)

// Up ...
type Up struct {
}

//
func (up Up) Run() error {
	// run a nanobox start
	processorStart := processor.Start{}
	if err := processorStart.Run(); err != nil {
		return err
	}

	// run a nanobox build
	processorBuild := processor.Build{}
	if err := processorBuild.Run(); err != nil {
		return err
	}

	// run a sim start
	simStart := Start{Env: processorBuild.Env}
	if err := simStart.Run(); err != nil {
		return err
	}

	// run a sim deploy
	simDeploy := Deploy{Env: processorBuild.Env, App: simStart.App}
	if err := simDeploy.Run(); err != nil {
		return err
	}

	return nil
}
