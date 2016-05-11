package processor

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

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
	box.Data, err = ioutil.ReadFile(util.BoxfileLocation())
	lumber.Debug("ioutil err %+v", err)

	oldBoxData := models.Boxfile{}
	data.Get(util.AppName()+"_meta", "boxfile", &oldBoxData)

	if string(oldBoxData.Data) != string(box.Data) || len(box.Data) == 0 {
		lumber.Debug("old boxfile:(%s)\nnew boxfile:(%s)", oldBoxData.Data, box.Data)
		err = data.Put(util.AppName()+"_meta", "boxfile", box)
		if err != nil {
			fmt.Println("unable to store new boxfile:", err)
		}

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
		box.Data = []byte(buildResult.Meta["boxfile"])
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
	lumber.Debug("incrementing usagecount toto %d", app.UsageCount)
	if app.UsageCount < 0 {
		app.UsageCount = 0
	}
	app.UsageCount = app.UsageCount + 1
	lumber.Debug("incrementing usagecount to %d", app.UsageCount)
	err = data.Put("apps", util.AppName(), app)
	lumber.Error("dataputerr: %+v", err)

	appAfter := models.App{}
	data.Get("apps", util.AppName(), &appAfter)
	lumber.Debug("incrementing usagecount after %d", appAfter.UsageCount)

	locker.LocalUnlock()

	// get the working dir from the last build
	self.config.Meta["working_dir"] = "/app"

	bBox := models.Boxfile{}
	data.Get(util.AppName()+"_meta", "build_boxfile", &bBox)
	lumber.Debug("dev: buildBox: %s", bBox.Data)
	boxf := boxfile.New(bBox.Data)
	if boxf.Node("dev").StringValue("cwd") != "" {
		self.config.Meta["working_dir"] = boxf.Node("dev").StringValue("cwd")
	}

	self.config.Meta["name"] = "dev"
	err = Run("code_dev", self.config)
	// make sure we stop let the db know we
	// are done with the app and it can be
	// shut down
	lumber.Debug("decrementing usagecount fromfrom %d", app.UsageCount)
	app = models.App{}
	err = data.Get("apps", util.AppName(), &app)
	lumber.Error("errfromdata:%+v",err)
	lumber.Debug("decrementing usagecount from %d", app.UsageCount)
	app.UsageCount = app.UsageCount - 1
	if app.UsageCount < 0 {
		app.UsageCount = 0
	}
	lumber.Debug("decrementing usagecount to %d", app.UsageCount)
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

	return nil
}
