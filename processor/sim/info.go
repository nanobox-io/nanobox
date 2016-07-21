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

// processSimInfo ...
type processSimInfo struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("sim_info", simInfoFn)
}

//
func simInfoFn(control processor.ProcessControl) (processor.Processor, error) {
	return processSimInfo{control}, nil
}

//
func (simInfo processSimInfo) Results() processor.ProcessControl {
	return simInfo.control
}

//
func (simInfo processSimInfo) Process() error {

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
	envVars := models.Evars{}
	data.Get(config.AppName()+"_meta", "sim_env", &envVars)
	bytes, _ := json.MarshalIndent(envVars, "", "  ")
	fmt.Printf("%s\n", bytes)

	return nil
}
