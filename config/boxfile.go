//
package config

import (
	"fmt"
	"os"
)

// BoxfileConfig represents all available/expected Boxfile configurable options
type BoxfileConfig struct {
	Build struct {
		Engine string `json:"engine"`
	}
}

// ParseBoxfile
func ParseBoxfile() (boxfile BoxfileConfig) {

	boxfilePath := "./Boxfile"

	// return early here with an empty boxfile; the error doesn't really matter because
	// that will be represented by the empty boxfile
	if _, err := os.Stat(boxfilePath); err != nil {
		return
	}

	//
	if err := ParseConfig(boxfilePath, &boxfile); err != nil {
		fmt.Printf("Nanobox failed to parse Boxfile. Please ensure it is valid YAML and try again.\n")
		Exit(1)
	}

	//
	return
}
