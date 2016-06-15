package dev

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processDevInfo ...
type processDevInfo struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("dev_info", devInfoFunc)
}

//
func devInfoFunc(control processor.ProcessControl) (processor.Processor, error) {
	return processDevInfo{control}, nil
}

//
func (devInfo processDevInfo) Results() processor.ProcessControl {
	return devInfo.control
}

//
func (devInfo processDevInfo) Process() error {

	//
	services, err := data.Keys(config.AppName())
	if err != nil {
		fmt.Println("data keys:", err)
		lumber.Close()
		os.Exit(1)
	}

	//
	for _, service := range services {
		if service != "builds" {
			svc := models.Service{}
			data.Get(config.AppName(), service, &svc)
			bytes, _ := json.MarshalIndent(svc, "", "  ")
			fmt.Printf("%s\n", bytes)
		}
	}

	//
	envVars := models.EnvVars{}
	data.Get(config.AppName()+"_meta", "env", &envVars)
	bytes, _ := json.MarshalIndent(envVars, "", "  ")
	fmt.Printf("%s\n", bytes)

	return nil
}
