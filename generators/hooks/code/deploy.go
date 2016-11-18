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
		BeforeLive    interface{}    `json:"before_live,omitempty"`
		BeforeLiveAll interface{}    `json:"before_live_all,omitempty"`
		AfterLive     interface{}    `json:"after_live,omitempty"`
		AfterLiveAll  interface{}    `json:"after_live_all,omitempty"`
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
		BeforeLive:    boxfile.Node("deploy.config").Node("before_live").Value(componentModel.Name),
		BeforeLiveAll: boxfile.Node("deploy.config").Node("before_live_all").Value(componentModel.Name),
		AfterLive:     boxfile.Node("deploy.config").Node("after_live").Value(componentModel.Name),
		AfterLiveAll:  boxfile.Node("deploy.config").Node("after_live_all").Value(componentModel.Name),
	}

	// turn it into json
	j, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}

	return string(j)
}
