package code

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox/models"
)

type (
	fetch struct {
		Component      component      `json:"component"`
		LogvacHost     string         `json:"logvac_host"`
		Member         map[string]int `json:"member"`
		Build          string         `json:"build"`
		Warehouse      string         `json:"warehouse"`
		WarehouseToken string         `json:"warehouse_token"`
	}
)

// Fetch payload
func FetchPayload(componentModel *models.Component, warehouse string) string {

	logvac, _ := models.FindComponentBySlug(componentModel.AppID, "logvac")

	pload := fetch{
		LogvacHost: logvac.IPAddr(),
		Component: component{
			Name: componentModel.Name,
			UID:  componentModel.Name,
			ID:   componentModel.ID,
		},
		Member:         map[string]int{"uid": 1},
		Build:          "1234",
		Warehouse:      warehouse,
		WarehouseToken: "123",
	}

	bytes, err := json.Marshal(pload)
	if err != nil {
		return "{}"
	}

	return string(bytes)
}
