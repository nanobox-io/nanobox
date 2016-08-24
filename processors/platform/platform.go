// Package platform ...
package platform

// these constants represent different potential states a platform can end up in
const (
	ACTIVE = "active"
)

// these represent different protocols that a platform might use
const (
	HTTP  = "http"
	HTTPS = "https"
	TCP   = "tcp"
	UDP   = "udp"
)

// Component ...
type PlatformComponent struct {
	label string
	name  string
	image string
}

// Components ...
var setupComponents = []PlatformComponent{
	{
		label: "Logger",
		name:  "logvac",
		image: "nanobox/logvac",
	},
	{
		label: "Message Bus",
		name:  "mist",
		image: "nanobox/mist",
	},
}

// Components ...
var deployComponents = []PlatformComponent{
	{
		label: "Router",
		name:  "portal",
		image: "nanobox/portal",
	},
	{
		label: "Storage",
		name:  "hoarder",
		image: "nanobox/hoarder",
	},
}
