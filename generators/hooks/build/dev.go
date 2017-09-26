package build

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dns"
)

func DevPayload(appModel *models.App) string {
	// create an APP_IP evar
	evars := appModel.Evars
	evars["APP_IP"] = appModel.LocalIPs["env"]

	rtn := map[string]interface{}{}
	rtn["env"] = evars
	rtn["boxfile"] = appModel.DeployedBoxfile
	rtn["dns_entries"] = dns.List(" by nanobox")
	bytes, _ := json.Marshal(rtn)
	return string(bytes)
}
