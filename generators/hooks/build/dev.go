package build

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox/models"
)

func DevPayload(appModel *models.App) string {
	rtn := map[string]interface{}{}
	rtn["env"] = appModel.Evars
	rtn["boxfile"] = appModel.DeployedBoxfile
	bytes, _ := json.Marshal(rtn)
	return string(bytes)
}
