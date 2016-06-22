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

// Service ...
type Service struct {
	label string
	name  string
	image string
}

// Services ...
var SetupServices = []Service{
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

// Services ...
var DeployServices = []Service{
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
