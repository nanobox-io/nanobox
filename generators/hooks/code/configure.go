package code

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dns"
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
		Cwd          interface{}       `json:"cwd"`
		Stop         interface{}       `json:"stop"`
		StopTimeout  interface{}       `json:"stop_timeout"`
		StopForce    interface{}       `json:"stop_force"`
		CronJobs     interface{}       `json:"cron_jobs"`
		DnsEntries   interface{}       `json:"dns_entries"`
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
		LogvacHost: logvac.IPAddr(),
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
		Cwd:          boxfile.Node(componentModel.Name).Value("cwd"),
		Stop:         boxfile.Node(componentModel.Name).Value("stop"),
		StopTimeout:  boxfile.Node(componentModel.Name).Value("stop_timeout"),
		StopForce:    boxfile.Node(componentModel.Name).Value("stop_force"),
		CronJobs:     boxfile.Node(componentModel.Name).Value("cron"),
		DnsEntries:   dns.List(" by nanobox"),
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
		m = append(m, mount{component.IPAddr(), component.Plan.MountProtocol, boxNetworkDirs.StringSliceValue(node)})

	}

	return m
}
