package processor

import (
	"io/ioutil"
	"os"
	"errors"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type dev struct {
	config 				ProcessConfig
	oldBoxfile		models.Boxfile
	newBoxfile		models.Boxfile
	buildBoxfile 	models.Boxfile
}

func init() {
	Register("dev", devFunc)
}

func devFunc(config ProcessConfig) (Processor, error) {
	// config.Meta["dev-config"]
	// do some config validation
	// check on the meta for the flags and make sure they work

	return dev{config: config}, nil
}

func (self dev) Results() ProcessConfig {
	return self.config
}

func (self dev) Process() error {

	// defer the clean up so if we exit early the
	// cleanup will always happen
	defer func() {
		if err := Run("dev_teardown", self.config); err != nil {
			// this is bad, really bad...
			// we should probably print a pretty message explaining that the app
			// was left in a bad state and needs to be reset
			os.Exit(1)
		}
	}()

	if err := Run("dev_setup", self.config); err != nil {
		// todo: how to display this?
		return err
	}

	if err := self.runBuild(); err != nil {
		// todo: how to display this?
		return err
	}

	// startDataServices will start all data services
	if err := Run("service_start_all", self.config); err != nil {
		// todo: how to display this?
		return err
	}

	// starts a dev container and establishes a console session
	if err := Run("code_dev", self.config); err != nil {
		// todo: how to display this?
		return err
	}

	return nil
}

func (self *dev) runBuild() error {
	if err := self.fetchOldBoxfile(); err != nil {
		return err
	}

	if err := self.fetchNewBoxfile(); err != nil {
		return err
	}

	// todo: we need to consider a more stateful way of determining if a dev
	// was successful previously. Otherwise a failure on the first run won't
	// try a subsequent build
	if self.hasBoxfileChanged() {
		// build the code
		if err := Run("code_build", self.config); err != nil {
			return err
		}

		// syncronize the data services with the new build
		if err := Run("service_sync", self.config); err != nil {
			return err
		}

		// persist the new boxfile so we know not to build next time.
		if err := self.persistNewBoxfile(); err != nil {
			return err
		}

	}
	return nil
}

// fetchOldBoxfile fetches the old boxfile from the db
func (self *dev) fetchOldBoxfile() error {
	// we don't care about the error here because it's very likely
	// that there won't be an old boxfile.
	data.Get(util.AppName()+"_meta", "boxfile", &self.oldBoxfile)

	return nil
}

// fetchNewBoxfile fetches the new boxfile
func (self *dev) fetchNewBoxfile() error {
	rawData, err := ioutil.ReadFile(util.BoxfileLocation())

	if err != nil {
		return errors.New("unable to load boxfile.yml")
	}

	self.newBoxfile.Data = rawData

	return nil
}

// persistNewBoxfile persists the new boxfile to the database
func (self *dev) persistNewBoxfile() error {

	key := util.AppName() + "_meta"
	if err := data.Put(key, "boxfile", self.newBoxfile); err != nil {
		return err
	}

	return nil
}

// hasBoxfileChanged returns true if the boxfile has changed
func (self *dev) hasBoxfileChanged() bool {

	if string(self.oldBoxfile.Data) != string(self.newBoxfile.Data) || len(self.oldBoxfile.Data) == 0 {
		return true
	}

	return false
}
