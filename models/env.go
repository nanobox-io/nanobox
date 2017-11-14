package models

import (
	"fmt"
	"time"

	"github.com/nanobox-io/nanobox/util/config"
)

// Env ...
type Env struct {
	ID        string
	Directory string
	Name      string

	// Remotes map a local app to multiple production apps, by an alias
	Remotes map[string]Remote
	// the boxfile from the most recent build
	BuiltBoxfile  string
	UserBoxfile   string
	BuiltID       string
	DeployedID    string
	LastBuild     time.Time
	LastCompile   time.Time
	BuildTriggers map[string]string
}

// Remote ...
type Remote struct {
	ID       string
	Name     string
	Endpoint string
}

// IsNew returns true if the Env hasn't been created yet
func (e *Env) IsNew() bool {
	return e.ID == ""
}

// Save persists the Env to the database
func (e *Env) Save() error {

	if err := put("envs", e.ID, e); err != nil {
		return fmt.Errorf("failed to save env: %s", err.Error())
	}

	return nil
}

// Delete deletes the app record from the database
func (e *Env) Delete() error {

	if err := destroy("envs", e.ID); err != nil {
		return fmt.Errorf("failed to delete env: %s", err.Error())
	}

	return nil
}

// Generate populates an Env from config data and persists the record
func (e *Env) Generate() error {

	// short-circuit if this record has already been generated
	if !e.IsNew() {
		return nil
	}

	// populate the data from the config package
	e.ID = config.EnvID()
	e.Directory = config.LocalDir()
	e.Name = config.LocalDirName()
	e.Remotes = map[string]Remote{}

	return e.Save()
}

// Apps get a list of the apps that belong to this
func (e *Env) Apps() ([]*App, error) {
	return AllAppsByEnv(e.ID)
}

// FindEnvByID finds an app by an ID
func FindEnvByID(ID string) (*Env, error) {

	env := &Env{}

	if err := get("envs", ID, &env); err != nil {
		return env, fmt.Errorf("failed to load env: %s", err.Error())
	}

	return env, nil
}

// FindEnvByName finds an app by a name
func FindEnvByName(name string) (*Env, error) {
	apps := []*Env{}

	if err := getAll("envs", &apps); err != nil {
		return nil, fmt.Errorf("failed to get all envs: %s", err.Error())
	}

	for i := range apps {
		if apps[i].Name == name {
			return apps[i], nil
		}
	}

	return &Env{}, fmt.Errorf("failed to find env '%s' by name", name)
}

// AllEnvs loads all of the Envs in the database
func AllEnvs() ([]*Env, error) {
	// list of apps to return
	apps := []*Env{}

	return apps, getAll("envs", &apps)
}
