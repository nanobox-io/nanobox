package models

import (
	"fmt"

	"github.com/nanobox-io/nanobox/util/config"
)

type App struct {
	EnvID string
	ID    string
	Name  string
	// State is used to ensure we don't setup this environment multiple times
	State  string
	Status string
	// Appironment variables available to the environment
	Evars map[string]string
	// There are certain global ips that need to be reserved across container
	// lifetimes. The dev ip and preview ip are examples. We'll store those here.
	GlobalIPs map[string]string
	// There are also certain platform service ips that need to 1) remain constant
	// even if the component were repaired and 2) be available even before the
	// component is. logvac and mist ips are examples. We'll store those here.
	LocalIPs map[string]string
	// the boxfile from the most recent deploy
	DeployedBoxfile string
}

// InNew returns true if the App hasn't been created yet
func (a *App) IsNew() bool {
	return a.ID == ""
}

// Save persists the App to the database
func (a *App) Save() error {

	if err := put(a.EnvID, a.ID, a); err != nil {
		return fmt.Errorf("failed to save app %s", err.Error())
	}

	return nil
}

// Delete deletes the app record from the database
func (a *App) Delete() error {

	if err := destroy(a.EnvID, a.ID); err != nil {
		return fmt.Errorf("failed to destroy app %s", err.Error())
	}

	return nil
}

// Generate populates an App with data and persists the record
func (a *App) Generate(env *Env, name string) error {

	// short-circuit if this record has already been generated
	if !a.IsNew() {
		return nil
	}

	// if no env is present we will need to create one
	if env == nil {
		env = &Env{ID: config.EnvID()}
	}

	a.EnvID = env.ID
	a.ID = fmt.Sprintf("%s_%s", env.ID, name)
	a.Name = name
	a.State = "initialized"
	a.GlobalIPs = map[string]string{}
	a.LocalIPs = map[string]string{}
	a.Evars = map[string]string{
		"APP_NAME": name,
	}

	return a.Save()
}

func (a *App) Env() (*Env, error) {
	return FindEnvByID(a.EnvID)
}

func (a *App) Components() ([]*Component, error) {
	return AllComponentsByApp(a.ID)
}

// FindBySlug finds an app by an appID and name
func FindAppBySlug(envID, name string) (*App, error) {

	app := &App{}

	key := fmt.Sprintf("%s_%s", envID, name)

	if err := get(envID, key, &app); err != nil {
		return app, fmt.Errorf("failed to load app: %s", err.Error())
	}

	return app, nil
}

// AllApps loads all of the Apps across all environments
func AllApps() ([]*App, error) {
	apps := []*App{}

	// load all envs
	envs, err := AllEnvs()
	if err != nil {
		return apps, fmt.Errorf("failed to load envs: %s", err.Error())
	}

	for _, env := range envs {
		// load all apps by env
		envApps, err := AllAppsByEnv(env.ID)
		if err != nil {
			return apps, fmt.Errorf("failed to load env apps: %s", err.Error())
		}

		for _, app := range envApps {
			// append to the apps that we'll return
			apps = append(apps, app)
		}
	}

	return apps, nil
}

// AllApps loads all of the Apps in the database
func AllAppsByEnv(envID string) ([]*App, error) {
	// list of envs to return
	apps := []*App{}

	return apps, getAll(envID, &apps)
}

// AllAppsByStatus loads all of the Apps filtering by status
func AllAppsByStatus(status string) ([]*App, error) {
	apps := []*App{}

	all, err := AllApps()
	if err != nil {
		return apps, fmt.Errorf("failed to load all apps: %s", err.Error())
	}

	for _, app := range all {
		if app.Status == status {
			apps = append(apps, app)
		}
	}

	return apps, nil
}

// lookup an app id from an alias
func AppIDByAlias(alias string) string {
	env, _ := FindEnvByID(config.EnvID())
	
	// set alias to default if none provided
	if alias == "" {
		alias = "default"
	}
	
	app, ok := env.Links[alias]
	if !ok {
		return ""
	}

	return app.ID
}
