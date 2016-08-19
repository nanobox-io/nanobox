package sim

import (
	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/util/display"
)

// Up ...
type Up struct {
}

//
func (up Up) Run() error {

	display.OpenContext("start provider")
	// run a nanobox start
	processorStart := processors.Start{}
	if err := processorStart.Run(); err != nil {
		lumber.Error("simUp:processorStart: %s", err)
		return err
	}
	display.CloseContext()

	display.OpenContext("running build")
	// run a nanobox build
	processorBuild := &processors.Build{}
	if err := processorBuild.Run(); err != nil {
		lumber.Error("simUp:processorBuild: %s", err)
		return err
	}
	display.CloseContext()

	display.OpenContext("starting sim")
	// run a sim start
	simStart := &Start{Env: processorBuild.Env}
	if err := simStart.Run(); err != nil {
		lumber.Error("simUp:simStart: %s", err)
		return err
	}
	display.CloseContext()

	display.OpenContext("deploying sim")
	// run a sim deploy
	simDeploy := Deploy{Env: processorBuild.Env, App: simStart.App}
	if err := simDeploy.Run(); err != nil {
		lumber.Error("simUp:simDeploy: %s", err)
		return err
	}
	display.CloseContext()

	return nil
}
