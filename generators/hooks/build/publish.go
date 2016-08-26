package build

import (
	"encoding/json"
)

func PublishPayload(buildID, warehouse, warehouseToken, boxfile, previousBuild string) string {
	pload := map[string]interface{}{}
	if previousBuild != "" {
		pload["previous_build"] = previousBuild
	}
	pload["build"] = buildID
	pload["warehouse"] = warehouse
	pload["warehouse_token"] = warehouseToken
	pload["boxfile"] = boxfile
	b, _ := json.Marshal(pload)

	return string(b)	
}

