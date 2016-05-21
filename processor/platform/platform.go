package platform

type PlatformService struct {
	label 	string
	name 		string
	image 	string
}

var (
	PlatformServices = []PlatformService{
		PlatformService{
			label: 	"Logger",
			name: 	"logvac",
			image: 	"nanobox/logvac",
		},
		PlatformService{
			label: 	"Router",
			name: 	"portal",
			image: 	"nanobox/portal",
		},
		PlatformService{
			label: 	"Message Bus",
			name: 	"mist",
			image: 	"nanobox/mist",
		},
		PlatformService{
			label: 	"Storage",
			name: 	"hoarder",
			image: 	"nanobox/hoarder",
		},
	}
)
