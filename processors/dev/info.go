package dev

import (
	"encoding/json"
	"fmt"

	"github.com/nanobox-io/nanobox/models"
)

// Info ...
func Info(appModel *models.App) error {

	//
	components, _ := models.AllComponentsByApp(appModel.ID)

	//
	for _, component := range components {
		if component.Name != "builds" {
			bytes, _ := json.MarshalIndent(component, "", "  ")
			fmt.Printf("%s\n", bytes)
		}
	}

	//
	bytes, _ := json.MarshalIndent(appModel.Evars, "", "  ")
	fmt.Printf("%s\n", bytes)

	return nil
}
