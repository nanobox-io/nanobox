package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"
	"github.com/nanobox-io/nanobox-golang-stylish"

	dockType "github.com/docker/engine-api/types"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/data"
	"github.com/nanobox-io/nanobox/util/dhcp"
	"github.com/nanobox-io/nanobox/util/print"
)

type (

	// processServiceSetup
	processServiceSetup struct {
		control    processor.ProcessControl
		service    models.Service
		localIP    net.IP
		globalIP   net.IP
		container  dockType.ContainerJSON
		plan       string
		fail       bool
		cleanFuncs []cleanFunc
	}

	// cleanFunc
	cleanFunc func() error
)

//
func init() {
	processor.Register("service_setup", serviceSetupFn)
}

//
func serviceSetupFn(control processor.ProcessControl) (processor.Processor, error) {
	// confirm the provider is an accessable one that we support.
	// ensure we have a name and immage
	if control.Meta["name"] == "" || control.Meta["image"] == "" {
		return nil, errors.New("missing image or name")
	}

	// add a label if im missing one
	if control.Meta["label"] == "" {
		control.Meta["label"] = control.Meta["name"]
	}

	return &processServiceSetup{
		control:    control,
		cleanFuncs: make([]cleanFunc, 0),
	}, nil
}

//
func (serviceSetup processServiceSetup) Results() processor.ProcessControl {
	return serviceSetup.control
}

//
func (serviceSetup *processServiceSetup) Process() error {

	header := fmt.Sprintf("Launching %s...", serviceSetup.control.Meta["label"])
	serviceSetup.control.Display(stylish.Bullet(header))

	// call the cleanup function to ensure we don't leave any bad state
	defer serviceSetup.clean()

	if err := serviceSetup.loadService(); err != nil {
		serviceSetup.fail = true
		return err
	}

	// short-circuit if the service has already progressed past this point
	if serviceSetup.service.State != "initialized" {
		return nil
	}

	if err := serviceSetup.downloadImage(); err != nil {
		serviceSetup.fail = true
		return err
	}

	if err := serviceSetup.reserveIps(); err != nil {
		serviceSetup.fail = true
		return err
	}

	if err := serviceSetup.launchContainer(); err != nil {
		serviceSetup.fail = true
		return err
	}

	if err := serviceSetup.attachNetwork(); err != nil {
		serviceSetup.fail = true
		return err
	}

	if err := serviceSetup.planService(); err != nil {
		serviceSetup.fail = true
		return err
	}

	if err := serviceSetup.persistService(); err != nil {
		serviceSetup.fail = true
		return err
	}

	if err := serviceSetup.addEnvVars(); err != nil {
		serviceSetup.fail = true
		return err
	}

	return nil
}

// clean will iterate through the cleanup functions that were registered and
// call them one-by-one
func (serviceSetup *processServiceSetup) clean() error {
	// short-circuit if we haven't failed
	if !serviceSetup.fail {
		return nil
	}

	// iterate through the cleanup functions that were registered and call them
	for _, cleanF := range serviceSetup.cleanFuncs {
		if err := cleanF(); err != nil {
			return err
		}
	}

	return nil
}

// validateMeta ensures we were given a name and image
func (serviceSetup *processServiceSetup) validateMeta() error {
	return nil
}

// loadService fetches the service from the database
func (serviceSetup *processServiceSetup) loadService() error {
	// the service really shouldn't exist yet, so let's not return the error if it fails
	data.Get(config.AppName(), serviceSetup.control.Meta["name"], &serviceSetup.service)

	// set the default state
	if serviceSetup.service.State == "" {
		serviceSetup.service.State = "initialized"
	}

	return nil
}

// downloadImage downloads the docker image
func (serviceSetup *processServiceSetup) downloadImage() error {
	// Create a pipe to for the JSONMessagesStream to read from
	// pr, pw := io.Pipe()
	prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(serviceSetup.control.DisplayLevel+1), serviceSetup.control.Meta["image"])
	//  go print.DisplayJSONMessagesStream(pr, os.Stdout, os.Stdout.Fd(), true, prefix, nil)
	if _, err := docker.ImagePull(serviceSetup.control.Meta["image"], &print.DockerPercentDisplay{Prefix: prefix}); err != nil {
		return err
	}

	return nil
}

// reserveIps reserves a global and local ip for the container
func (serviceSetup *processServiceSetup) reserveIps() error {
	var err error
	serviceSetup.localIP, err = dhcp.ReserveLocal()
	if err != nil {
		return err
	}

	serviceSetup.cleanFuncs = append(serviceSetup.cleanFuncs, func() error {
		return dhcp.ReturnIP(serviceSetup.localIP)
	})

	serviceSetup.globalIP, err = dhcp.ReserveGlobal()
	if err != nil {
		return err
	}

	serviceSetup.cleanFuncs = append(serviceSetup.cleanFuncs, func() error {
		return dhcp.ReturnIP(serviceSetup.globalIP)
	})

	return nil
}

// launchContainer launches and starts a docker container
func (serviceSetup *processServiceSetup) launchContainer() error {
	config := docker.ContainerConfig{
		Name:    fmt.Sprintf("nanobox-%s-%s", config.AppName(), serviceSetup.control.Meta["name"]),
		Image:   serviceSetup.control.Meta["image"],
		Network: "virt",
		IP:      serviceSetup.localIP.String(),
	}

	serviceSetup.control.Info(stylish.SubBullet("Starting container..."))
	container, err := docker.CreateContainer(config)
	if err != nil {
		return err
	}

	serviceSetup.cleanFuncs = append(serviceSetup.cleanFuncs, func() error {
		return docker.ContainerRemove(container.ID)
	})

	serviceSetup.container = container

	return nil
}

// attachNetwork attaches the IP addresses to the container
func (serviceSetup *processServiceSetup) attachNetwork() error {
	label := "Bridging container to host network..."
	serviceSetup.control.Info(stylish.SubBullet(label))

	err := provider.AddIP(serviceSetup.globalIP.String())
	if err != nil {
		return err
	}

	serviceSetup.cleanFuncs = append(serviceSetup.cleanFuncs, func() error {
		return provider.RemoveIP(serviceSetup.globalIP.String())
	})

	err = provider.AddNat(serviceSetup.globalIP.String(), serviceSetup.localIP.String())
	if err != nil {
		return err
	}

	serviceSetup.cleanFuncs = append(serviceSetup.cleanFuncs, func() error {
		return provider.RemoveNat(serviceSetup.globalIP.String(), serviceSetup.localIP.String())
	})

	return nil
}

// planService runs the plan hook
func (serviceSetup *processServiceSetup) planService() error {
	serviceSetup.control.Info(stylish.SubBullet("Gathering service requirements..."))

	boxfile := boxfile.New([]byte(serviceSetup.control.Meta["boxfile"]))
	boxConfig := boxfile.Node(serviceSetup.control.Meta["name"]).Node("config")
	planPayload := map[string]interface{}{"config": boxConfig.Parsed}
	jsonPayload, _ := json.Marshal(planPayload)

	p, err := util.Exec(serviceSetup.container.ID, "plan", string(jsonPayload), processor.ExecWriter())
	if err != nil {
		return err
	}
	serviceSetup.plan = p

	return nil
}

// persistService saves the service in the database
func (serviceSetup *processServiceSetup) persistService() error {
	// save service in DB
	serviceSetup.service.ID = serviceSetup.container.ID
	serviceSetup.service.Name = serviceSetup.control.Meta["name"]
	serviceSetup.service.ExternalIP = serviceSetup.globalIP.String()
	serviceSetup.service.InternalIP = serviceSetup.localIP.String()
	serviceSetup.service.State = "planned"
	serviceSetup.service.Type = "data"

	err := json.Unmarshal([]byte(serviceSetup.plan), &serviceSetup.service.Plan)
	if err != nil {
		return fmt.Errorf("persistService:%s", err.Error())
	}
	for i := 0; i < len(serviceSetup.service.Plan.Users); i++ {
		serviceSetup.service.Plan.Users[i].Password = util.RandomString(10)
	}

	// save the service
	err = data.Put(config.AppName(), serviceSetup.control.Meta["name"], &serviceSetup.service)
	if err != nil {
		return err
	}

	return nil
}

// updateEnvVars will generate environment variables from the plan
func (serviceSetup *processServiceSetup) addEnvVars() error {
	// fetch the environment variables model
	envVars := models.EnvVars{}
	data.Get(config.AppName()+"_meta", "env", &envVars)

	// create a prefix for each of the environment variables.
	// for example, if the service is 'data.db' the prefix
	// would be DATA_DB. Dots are replaced with underscores,
	// and characters are uppercased.
	prefix := strings.ToUpper(strings.Replace(serviceSetup.service.Name, ".", "_", -1))

	// we need to create an host evar that holds the IP of the service
	envVars[fmt.Sprintf("%s_HOST", prefix)] = serviceSetup.service.InternalIP

	// we need to create evars that contain usernames and passwords
	//
	// during the plan phase, the service was informed of potentially
	// 	1 - users (all of the users)
	// 	2 - user (the default user)
	//
	// First, we need to create an evar that contains the list of users.
	// 	{prefix}_USERS
	//
	// Each user provided was given a password. For every user specified,
	// we need to create a corresponding evars to store the password:
	//  {prefix}_{username}_PASS
	//
	// Lastly, if a default user was provided, we need to create a pair
	// of environment variables as a convenience to the user:
	// 	1 - {prefix}_USER
	// 	2 - {prefix}_PASS

	// create a slice of user strings that we will use to generate the list of users
	users := []string{}

	// users will have been loaded into the service plan, so let's iterate
	for _, user := range serviceSetup.service.Plan.Users {
		// add this username to the list
		users = append(users, user.Username)

		// generate the corresponding evar for the password
		key := fmt.Sprintf("%s_%s_PASS", prefix, strings.ToUpper(user.Username))
		envVars[key] = user.Password

		// if this user is the default user
		// set additional default env vars
		if user.Username == serviceSetup.service.Plan.DefaultUser {
			envVars[fmt.Sprintf("%s_USER", prefix)] = user.Username
			envVars[fmt.Sprintf("%s_PASS", prefix)] = user.Password
		}
	}

	// if there are users, create an environment variable to represent the list
	if len(users) > 0 {
		envVars[fmt.Sprintf("%s_USERS", prefix)] = strings.Join(users, " ")
	}

	// persist the evars
	if err := data.Put(config.AppName()+"_meta", "env", envVars); err != nil {
		return err
	}

	return nil
}
