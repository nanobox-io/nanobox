package code

import (
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/counter"
	"github.com/nanobox-io/nanobox/util/data"
)

// Clean all the code services after a dev deploy
// unlike the service clean it doenst clean ones no longer in the box file
// but instead removes them all.
type codeClean struct {
	control processor.ProcessControl
}

func init() {
	processor.Register("code_clean", codeCleanFunc)
}

func codeCleanFunc(control processor.ProcessControl) (processor.Processor, error) {
	return &codeClean{control: control}, nil
}

func (self codeClean) Results() processor.ProcessControl {
	return self.control
}

func (self *codeClean) Process() error {
	// if background do nothing
	if processor.DefaultConfig.Background {
		return nil
	}

	count, _ := counter.Decrement(util.AppName() + "_deploy")
	// if other deploys are in progress do nothing
	if count != 0 {
		return nil
	}

	// remove all the existing code services
	keys, err := data.Keys(util.AppName())
	if err != nil {
		return err
	}

	// get all the code services and remove them
	for _, key := range keys {
		service := models.Service{}
		data.Get(util.AppName(), key, &key)
		if service.Type == "code" {
			self.control.Meta["name"] = key
			if err := processor.Run("code_destroy", self.control); err != nil {
				// we probably dont wnat to break the process just try the rest
				// and report the errors.
				// TODO: output
			}
		}
	}
	return nil
}
