package models

import "fmt"

var (
	// will be set with build flags, defaults for one-off `go-build`
	nanoVersion string = "0.0.0"  // git tag
	nanoCommit  string = "custom" // commit id of build
	nanoBuild   string = "now"    // date of build
)

func VersionString() string {
	return fmt.Sprintf("Nanobox Version %s-%s (%s)", nanoVersion, nanoBuild, nanoCommit)
}
