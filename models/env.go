package models

import (
	"fmt"
)

type Env struct {
	ID        string
	Directory string
	Name      string

	// Links map a local app to multiple production apps, by an alias
	Links map[string]string
	// the boxfile from the most recent build
	BuiltBoxfile string
}

// Save persists the Env to the database
func (e *Env) Save() error {

	if err := put("envs", e.ID, e); err != nil {
		return fmt.Errorf("failed to save app: %s", err.Error())
	}

	return nil
}

// Delete deletes the app record from the database
func (e *Env) Delete() error {

	if err := delete("envs", e.ID); err != nil {
		return fmt.Errorf("failed to delete app: %s", err.Error())
	}

	return nil
}

// FindEnvByID finds an app by an ID
func FindEnvByID(ID string) (Env, error) {

	env := Env{}

	if err := get("envs", ID, &env); err != nil {
		return env, fmt.Errorf("failed to load env: %s", err.Error())
	}

	return env, nil
}

// AllEnvs loads all of the Envs in the database
func AllEnvs() ([]Env, error) {
	// list of apps to return
	apps := []Env{}

	return apps, getAll("envs", &apps)
}
