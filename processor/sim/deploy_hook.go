package sim

import (
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
)

// processSimDeployHook ...
type (
	processSimDeployHook struct {
		control processor.ProcessControl
		app     models.App
		service models.Service
		box     boxfile.Boxfile
	}
	// member ...
	member struct {
		UID int `json:"uid,omitempty"`
	}

	// component ...
	component struct {
		UID string `json:"uid,omitempty"`
		ID  string `json:"id,omitempty"`
	}

	// payload ...
	payload struct {
		LogvacHost      string      `json:"logvac_host"`
		Platform        string      `json:"platform"`
		Member          member      `json:"member"`
		Component       component   `json:"component"`
		BeforeDeploy    interface{} `json:"before_deploy,omitempty"`
		BeforeDeployAll interface{} `json:"before_deploy_all,omitempty"`
		AfterDeploy     interface{} `json:"after_deploy,omitempty"`
		AfterDeployAll  interface{} `json:"after_deploy_all,omitempty"`
	}
)

//
func init() {
	processor.Register("sim_deploy_hook", simDeployHookFn)
}

//
func simDeployHookFn(control processor.ProcessControl) (processor.Processor, error) {
	simDeployHook := &processSimDeployHook{control: control}
	return simDeployHook, simDeployHook.validateMeta()
}

func (simDeployHook *processSimDeployHook) validateMeta() error {

	if simDeployHook.control.Meta["hook_type"] == "" {
		return fmt.Errorf("Deploy hooks need a hook type")
	}

	if simDeployHook.control.Meta["service_name"] == "" {
		return fmt.Errorf("Deploy hooks need a service to run on")
	}

	if simDeployHook.control.Meta["boxfile"] == "" {
		return fmt.Errorf("Deploy hooks need a boxfile")
	}

	simDeployHook.box = boxfile.New([]byte(simDeployHook.control.Meta["boxfile"]))
	if !simDeployHook.box.Valid {
		return fmt.Errorf("Given boxfile was not valid")
	}

	return nil
}

//
func (simDeployHook processSimDeployHook) Results() processor.ProcessControl {
	return simDeployHook.control
}

//
func (simDeployHook *processSimDeployHook) Process() error {

	if err := simDeployHook.loadApp(); err != nil {
		return err
	}

	if err := simDeployHook.loadService(); err != nil {
		return err
	}

	simDeployHook.control.Info(stylish.SubBullet("running %s for %s...", simDeployHook.control.Meta["hook_type"], simDeployHook.control.Meta["service_name"]))
	_, err := util.Exec(simDeployHook.service.ID, simDeployHook.control.Meta["hook_type"], simDeployHook.hookPayload(), nil)

	return err
}

// loadApp loads the app from the database
func (simDeployHook *processSimDeployHook) loadApp() error {
	key := fmt.Sprintf("%s_%s", config.AppID(), simDeployHook.control.Env)
	return data.Get("apps", key, &simDeployHook.app)
}

// loadService loads the service from the database
func (simDeployHook *processSimDeployHook) loadService() error {
	bucket := fmt.Sprintf("%s_%s", config.AppID(), simDeployHook.control.Env)
	return data.Get(bucket, simDeployHook.control.Meta["service_name"], &simDeployHook.service)
}

// hookPayload ...
func (simDeployHook processSimDeployHook) hookPayload() string {
	key := fmt.Sprintf("%s_%s", config.AppID(), simDeployHook.control.Env)

	// load the app
	app := models.App{}
	data.Get("apps", key, &app)

	// load the service
	serviceName := simDeployHook.control.Meta["service_name"]
	service := models.Service{}
	data.Get(key, serviceName, &service)

	// build the payload
	pload := payload{
		LogvacHost: app.LocalIPs["logvac"],
		Platform:   "local",
		Member: member{
			UID: 1,
		},
		Component: component{
			UID: serviceName,
			ID:  service.ID,
		},
		BeforeDeploy:    simDeployHook.box.Node("code.deploy").Node("before_deploy").Value(serviceName),
		BeforeDeployAll: simDeployHook.box.Node("code.deploy").Node("before_deploy_all").Value(serviceName),
		AfterDeploy:     simDeployHook.box.Node("code.deploy").Node("after_deploy").Value(serviceName),
		AfterDeployAll:  simDeployHook.box.Node("code.deploy").Node("after_deploy_all").Value(serviceName),
	}

	// turn it into json
	j, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}

	return string(j)
}
