package processor

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/locker"
)

type dev struct {
	config ProcessConfig
}

func init() {
	Register("dev", devFunc)
}

func devFunc(config ProcessConfig) (Processor, error) {
	// config.Meta["dev-config"]
	// do some config validation
	// check on the meta for the flags and make sure they work

	return dev{config}, nil
}

func (self dev) Results() ProcessConfig {
	return self.config
}

func (self dev) Process() error {
	// setup the environment (boot vm)
	err := Run("provider_setup", self.config)
	if err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		os.Exit(1)
	}

	// start all the services that are in standby
	err = Run("service_start_all", self.config)
	if err != nil {
		fmt.Printf("service_start_all: %s\n", err.Error())
		os.Exit(1)
	}

	// start nanopack service
	err = Run("nanopack_setup", self.config)
	if err != nil {
		fmt.Println("nanoagent_setup:", err)
		os.Exit(1)
	}

	locker.LocalLock()
	box := models.Boxfile{}
	box.Data, _ = ioutil.ReadFile(util.BoxfileLocation())

	oldBoxData := models.Boxfile{}
	data.Get(util.AppName()+"_meta", "boxfile", &oldBoxData)

	if string(oldBoxData.Data) != string(box.Data) {
		// build code (without build hook)
		buildProcessor, err := Build("code_build", self.config)
		if err != nil {
			fmt.Println("code_build:", err)
			os.Exit(1)
		}
		err = buildProcessor.Process()
		if err != nil {
			fmt.Println("code_build:", err)
			os.Exit(1)
		}

		// combine the boxfiles
		buildResult := buildProcessor.Results()
		if buildResult.Meta["boxfile"] == "" {
			fmt.Println("boxfile is empty!")
			os.Exit(1)
		}
		self.config.Meta["boxfile"] = buildResult.Meta["boxfile"]

		// syncronize the services as per the new boxfile
		err = Run("service_sync", self.config)
		if err != nil {
			fmt.Println("service_sync:", err)
			lumber.Close()
			os.Exit(1)
		}
	}

	// make sure everyone knows im using the app (so dont shut down)
	app := models.App{}
	data.Get("apps", util.AppName(), &app)
	app.UsageCount++
	data.Put("apps", util.AppName(), app)

	locker.LocalUnlock()

	// syncronize the services as per the new boxfile
	self.config.Meta["name"] = "dev"
	self.config.Meta["workding_dir"] = "/code"
	err = Run("code_dev", self.config)
	// make sure we stop let the db know we
	// are done with the app and it can be
	// shut down
	data.Get("apps", util.AppName(), &app)
	app.UsageCount--
	data.Put("apps", util.AppName(), app)

	if err != nil {
		fmt.Println("code_dev:", err)
		lumber.Close()
		os.Exit(1)
	}

	err = Run("dev_stop", self.config)
	if err != nil {
		fmt.Println("dev_stop:", err)
		lumber.Close()
		os.Exit(1)
	}

	return data.Put(util.AppName()+"_meta", "boxfile", box)
}
