package build

import (
	"encoding/json"

	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/models"
)
type (
	mount struct {
		Host     string   `json:"host"`
		Protocol string   `json:"protocol"`
		Shares   []string `json:"shares"`
	}
)


func DevPayload(appModel *models.App) string {
	rtn := map[string]interface{}{}
	rtn["env"] = appModel.Evars
	rtn["mounts"] = mounts(appModel)
	bytes, _ := json.Marshal(rtn)
	return string(bytes)
}


// mounts ...
func mounts(appModel *models.App) []mount {
	boxfile := boxfile.New([]byte(appModel.DeployedBoxfile))
	boxNetworkDirs := boxfile.Node("dev").Node("network_dirs")

	m := []mount{}
	for _, node := range boxNetworkDirs.Nodes() {
		// i think i store these as data.name
		// cleanNode := regexp.MustCompile(`.+\.`).ReplaceAllString(node, "")
		component, err := models.FindComponentBySlug(appModel.ID, node)
		if err != nil {
			// skip because of problems
			continue
		}
		if !component.Plan.BehaviorPresent("mountable") || component.Plan.MountProtocol == "" {
			// skip because of problems
			continue
		}
		m = append(m, mount{component.InternalIP, component.Plan.MountProtocol, boxNetworkDirs.StringSliceValue(node)})

	}

	return m
}
