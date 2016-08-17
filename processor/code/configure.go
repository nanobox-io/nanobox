package code

import (
	"encoding/json"
	"fmt"

	"github.com/jcelliott/lumber"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/display"
)

type (

	// Configure ...
	Configure struct {
		App            models.App
		Component      models.Component
		BuildID        string
		WarehouseURL   string
		WarehouseToken string
	}

	// payload ...
	payload struct {
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

	fetchPayload struct {
		Component      component      `json:"component"`
		LogvacHost     string         `json:"logvac_host"`
		Member         map[string]int `json:"member"`
		Build          string         `json:"build"`
		Warehouse      string         `json:"warehouse"`
		WarehouseToken string         `json:"warehouse_token"`
	}
)

//
func (configure *Configure) Run() error {

	// quit now if the service was activated already
	if configure.Component.State == ACTIVE {
		return nil
	}

	// set the prefix so the utilExec lumber logging has context
	lumber.Prefix("code:Configure")
	defer lumber.Prefix("")

	display.OpenContext("configuring %s", configure.Component.Name)
	defer display.CloseContext()

	streamer := display.NewStreamer("info")

	// run fetch build command
	fetchPayload := configure.fetchPayload()

	display.StartTask("fetching code")
	if _, err := util.Exec(configure.Component.ID, "fetch", fetchPayload, streamer); err != nil {
		display.ErrorTask()
		return err
	}
	display.StopTask()

	// run configure command
	payload := configure.configurePayload()

	//
	display.StartTask("configuring code")
	if _, err := util.Exec(configure.Component.ID, "configure", payload, streamer); err != nil {
		display.ErrorTask()
		return err
	}
	display.StopTask()

	// run start command
	display.StartTask("starting code")
	if _, err := util.Exec(configure.Component.ID, "start", payload, streamer); err != nil {
		display.ErrorTask()
		return err
	}
	display.StopTask()

	//
	configure.Component.State = ACTIVE
	err := configure.Component.Save()
	if err != nil {
		lumber.Error("code:Configure:Component.Save(): %s", err.Error())
	}
	return err
}

// startPayload ...
func (configure Configure) startPayload() string {
	boxfile := boxfile.New([]byte(configure.App.DeployedBoxfile))
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
func (configure *Configure) configurePayload() string {

	boxfile := boxfile.New([]byte(configure.App.DeployedBoxfile))

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
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

// fetch payload
func (configure *Configure) fetchPayload() string {

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
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

// mounts ...
func (configure *Configure) mounts() []mount {
	boxfile := boxfile.New([]byte(configure.App.DeployedBoxfile))
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
