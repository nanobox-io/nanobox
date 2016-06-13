// Package code ...
package code

import "errors"

// these constants represent different potential names an services can have
const (
	BUILD = "build"
)

// these constants represent different potential states an app can end up in
const (
	ACTIVE = "active"
)

var errMissingImageOrName = errors.New("missing image or name")
