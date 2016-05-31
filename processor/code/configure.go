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

type codeConfigure struct {
	control processor.ProcessControl
}

type payload struct {
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

type component struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
	ID   string `json:"id"`
}

type mount struct {
	Host     string   `json:"host"`
	Protocol string   `json:"protocol"`
	Shares   []string `json:"shares"`
}

func init() {
	processor.Register("code_configure", codeConfigureFunc)
}

func codeConfigureFunc(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// make sure i was given a name and image
	if control.Meta["name"] == "" || control.Meta["boxfile"] == "" {
		return nil, missingImageOrName
	}

	return &codeConfigure{control: control}, nil
}

func (self codeConfigure) startPayload() string {
	boxfile := boxfile.New([]byte(self.control.Meta["boxfile"]))
	pload := payload{
		Config: boxfile.Node(self.control.Meta["name"]).Value("config"),
		Start:  boxfile.Node(self.control.Meta["name"]).StringValue("start"),
	}

	bytes, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (self *codeConfigure) configurePayload() (string, error) {

	me := models.Service{}
	err := data.Get(util.AppName(), self.control.Meta["name"], &me)
	boxfile := boxfile.New([]byte(self.control.Meta["boxfile"]))

	logvac := models.Service{}
	data.Get(util.AppName(), "logvac", &logvac)

	pload := payload{
		LogvacHost: logvac.InternalIP,
		Config:     boxfile.Node(self.control.Meta["name"]).Value("config"),
		Component: component{
			Name: "whydoesthismatter",
			UID:  self.control.Meta["name"],
			ID:   me.ID,
		},
		Mounts:       self.mounts(),
		WritableDirs: boxfile.Node(self.control.Meta["name"]).Value("writable_dirs"),
		Transform:    boxfile.Node("code.deploy").Value("transform"),
		Env:          self.env(),
		LogWatches:   boxfile.Node(self.control.Meta["name"]).Value("log_watch"),
		Start:        boxfile.Node(self.control.Meta["name"]).Value("start"),
	}

	bytes, err := json.Marshal(pload)
	return string(bytes), err
}

func (self *codeConfigure) mounts() []mount {
	boxfile := boxfile.New([]byte(self.control.Meta["boxfile"]))
	boxNetworkDirs := boxfile.Node(self.control.Meta["name"]).Node("network_dirs")

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

func (self *codeConfigure) env() map[string]string {
	envVars := models.EnvVars{}
	data.Get(util.AppName()+"_meta", "env", &envVars)
	return envVars
}

func (self *codeConfigure) Results() processor.ProcessControl {
	return self.control
}

func (self *codeConfigure) Process() error {

	// get the service from the database
	service := models.Service{}
	err := data.Get(util.AppName(), self.control.Meta["name"], &service)
	if err != nil {
		// cannot start a service that wasnt setup (ie saved in the database)
		return err
	}

	fmt.Println("-> configuring", self.control.Meta["name"])

	// quit now if the service was activated already
	if service.State == "active" {
		return nil
	}

	// run fetch build command
	_, err = util.Exec(service.ID, "fetch", fmt.Sprintf(`{"build":"%s","warehouse":"%s","warehouse_token":"%s"}`, self.control.Meta["build_id"], self.control.Meta["warehouse_url"], self.control.Meta["warehouse_token"]), processor.ExecWriter())
	if err != nil {
		return err
	}

	// run configure command
	payload, err := self.configurePayload()
	if err != nil {
		return err
	}
	_, err = util.Exec(service.ID, "configure", payload, nil)
	if err != nil {
		return err
	}

	// run start command
	_, err = util.Exec(service.ID, "start", self.startPayload(), nil)
	if err != nil {
		return err
	}

	service.State = "active"
	err = data.Put(util.AppName(), self.control.Meta["name"], service)
	if err != nil {
		return err
	}

	return nil
}
