package models

import (
	"fmt"
)

type (

	Component struct {
		// the docker id
		ID         string         `json:"id"`
		AppID      string 				`json:"app_id"`          
		// name is used for boltdb storage
		Name       string 				`json:"name"`        
		Type       string 				`json:"type"`        
		ExternalIP string 				`json:"external_ip"` 
		InternalIP string 				`json:"internal_ip"` 
		Plan       ComponentPlan  `json:"plan"`        
		State      string 				`json:"state"`       
	}

	ComponentUser struct {
		Username string                 `json:"username"` 
		Password string                 `json:"password"` 
		Meta     map[string]interface{} `json:"meta"`     
	}
)

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
	if err := delete(c.AppID, c.Name); err != nil {
		return fmt.Errorf("failed to delete component: %s", err.Error())
	}
	
	return nil
}

// FindBySlug finds an app by an appID and name
func FindComponentBySlug(appID, name string) (Component, error) {
	
	component := Component{}
	
	if err := get(appID, name, &component); err != nil {
		return component, fmt.Errorf("failed to load component: %s", err.Error())
	}
	
	return component, nil
}

// AllApps loads all of the Apps in the database
func AllComponentsByApp(appID string) ([]Component, error) {
	// list of components to return
	components := []Component{}
	return components, getAll(appID, &components)
}
