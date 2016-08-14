package component

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/nanobox-io/golang-docker-client"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/provider"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/dhcp"
)

type (

	// Setup
	Setup struct {
		App        models.App
		Name       string
		Image      string
		Component    models.Component
		fail       bool
		cleanFuncs []cleanFunc
	}

	// cleanFunc
	cleanFunc func() error
)

//
func (setup *Setup) Run() error {

	// call the cleanup function to ensure we don't leave any bad state
	defer setup.clean()

	// attempt to load the component
	// if the component has not been created it is fine
	// so no errors are caught here
	setup.loadComponent()

	// short-circuit if the service has already progressed past this point
	if setup.Component.State != "initialized" {
		return nil
	}

	if err := setup.downloadImage(); err != nil {
		setup.fail = true
		return err
	}

	if err := setup.reserveIps(); err != nil {
		setup.fail = true
		return err
	}

	if err := setup.launchContainer(); err != nil {
		setup.fail = true
		return err
	}

	if err := setup.attachNetwork(); err != nil {
		setup.fail = true
		return err
	}

	if err := setup.planService(); err != nil {
		setup.fail = true
		return err
	}

	if err := setup.persistService(); err != nil {
		setup.fail = true
		return err
	}

	if err := setup.addEvars(); err != nil {
		setup.fail = true
		return err
	}

	return nil
}

// clean will iterate through the cleanup functions that were registered and
// call them one-by-one
func (setup *Setup) clean() error {
	// short-circuit if we haven't failed
	if !setup.fail {
		return nil
	}

	// iterate through the cleanup functions that were registered and call them
	for _, cleanF := range setup.cleanFuncs {
		if err := cleanF(); err != nil {
			return err
		}
	}

	return nil
}

// loadService fetches the service from the database
func (setup *Setup) loadComponent() error {
	setup.Component, _ = models.FindComponentBySlug(setup.App.ID, setup.Name)

	// set the default state
	if setup.Component.State == "" {
		setup.Component.State = "initialized"
	}

	return nil
}

// downloadImage downloads the docker image
func (setup *Setup) downloadImage() error {
	// Create a pipe to for the JSONMessagesStream to read from
	// pr, pw := io.Pipe()
	// prefix := fmt.Sprintf("%s+ Pulling %s -", stylish.GenerateNestedPrefix(setup.control.DisplayLevel+1), setup.control.Meta["image"])
	//  go print.DisplayJSONMessagesStream(pr, os.Stdout, os.Stdout.Fd(), true, prefix, nil)
	// if _, err := docker.ImagePull(setup.Image, &print.DockerPercentDisplay{Prefix: prefix}); err != nil {

	// TODO: the portion above is commented pending display updates
	if _, err := docker.ImagePull(setup.Image, nil); err != nil {
		return err
	}

	return nil
}

// reserveIps reserves a global and local ip for the container
func (setup *Setup) reserveIps() error {

	// dont reserve a new one if we already have this one
	if setup.Component.InternalIP == "" {
		// first let's see if our local IP was reserved during app creation
		if setup.App.LocalIPs[setup.Name] != "" {

			// assign the localIP from the pre-generated app cache
			setup.Component.InternalIP = setup.App.LocalIPs[setup.Name]
		} else {

			localIP, err := dhcp.ReserveLocal()
			if err != nil {
				return err
			}

			setup.Component.InternalIP = localIP.String()

			setup.cleanFuncs = append(setup.cleanFuncs, func() error {
				return dhcp.ReturnIP(net.ParseIP(setup.Component.InternalIP))
			})
		}
	}

	// dont reserve a new global ip if i already have on
	if setup.Component.ExternalIP == "" {
		// only if this service is portal, we need to use the preview IP
		// in a dev environment there will be no portal installed
		// so the env ip should be available
		// in dev the env ip is used for the dev container
		if setup.Name == "portal" {
			// portal's global ip is the preview ip
			setup.Component.ExternalIP = setup.App.GlobalIPs["env"]
		} else {

			globalIP, err := dhcp.ReserveGlobal()
			if err != nil {
				return err
			}

			setup.Component.ExternalIP = globalIP.String()

			setup.cleanFuncs = append(setup.cleanFuncs, func() error {
				return dhcp.ReturnIP(net.ParseIP(setup.Component.ExternalIP))
			})
		}
	}

	return nil
}

// launchContainer launches and starts a docker container
func (setup *Setup) launchContainer() error {

	name := fmt.Sprintf("nanobox_%s_%s", setup.App.ID, setup.Name)

	config := docker.ContainerConfig{
		Name:    name,
		Image:   setup.Image,
		Network: "virt",
		IP:      setup.Component.InternalIP,
	}

	container, err := docker.CreateContainer(config)
	if err != nil {
		return err
	}

	setup.cleanFuncs = append(setup.cleanFuncs, func() error {
		return docker.ContainerRemove(container.ID)
	})

	setup.Component.ID = container.ID

	return nil
}

// attachNetwork attaches the IP addresses to the container
func (setup *Setup) attachNetwork() error {

	err := provider.AddIP(setup.Component.ExternalIP)
	if err != nil {
		return err
	}

	setup.cleanFuncs = append(setup.cleanFuncs, func() error {
		return provider.RemoveIP(setup.Component.ExternalIP)
	})

	err = provider.AddNat(setup.Component.ExternalIP, setup.Component.InternalIP)
	if err != nil {
		return err
	}

	setup.cleanFuncs = append(setup.cleanFuncs, func() error {
		return provider.RemoveNat(setup.Component.ExternalIP, setup.Component.InternalIP)
	})

	return nil
}

// planService runs the plan hook
func (setup *Setup) planService() error {
	// get the environment so i can get the latest build boxfile
	env, _ := models.FindEnvByID(setup.App.EnvID)

	// get this services config from the boxfile
	boxfile := boxfile.New([]byte(env.BuiltBoxfile))
	boxConfig := boxfile.Node(setup.Name).Node("config")

	planPayload := map[string]interface{}{"config": boxConfig.Parsed}
	jsonPayload, _ := json.Marshal(planPayload)

	// TODO: replace nil with something from the display package 
	p, err := util.Exec(setup.Component.ID, "plan", string(jsonPayload), nil)
	if err != nil {
		return err
	}

	// now set the plans responses data as the components plan object
	err = json.Unmarshal([]byte(p), &setup.Component.Plan)
	if err != nil {
		return fmt.Errorf("persistService:%s", err.Error())
	}

	// set passwords for the users in the plan
	for i := 0; i < len(setup.Component.Plan.Users); i++ {
		setup.Component.Plan.Users[i].Password = util.RandomString(10)
	}

	return nil
}

// persistService saves the service in the database
func (setup *Setup) persistService() error {
	// save service in DB
	setup.Component.AppID = setup.App.ID
	setup.Component.Name = setup.Name
	setup.Component.State = "planned"
	setup.Component.Type = "data"

	// save the service
	return setup.Component.Save()
}

// updateEvars will generate environment variables from the plan
func (setup *Setup) addEvars() error {

	// fetch the environment variables model
	envVars := setup.App.Evars

	// create a prefix for each of the environment variables.
	// for example, if the service is 'data.db' the prefix
	// would be DATA_DB. Dots are replaced with underscores,
	// and characters are uppercased.
	prefix := strings.ToUpper(strings.Replace(setup.Component.Name, ".", "_", -1))

	// we need to create an host evar that holds the IP of the service
	envVars[fmt.Sprintf("%s_HOST", prefix)] = setup.Component.InternalIP

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
	for _, user := range setup.Component.Plan.Users {
		// add this username to the list
		users = append(users, user.Username)

		// generate the corresponding evar for the password
		key := fmt.Sprintf("%s_%s_PASS", prefix, strings.ToUpper(user.Username))
		envVars[key] = user.Password

		// if this user is the default user
		// set additional default env vars
		if user.Username == setup.Component.Plan.DefaultUser {
			envVars[fmt.Sprintf("%s_USER", prefix)] = user.Username
			envVars[fmt.Sprintf("%s_PASS", prefix)] = user.Password
		}
	}

	// if there are users, create an environment variable to represent the list
	if len(users) > 0 {
		envVars[fmt.Sprintf("%s_USERS", prefix)] = strings.Join(users, " ")
	}

	// persist the evars
	setup.App.Evars = envVars
	return setup.App.Save()
}
