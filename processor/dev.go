package processor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

// processDev ...
type processDev struct {
	control      ProcessControl
	oldBoxfile   models.Boxfile
	newBoxfile   models.Boxfile
	buildBoxfile models.Boxfile
}

//
func init() {
	Register("dev", devFunc)
}

//
func devFunc(control ProcessControl) (Processor, error) {
	// control.Meta["dev-control"]

	// do some control validation check on the meta for the flags and make sure they
	// work

	return processDev{control: control}, nil
}

//
func (dev processDev) Results() ProcessControl {
	return dev.control
}

//
func (dev processDev) Process() error {

	// defer the clean up so if we exit early the cleanup will always happen
	defer func() {
		if err := Run("dev_teardown", dev.control); err != nil {
			fmt.Println("teardown broke")
			fmt.Println(err)

			// this is bad, really bad...
			// we should probably print a pretty message explaining that the app
			// was left in a bad state and needs to be reset
			os.Exit(1)
		}
	}()

	// get the vm and app up.
	if err := Run("dev_setup", dev.control); err != nil {
		return err
	}

	// startDataServices will start all data services
	if err := Run("service_start_all", dev.control); err != nil {
		return err
	}

	if err := dev.runBuild(); err != nil {
		return err
	}

	// starts a dev container and establishes a console session
	if err := Run("code_dev", dev.control); err != nil {
		return err
	}

	return nil
}

// runBuild ...
func (dev *processDev) runBuild() error {
	if err := dev.fetchOldBoxfile(); err != nil {
		return err
	}

	if err := dev.fetchNewBoxfile(); err != nil {
		return err
	}

	// if the build has been done or not we always have
	// to check the boxfile to determine if we are going to
	// build/rebuild or use the existing one
	if dev.hasBoxfileChanged() {
		// build the code
		if err := Run("code_build", dev.control); err != nil {
			return err
		}

		// persist the new boxfile so we know not to build next time.
		if err := dev.persistNewBoxfile(); err != nil {
			return err
		}

	}

	// syncronize the data services
	if err := Run("service_sync", dev.control); err != nil {
		return err
	}

	return nil
}

// fetchOldBoxfile fetches the old boxfile from the db
func (dev *processDev) fetchOldBoxfile() error {
	// we don't care about the error here because it's very likely
	// that there won't be an old boxfile.
	data.Get(util.AppName()+"_meta", "boxfile", &dev.oldBoxfile)

	return nil
}

// fetchNewBoxfile fetches the new boxfile
func (dev *processDev) fetchNewBoxfile() error {
	rawData, err := ioutil.ReadFile(util.BoxfileLocation())

	if err != nil {
		return errors.New("unable to load boxfile.yml")
	}

	dev.newBoxfile.Data = rawData

	return nil
}

// persistNewBoxfile persists the new boxfile to the database
func (dev *processDev) persistNewBoxfile() error {

	key := util.AppName() + "_meta"
	if err := data.Put(key, "boxfile", dev.newBoxfile); err != nil {
		return err
	}

	return nil
}

// hasBoxfileChanged returns true if the boxfile has changed
func (dev *processDev) hasBoxfileChanged() bool {

	if string(dev.oldBoxfile.Data) != string(dev.newBoxfile.Data) || len(dev.oldBoxfile.Data) == 0 {
		return true
	}

	return false
}
