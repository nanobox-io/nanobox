package processor

import (
	"fmt"
	"os"

	"github.com/jcelliott/lumber"
	// "github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/production_api"
)

type deploy struct {
	config ProcessConfig
}

func init() {
	Register("deploy", deployFunc)
}

func deployFunc(config ProcessConfig) (Processor, error) {
	// config.Meta["deploy-config"]
	// do some config validation
	// check on the meta for the flags and make sure they work

	return deploy{config}, nil
}

func (self deploy) Results() ProcessConfig {
	return self.config
}

func (self deploy) Process() error {
	// setup the environment (boot vm)
	err := Run("provider_setup", self.config)
	if err != nil {
		fmt.Println("provider_setup:", err)
		lumber.Close()
		os.Exit(1)
	}

	// get remote hoarder credentials
	self.config.Meta["build_id"] = util.RandomString(30)
	self.config.Meta["warehouse_url"] = "??"
	self.config.Meta["warehouse_token"] = "??"

	// publish code
	publishProcessor, err := Build("code_publish", self.config)
	if err != nil {
		fmt.Println("code_publish:", err)
		os.Exit(1)
	}
	err = publishProcessor.Process()
	if err != nil {
		fmt.Println("code_publish:", err)
		os.Exit(1)
	}
	publishResult := publishProcessor.Results()
	if publishResult.Meta["boxfile"] == "" {
		fmt.Println("boxfile is empty!")
		os.Exit(1)
	}
	// boxfile := boxfile.New([]byte(publishResult.Meta["boxfile"]))
	self.config.Meta["boxfile"] = publishResult.Meta["boxfile"]

	// get the appid for the deploying app
	link := models.AppLinks{}
	data.Get(util.AppName(), "links", &link)

	// tell odin what happened
	return production_api.Deploy(link["default"], self.config.Meta["build_id"], self.config.Meta["boxfile"], self.config.Meta["message"])
}
