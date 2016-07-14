package code

import (
	"fmt"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/locker"
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

	// do not allow more then one process to run the
	// code sync or code clean at the same time
	locker.LocalLock()
	defer locker.LocalUnlock()

	// remove all the existing code services
	bucket := fmt.Sprintf("%s_%s", config.AppName(), codeClean.control.Env)
	keys, err := data.Keys(bucket)
	if err != nil {
		return err
	}

	// get all the code services and remove them
	for _, key := range keys {
		service := models.Service{}
		data.Get(bucket, key, &key)
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
