package code

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/data"
)

type codeConfigure struct {
	config processor.ProcessConfig
}

type payload struct {
	LogvacHost   string            `json:"logvac_host,omitempty"`
	Config       interface{}       `json:"config"`
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

func codeConfigureFunc(config processor.ProcessConfig) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// make sure i was given a name and image
	if config.Meta["name"] == "" || config.Meta["boxfile"] == "" {
		return nil, missingImageOrName
	}

	return &codeConfigure{config: config}, nil
}

func (self codeConfigure) startPayload() string {
	boxfile := boxfile.New([]byte(self.config.Meta["boxfile"]))
	pload := payload{
		Config: boxfile.Node(self.config.Meta["name"]).Value("config"),
		Start:  boxfile.Node(self.config.Meta["name"]).StringValue("start"),
	}

	bytes, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (self *codeConfigure) configurePayload() (string, error) {
	me := models.Service{}
	err := data.Get(util.AppName(), self.config.Meta["name"], &me)

	boxfile := boxfile.New([]byte(self.config.Meta["boxfile"]))

	logvac := models.Service{}
	data.Get(util.AppName(), "logvac", &logvac)

	pload := payload{
		LogvacHost: logvac.InternalIP,
		Config:     boxfile.Node(self.config.Meta["name"]).Value("config"),
		Component: component{
			Name: "whydoesthismatter",
			UID:  self.config.Meta["name"],
			ID:   me.ID,
		},
		Mounts:       self.mounts(),
		WritableDirs: boxfile.Node(self.config.Meta["name"]).Value("writable_dirs"),
		Transform:    boxfile.Node("code.deploy").Value("transform"),
		Env:          self.env(),
		LogWatches:   boxfile.Node(self.config.Meta["name"]).Value("log_watch"),
		Start:        boxfile.Node(self.config.Meta["name"]).Value("start"),
	}

	bytes, err := json.Marshal(pload)
	return string(bytes), err
}

func (self *codeConfigure) mounts() []mount {
	boxfile := boxfile.New([]byte(self.config.Meta["boxfile"]))
	boxNetworkDirs := boxfile.Node(self.config.Meta["name"]).Node("network_dirs")

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

	serviceNames, _ := data.Keys(util.AppName())
	for _, serviceName := range serviceNames {
		// only look at data services
		if strings.HasPrefix(serviceName, "data") {
			envName := strings.ToUpper(strings.Replace(serviceName, ".", "_", -1))
			service := models.Service{}
			users := []string{}
			data.Get(util.AppName(), serviceName, &service)
			envVars[envName+"_HOST"] = service.InternalIP
			for _, user := range service.Plan.Users {
				users = append(users, user.Username)
				envVars[fmt.Sprintf("%s_%s_PW", envName, strings.ToUpper(user.Username))] = user.Password
			}
			envVars[envName+"_USERS"] = strings.Join(users, " ")
		}

	}
	return envVars
}

func (self *codeConfigure) Results() processor.ProcessConfig {
	return self.config
}

func (self *codeConfigure) Process() error {

	// get the service from the database
	service := models.Service{}
	err := data.Get(util.AppName(), self.config.Meta["name"], &service)
	if err != nil {
		// cannot start a service that wasnt setup (ie saved in the database)
		return err
	}

	if service.Started {
		return nil
	}

	// run fetch build command
	output, err := util.Exec(service.ID, "fetch", fmt.Sprintf(`{"build":"%s","warehouse":"%s","warehouse_token":"%s"}`, self.config.Meta["build_id"], self.config.Meta["warehouse_url"], self.config.Meta["warehouse_token"]))
	if err != nil {
		fmt.Println(output)
		return err
	}

	// run configure command
	payload, err := self.configurePayload()
	if err != nil {
		fmt.Println("error building configure payload", err)
		return err
	}
	fmt.Println("configure payload", payload)
	output, err = util.Exec(service.ID, "configure", payload)
	if err != nil {
		fmt.Println(output)
		return err
	}

	// run start command
	fmt.Println("start payload", self.startPayload())
	output, err = util.Exec(service.ID, "start", self.startPayload())
	if err != nil {
		fmt.Println(output)
		return err
	}

	service.Started = true
	err = data.Put(util.AppName(), self.config.Meta["name"], service)
	if err != nil {
		return err
	}

	return nil
}
