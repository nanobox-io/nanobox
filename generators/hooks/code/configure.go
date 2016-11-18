package code

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
)

type (
	// configure ...
	configure struct {
		LogvacHost   string            `json:"logvac_host"`
		Config       interface{}       `json:"config"`
		Component    component         `json:"component"`
		Member       map[string]int    `json:"member"`
		Mounts       []mount           `json:"mounts"`
		WritableDirs interface{}       `json:"writable_dirs"`
		Transform    interface{}       `json:"transform"`
		Env          map[string]string `json:"env"`
		LogWatches   interface{}       `json:"log_watches"`
		Start        interface{}       `json:"start"`
	}

	// component ...
	component struct {
		Name string `json:"name"`
		UID  string `json:"uid"`
		ID   string `json:"id"`
	}

	// mount ...
	mount struct {
		Host     string   `json:"host"`
		Protocol string   `json:"protocol"`
		Shares   []string `json:"shares"`
	}
)

// configurePayload ...
func ConfigurePayload(appModel *models.App, componentModel *models.Component) string {

	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))

	logvac, _ := models.FindComponentBySlug(componentModel.AppID, "logvac")

	pload := configure{
		LogvacHost: logvac.InternalIP,
		Config:     boxfile.Node(componentModel.Name).Value("config"),
		Component: component{
			Name: componentModel.Name,
			UID:  componentModel.Name,
			ID:   componentModel.ID,
		},
		Member:       map[string]int{"uid": 1},
		Mounts:       mounts(appModel, componentModel),
		WritableDirs: boxfile.Node(componentModel.Name).Value("writable_dirs"),
		Transform:    boxfile.Node("deploy.config").Value("transform"),
		Env:          appModel.Evars,
		LogWatches:   boxfile.Node(componentModel.Name).Value("log_watch"),
		Start:        boxfile.Node(componentModel.Name).Value("start"),
	}

	bytes, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

// mounts ...
func mounts(appModel *models.App, componentModel *models.Component) []mount {
	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))
	boxNetworkDirs := boxfile.Node(componentModel.Name).Node("network_dirs")

	m := []mount{}
	for _, node := range boxNetworkDirs.Nodes() {
		// i think i store these as data.name
		// cleanNode := regexp.MustCompile(`.+\.`).ReplaceAllString(node, "")
		component, err := models.FindComponentBySlug(componentModel.AppID, node)
		if err != nil {
			// skip because of problems
			continue
		}
		if !component.Plan.BehaviorPresent("mountable") || component.Plan.MountProtocol == "" {
			// skip because of problems
			continue
		}
		m = append(m, mount{component.InternalIP, component.Plan.MountProtocol, boxNetworkDirs.StringSliceValue(node)})

	}

	return m
}
