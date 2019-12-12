package build

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/dns"
)

func DevPayload(appModel *models.App) string {
	evars := appModel.Evars

	// create an APP_IP evar
	evars["APP_IP"] = appModel.LocalIPs["env"]

	// create a NANOBOX evar
	evars["NANOBOX"] = "nanobox"

	rtn := map[string]interface{}{}
	rtn["env"] = evars
	rtn["boxfile"] = appModel.DeployedBoxfile
	rtn["dns_entries"] = dns.List(" by nanobox")
	bytes, _ := json.Marshal(rtn)
	return string(bytes)
}
