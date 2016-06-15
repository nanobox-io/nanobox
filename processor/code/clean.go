package code

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/data"
)

// processCodeClean ...
type processCodeClean struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("code_clean", codeCleanFn)
}

//
func codeCleanFn(control processor.ProcessControl) (processor.Processor, error) {
	return &processCodeClean{control: control}, nil
}

//
func (codeClean processCodeClean) Results() processor.ProcessControl {
	return codeClean.control
}

// clean all the code services after a dev deploy; unlike the service clean it
// doenst clean ones no longer in the box file but instead removes them all.
func (codeClean *processCodeClean) Process() error {

	// if background do nothing
	if processor.DefaultConfig.Background {
		return nil
	}

	// if other deploys are in progress do nothing
	count, _ := counter.Decrement(config.AppName() + "_deploy")
	if count != 0 {
		return nil
	}

	// remove all the existing code services
	keys, err := data.Keys(config.AppName())
	if err != nil {
		return err
	}

	// get all the code services and remove them
	for _, key := range keys {
		service := models.Service{}
		data.Get(config.AppName(), key, &key)
		if service.Type == "code" {
			codeClean.control.Meta["name"] = key
			if err := processor.Run("code_destroy", codeClean.control); err != nil {
				// we probably dont wnat to break the process just try the rest
				// and report the errors.
				// TODO: output
			}
		}
	}

	return nil
}
