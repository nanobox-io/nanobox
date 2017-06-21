package build

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox/models"
)

func DevPayload(appModel *models.App) string {
	// create an APP_IP evar
	evars := appModel.Evars
	evars["APP_IP"] = appModel.LocalIPs["env"]

	rtn := map[string]interface{}{}
	rtn["env"] = evars
	rtn["boxfile"] = appModel.DeployedBoxfile
	bytes, _ := json.Marshal(rtn)
	return string(bytes)
}
