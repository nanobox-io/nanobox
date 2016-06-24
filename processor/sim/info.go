package sim

import (
	"encoding/json"
	"fmt"

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
	processor.Register("sim_info", simInfoFn)
}

//
func simInfoFn(control processor.ProcessControl) (processor.Processor, error) {
	return processDevInfo{control}, nil
}

//
func (simInfo processDevInfo) Results() processor.ProcessControl {
	return simInfo.control
}

//
func (simInfo processDevInfo) Process() error {

	//
	bucket := fmt.Sprintf("%s_sim", config.AppName())
	services, err := data.Keys(bucket)
	if err != nil {
		fmt.Println("data keys:", err)
		lumber.Close()
		return err
	}

	//
	for _, service := range services {
		if service != "builds" {
			svc := models.Service{}
			data.Get(bucket, service, &svc)
			bytes, _ := json.MarshalIndent(svc, "", "  ")
			fmt.Printf("%s\n", bytes)
		}
	}

	//
	envVars := models.EnvVars{}
	data.Get(config.AppName()+"_meta", "sim_env", &envVars)
	bytes, _ := json.MarshalIndent(envVars, "", "  ")
	fmt.Printf("%s\n", bytes)

	return nil
}
