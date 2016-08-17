// Package service ...
package component

// these constants represent all the "platform components" thats nanobox uses;
// platform components are just our other go microservcies
const (
	LOGVAC  = "logvac"
	PORTAL  = "portal"
	MIST    = "mist"
	HOARDER = "hoarder"
)

// these constants represent different potential states a service can end up in
const (
	ACTIVE = "active"
)
