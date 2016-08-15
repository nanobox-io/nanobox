package sim

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
)

// DeployHook ...
type (
	DeployHook struct {
		// mandatory
		App       models.App
		Component models.Component
		HookType  string
		// internal
		box boxfile.Boxfile
	}

	// member ...
	member struct {
		UID int `json:"uid,omitempty"`
	}

	// component ...
	comp struct {
		UID string `json:"uid,omitempty"`
		ID  string `json:"id,omitempty"`
	}

	// payload ...
	payload struct {
		LogvacHost      string      `json:"logvac_host"`
		Platform        string      `json:"platform"`
		Member          member      `json:"member"`
		Component       comp        `json:"component"`
		BeforeDeploy    interface{} `json:"before_deploy,omitempty"`
		BeforeDeployAll interface{} `json:"before_deploy_all,omitempty"`
		AfterDeploy     interface{} `json:"after_deploy,omitempty"`
		AfterDeployAll  interface{} `json:"after_deploy_all,omitempty"`
	}
)

//
func (deployHook *DeployHook) Run() error {
	// set the boxfile so it is easy to access
	deployHook.box = boxfile.New([]byte(deployHook.App.DeployedBoxfile))

	_, err := util.Exec(deployHook.Component.ID, deployHook.HookType, deployHook.hookPayload(), nil)

	return err
}

// hookPayload ...
func (deployHook DeployHook) hookPayload() string {
	// build the payload
	pload := payload{
		LogvacHost: deployHook.App.LocalIPs["logvac"],
		Platform:   "local",
		Member: member{
			UID: 1,
		},
		Component: comp{
			UID: deployHook.Component.Name,
			ID:  deployHook.Component.ID,
		},
		BeforeDeploy:    deployHook.box.Node("code.deploy").Node("before_deploy").Value(deployHook.Component.Name),
		BeforeDeployAll: deployHook.box.Node("code.deploy").Node("before_deploy_all").Value(deployHook.Component.Name),
		AfterDeploy:     deployHook.box.Node("code.deploy").Node("after_deploy").Value(deployHook.Component.Name),
		AfterDeployAll:  deployHook.box.Node("code.deploy").Node("after_deploy_all").Value(deployHook.Component.Name),
	}

	// turn it into json
	j, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}

	return string(j)
}
