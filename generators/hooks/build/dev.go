package build

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox/models"
)

func DevPayload(appModel *models.App) string {
	rtn := map[string]interface{}{}
	rtn["env"] = appModel.Evars
	bytes, _ := json.Marshal(rtn)
	return string(bytes)
}
