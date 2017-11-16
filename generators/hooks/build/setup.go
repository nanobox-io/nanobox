package build

import (
	"encoding/json"
)

var ClearPkgCache bool

// SetupPayload returns a string for the user hook payload
func SetupPayload() string {
	if ClearPkgCache {
		rtn := map[string]string{}
		rtn["clear_cache"] = "true"
		bytes, _ := json.Marshal(rtn)
		return string(bytes)
	}
	// currently, this payload is empty. This may change at some point
	return emptyPayload()
}
