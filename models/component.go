package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nanobox-io/nanobox/util"
)

type (

	// Component ...
	Component struct {
		// the docker id
		ID    string `json:"id"`
		AppID string `json:"app_id"`
		EnvID string `json:"env_id"`
		// name is used for boltdb storage
		Name       string        `json:"name"`
		Label      string        `json:"label"`
		Image      string        `json:"image"`
		Type       string        `json:"type"`
		ExternalIP string        `json:"external_ip"`
		InternalIP string        `json:"internal_ip"`
		Plan       ComponentPlan `json:"plan"`
		State      string        `json:"state"`
	}

	// ComponentUser ...
	ComponentUser struct {
		Username string                 `json:"username"`
		Password string                 `json:"password"`
		Meta     map[string]interface{} `json:"meta"`
	}
)

// IsNew returns true if the Component hasn't been created yet
func (c *Component) IsNew() bool {
	return c.ID == ""
}

// Save persists the Component to the database
func (c *Component) Save() error {
	// store under the apps id and
	if err := put(c.AppID, c.Name, c); err != nil {
		return fmt.Errorf("failed to save component: %s", err.Error())
	}

	return nil
}

// Delete deletes the component record from the database
func (c *Component) Delete() error {
	if err := destroy(c.AppID, c.Name); err != nil {
		return fmt.Errorf("failed to delete component: %s", err.Error())
	}

	return nil
}

// Generate populates a Component with data and persists the record
func (c *Component) Generate(app *App, ttype string) error {
	// short-circuit if the component is already created
	if !c.IsNew() {
		return nil
	}

	c.AppID = app.ID
	c.EnvID = app.EnvID
	c.State = "initialized"
	c.Type = ttype

	return c.Save()
}

// GeneratePlan generates the plan from the plan hook output
func (c *Component) GeneratePlan(data string) error {
	// if there is no plan data then i cant do any planning
	if data == "" {
		return nil
	}
	// decode the json directly into the component plan
	if err := json.Unmarshal([]byte(data), &c.Plan); err != nil {
		return fmt.Errorf("failed to decode the plan data:%s", err.Error())
	}

	// set passwords for the users in the plan
	for i := 0; i < len(c.Plan.Users); i++ {
		c.Plan.Users[i].Password = util.RandomString(10)
	}

	return c.Save()
}

// GenerateEvars generates the evars for this component
func (c *Component) GenerateEvars(app *App) error {
	// create a prefix for each of the environment variables.
	// for example, if the service is 'data.db' the prefix
	// would be DATA_DB. Dots are replaced with underscores,
	// and characters are uppercased.
	prefix := strings.ToUpper(strings.Replace(c.Name, ".", "_", -1))

	// we need to create an host evar that holds the IP of the service
	app.Evars[fmt.Sprintf("%s_HOST", prefix)] = c.InternalIP

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
	for _, user := range c.Plan.Users {
		// add this username to the list
		users = append(users, user.Username)

		// generate the corresponding evar for the password
		key := fmt.Sprintf("%s_%s_PASS", prefix, strings.ToUpper(user.Username))
		app.Evars[key] = user.Password

		// if this user is the default user
		// set additional default env vars
		if user.Username == c.Plan.DefaultUser {
			app.Evars[fmt.Sprintf("%s_USER", prefix)] = user.Username
			app.Evars[fmt.Sprintf("%s_PASS", prefix)] = user.Password
		}
	}

	// if there are users, create an environment variable to represent the list
	if len(users) > 0 {
		app.Evars[fmt.Sprintf("%s_USERS", prefix)] = strings.Join(users, " ")
	}

	return app.Save()
}

// PurgeEvars purges the generated evars for a component
func (c *Component) PurgeEvars(a *App) error {

	// create a prefix for each of the environment variables.
	// for example, if the service is 'data.db' the prefix
	// would be DATA_DB. Dots are replaced with underscores,
	// and characters are uppercased.
	prefix := strings.ToUpper(strings.Replace(c.Name, ".", "_", -1))

	// we loop over all environment variables and see if the key contains
	// the prefix above. If so, we delete the item.
	for key := range a.Evars {
		if strings.HasPrefix(key, prefix) {
			delete(a.Evars, key)
		}
	}

	// persist the app with the new env vars
	return a.Save()
}

// FindComponentBySlug finds a component by an appID and name
func FindComponentBySlug(appID, name string) (*Component, error) {

	component := &Component{}

	if err := get(appID, name, &component); err != nil {
		return component, fmt.Errorf("failed to load component: %s", err.Error())
	}

	return component, nil
}

// AllComponentsByApp loads all of the components in the database
func AllComponentsByApp(appID string) ([]*Component, error) {
	// list of components to return
	components := []*Component{}
	return components, getAll(appID, &components)
}
