package code

import (
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
)

type (

	// Configure ...
	Configure struct {
		Env models.Env
		App models.App
		Component models.Component
		BuildID string
		WarehouseURL string
		WarehouseToken string
	}

	// payload ...
	payload struct {
		LogvacHost   string            `json:"logvac_host,omitempty"`
		Config       interface{}       `json:"config,omitempty"`
		Component    component         `json:"component,omitempty"`
		Member       map[string]int    `json:"member,omitempty"`
		Mounts       []mount           `json:"mounts,omitempty"`
		WritableDirs interface{}       `json:"writable_dirs,omitempty"`
		Transform    interface{}       `json:"transform,omitempty"`
		Env          map[string]string `json:"env,omitempty"`
		LogWatches   interface{}       `json:"log_watches,omitempty"`
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

	fetchPayload struct {
		Component      component      `json:"component,omitempty"`
		LogvacHost     string         `json:"logvac_host,omitempty"`
		Member         map[string]int `json:"member,omitempty"`
		Build          string         `json:"build,omitempty"`
		Warehouse      string         `json:"warehouse,omitempty"`
		WarehouseToken string         `json:"warehouse_token,omitempty"`
	}
)

//
func (configure *Configure) Run() error {

	// quit now if the service was activated already
	if configure.Component.State == ACTIVE {
		return nil
	}

	// run fetch build command
	fetchPayload, err := configure.fetchPayload()
	if err != nil {
		return err
	}

	if _, err := util.Exec(configure.Component.ID, "fetch", fetchPayload, nil); err != nil {
		return err
	}

	// run configure command
	payload, err := configure.configurePayload()
	if err != nil {
		return err
	}

	//
	if _, err = util.Exec(configure.Component.ID, "configure", payload, nil); err != nil {
		return err
	}

	// run start command
	if _, err = util.Exec(configure.Component.ID, "start", payload, nil); err != nil {
		return err
	}

	//
	configure.Component.State = ACTIVE


	return configure.Component.Save()
}

// startPayload ...
func (configure Configure) startPayload() string {
	boxfile := boxfile.New([]byte(configure.Env.BuiltBoxfile))
	pload := payload{
		Config: boxfile.Node(configure.Component.Name).Value("config"),
		Start:  boxfile.Node(configure.Component.Name).StringValue("start"),
	}

	bytes, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

// configurePayload ...
func (configure *Configure) configurePayload() (string, error) {

	boxfile := boxfile.New([]byte(configure.Env.BuiltBoxfile))

	logvac, _ := models.FindComponentBySlug(configure.Component.AppID, "logvac")

	pload := payload{
		LogvacHost: logvac.InternalIP,
		Config:     boxfile.Node(configure.Component.Name).Value("config"),
		Component: component{
			Name: "whydoesthismatter",
			UID:  configure.Component.Name,
			ID:   configure.Component.ID,
		},
		Member:       map[string]int{"uid": 1},
		Mounts:       configure.mounts(),
		WritableDirs: boxfile.Node(configure.Component.Name).Value("writable_dirs"),
		Transform:    boxfile.Node("code.deploy").Value("transform"),
		Env:          configure.App.Evars,
		LogWatches:   boxfile.Node(configure.Component.Name).Value("log_watch"),
		Start:        boxfile.Node(configure.Component.Name).Value("start"),
	}

	bytes, err := json.Marshal(pload)
	return string(bytes), err
}

// fetch payload
func (configure *Configure) fetchPayload() (string, error) {

	logvac, _ := models.FindComponentBySlug(configure.Component.AppID, "logvac")

	pload := fetchPayload{
		LogvacHost: logvac.InternalIP,
		Component: component{
			Name: "whydoesthismatter",
			UID:  configure.Component.Name,
			ID:   configure.Component.ID,
		},
		Member:         map[string]int{"uid": 1},
		Build:          configure.BuildID,
		Warehouse:      configure.WarehouseURL,
		WarehouseToken: configure.WarehouseToken,
	}

	bytes, err := json.Marshal(pload)
	return string(bytes), err
}

// mounts ...
func (configure *Configure) mounts() []mount {
	boxfile := boxfile.New([]byte(configure.Env.BuiltBoxfile))
	boxNetworkDirs := boxfile.Node(configure.Component.Name).Node("network_dirs")

	m := []mount{}
	for _, node := range boxNetworkDirs.Nodes() {
		// i think i store these as data.name
		// cleanNode := regexp.MustCompile(`.+\.`).ReplaceAllString(node, "")
		component, err := models.FindComponentBySlug(configure.Component.AppID, node)
		if err != nil {
			// skip because of problems
			fmt.Println("cant get component:", err)
			continue
		}
		if !component.Plan.BehaviorPresent("mountable") || component.Plan.MountProtocol == "" {
			// skip because of problems
			fmt.Println("non mountable component", component.Name)
			continue
		}
		m = append(m, mount{component.InternalIP, component.Plan.MountProtocol, boxNetworkDirs.StringSliceValue(node)})

	}

	return m
}
