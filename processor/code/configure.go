package code

import (
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type (

	// processCodeConfigure ...
	processCodeConfigure struct {
		control processor.ProcessControl
	}

	// payload ...
	payload struct {
		LogvacHost   string            `json:"logvac_host,omitempty"`
		Config       interface{}       `json:"config,omitempty"`
		Component    component         `json:"component,omitempty"`
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
)

//
func init() {
	processor.Register("code_configure", codeConfigureFunc)
}

//
func codeConfigureFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// make sure i was given a name and image
	if control.Meta["name"] == "" || control.Meta["boxfile"] == "" {
		return nil, errMissingImageOrName
	}

	return &processCodeConfigure{control: control}, nil
}

//
func (codeConfigure *processCodeConfigure) Results() processor.ProcessControl {
	return codeConfigure.control
}

//
func (codeConfigure *processCodeConfigure) Process() error {

	// get the service from the database
	service := models.Service{}

	//
	if err := data.Get(util.AppName(), codeConfigure.control.Meta["name"], &service); err != nil {
		return err
	}

	fmt.Println("-> configuring", codeConfigure.control.Meta["name"])

	// quit now if the service was activated already
	if service.State == ACTIVE {
		return nil
	}

	// run fetch build command
	if _, err := util.Exec(service.ID, "fetch", fmt.Sprintf(`{"build":"%s","warehouse":"%s","warehouse_token":"%s"}`, codeConfigure.control.Meta["build_id"], codeConfigure.control.Meta["warehouse_url"], codeConfigure.control.Meta["warehouse_token"]), processor.ExecWriter()); err != nil {
		return err
	}

	// run configure command
	payload, err := codeConfigure.configurePayload()
	if err != nil {
		return err
	}

	//
	if _, err = util.Exec(service.ID, "configure", payload, nil); err != nil {
		return err
	}

	// run start command
	if _, err = util.Exec(service.ID, "start", codeConfigure.startPayload(), nil); err != nil {
		return err
	}

	//
	service.State = ACTIVE
	if err := data.Put(util.AppName(), codeConfigure.control.Meta["name"], service); err != nil {
		return err
	}

	return nil
}

// startPayload ...
func (codeConfigure processCodeConfigure) startPayload() string {
	boxfile := boxfile.New([]byte(codeConfigure.control.Meta["boxfile"]))
	pload := payload{
		Config: boxfile.Node(codeConfigure.control.Meta["name"]).Value("config"),
		Start:  boxfile.Node(codeConfigure.control.Meta["name"]).StringValue("start"),
	}

	bytes, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}

// configurePayload ...
func (codeConfigure *processCodeConfigure) configurePayload() (string, error) {

	me := models.Service{}
	err := data.Get(util.AppName(), codeConfigure.control.Meta["name"], &me)
	boxfile := boxfile.New([]byte(codeConfigure.control.Meta["boxfile"]))

	logvac := models.Service{}
	data.Get(util.AppName(), "logvac", &logvac)

	pload := payload{
		LogvacHost: logvac.InternalIP,
		Config:     boxfile.Node(codeConfigure.control.Meta["name"]).Value("config"),
		Component: component{
			Name: "whydoesthismatter",
			UID:  codeConfigure.control.Meta["name"],
			ID:   me.ID,
		},
		Mounts:       codeConfigure.mounts(),
		WritableDirs: boxfile.Node(codeConfigure.control.Meta["name"]).Value("writable_dirs"),
		Transform:    boxfile.Node("code.deploy").Value("transform"),
		Env:          codeConfigure.env(),
		LogWatches:   boxfile.Node(codeConfigure.control.Meta["name"]).Value("log_watch"),
		Start:        boxfile.Node(codeConfigure.control.Meta["name"]).Value("start"),
	}

	bytes, err := json.Marshal(pload)
	return string(bytes), err
}

// mounts ...
func (codeConfigure *processCodeConfigure) mounts() []mount {
	boxfile := boxfile.New([]byte(codeConfigure.control.Meta["boxfile"]))
	boxNetworkDirs := boxfile.Node(codeConfigure.control.Meta["name"]).Node("network_dirs")

	m := []mount{}
	for _, node := range boxNetworkDirs.Nodes() {
		// i think i store these as data.name
		// cleanNode := regexp.MustCompile(`.+\.`).ReplaceAllString(node, "")
		service := models.Service{}
		err := data.Get(util.AppName(), node, &service)
		if err != nil {
			// skip because of problems
			fmt.Println("cant get service:", err)
			continue
		}
		if !service.Plan.BehaviorPresent("mountable") || service.Plan.MountProtocol == "" {
			// skip because of problems
			fmt.Println("non mountable service", service.Name)
			continue
		}
		m = append(m, mount{service.InternalIP, service.Plan.MountProtocol, boxNetworkDirs.StringSliceValue(node)})

	}

	return m
}

// env ...
func (codeConfigure *processCodeConfigure) env() map[string]string {
	envVars := models.EnvVars{}
	data.Get(util.AppName()+"_meta", "env", &envVars)

	return envVars
}
