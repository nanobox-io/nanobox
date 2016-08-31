package code

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
)

type (
	deploy struct {
		LogvacHost      string         `json:"logvac_host"`
		Platform        string         `json:"platform"`
		Member          map[string]int `json:"member"`
		Component       component      `json:"component"`
		BeforeDeploy    interface{}    `json:"before_deploy,omitempty"`
		BeforeDeployAll interface{}    `json:"before_deploy_all,omitempty"`
		AfterDeploy     interface{}    `json:"after_deploy,omitempty"`
		AfterDeployAll  interface{}    `json:"after_deploy_all,omitempty"`
	}
)

// hookPayload ...
func DeployPayload(appModel *models.App, componentModel *models.Component) string {
	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))
	// build the payload
	pload := deploy{
		LogvacHost: appModel.LocalIPs["logvac"],
		Platform:   "local",
		Member:     map[string]int{"uid": 1},
		Component: component{
			Name: componentModel.Name,
			UID:  componentModel.Name,
			ID:   componentModel.ID,
		},
		BeforeDeploy:    boxfile.Node("code.deploy").Node("before_deploy").Value(componentModel.Name),
		BeforeDeployAll: boxfile.Node("code.deploy").Node("before_deploy_all").Value(componentModel.Name),
		AfterDeploy:     boxfile.Node("code.deploy").Node("after_deploy").Value(componentModel.Name),
		AfterDeployAll:  boxfile.Node("code.deploy").Node("after_deploy_all").Value(componentModel.Name),
	}

	// turn it into json
	j, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}

	return string(j)
}
