package platform

type PlatformService struct {
	label string
	name  string
	image string
}

var (
	PlatformServices = []PlatformService{
		{
			label: "Logger",
			name:  "logvac",
			image: "nanobox/logvac",
		},
		{
			label: "Router",
			name:  "portal",
			image: "nanobox/portal",
		},
		{
			label: "Message Bus",
			name:  "mist",
			image: "nanobox/mist",
		},
		{
			label: "Storage",
			name:  "hoarder",
			image: "nanobox/hoarder",
		},
	}
)
