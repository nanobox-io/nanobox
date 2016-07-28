package dev

import (
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
)

// processDevDestroy ...
type processDevDestroy struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("dev_destroy", devDestroyFn)
}

//
func devDestroyFn(control processor.ProcessControl) (processor.Processor, error) {
	devDestroy := &processDevDestroy{control}
	fmt.Println("meta:", control.Meta)
	return devDestroy, devDestroy.validateMeta()
}

func (devDestroy *processDevDestroy) validateMeta() error {
	if devDestroy.control.Meta["app_name"] == "" {
		devDestroy.control.Meta["app_name"] = fmt.Sprintf("%s_dev", config.AppID())
	}

	return nil
}

//
func (devDestroy processDevDestroy) Results() processor.ProcessControl {
	return devDestroy.control
}

//
func (devDestroy processDevDestroy) Process() error {
	devDestroy.control.Env = "dev"

	return processor.Run("env_destroy", devDestroy.control)
}

