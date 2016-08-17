package models

import (
	"fmt"
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

// Save persists the App to the database
func (a *App) Save() error {

	if err := put(a.EnvID, a.ID, a); err != nil {
		return fmt.Errorf("failed to save app %s", err.Error())
	}

	return nil
}

// Delete deletes the app record from the database
func (a *App) Delete() error {

	if err := delete(a.EnvID, a.ID); err != nil {
		return fmt.Errorf("failed to delete app %s", err.Error())
	}

	return nil
}

// FindBySlug finds an app by an appID and name
func FindAppBySlug(envID, name string) (App, error) {

	app := App{}

	key := fmt.Sprintf("%s_%s", envID, name)

	if err := get(envID, key, &app); err != nil {
		return app, fmt.Errorf("failed to load app: %s", err.Error())
	}

	return app, nil
}

// AllApps loads all of the Apps in the database
func AllAppsByEnv(envID string) ([]App, error) {
	// list of envs to return
	apps := []App{}

	return apps, getAll(envID, &apps)
}
