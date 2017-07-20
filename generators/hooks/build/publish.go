package build

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox/models"
)

type WarehouseConfig struct {
	BuildID        string
	WarehouseURL   string
	WarehouseToken string
	PreviousBuild  string
}

func PublishPayload(envModel *models.Env, warehouseConfig WarehouseConfig) string {
	config, _ := models.LoadConfig()

	pload := map[string]interface{}{
		"build":           warehouseConfig.BuildID,
		"warehouse":       warehouseConfig.WarehouseURL,
		"warehouse_token": warehouseConfig.WarehouseToken,
		"boxfile":         envModel.BuiltBoxfile,
		"sync_verbose":    !(config.CIMode && !config.CISyncVerbose),
	}

	if warehouseConfig.PreviousBuild != "" {
		pload["previous_build"] = warehouseConfig.PreviousBuild
	}

	b, _ := json.Marshal(pload)

	return string(b)
}
