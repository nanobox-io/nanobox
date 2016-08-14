package dev

import (
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/nanobox/models"
)

// Info ...
type Info struct {
	App models.App
}

//
func (info Info) Run() error {

	//
	components, _ := models.AllComponentsByApp(info.App.ID)

	//
	for _, component := range components {
		if component.Name != "builds" {
			bytes, _ := json.MarshalIndent(component, "", "  ")
			fmt.Printf("%s\n", bytes)
		}
	}

	//
	bytes, _ := json.MarshalIndent(info.App.Evars, "", "  ")
	fmt.Printf("%s\n", bytes)

	return nil
}
