package bridge

import  (
	"github.com/nanobox-io/nanobox/util/service"
)

func StopService() error {
	return service.Stop("nanobox-vpn")
}


func Remove() error {
	return service.Stop("nanobox-vpn")
}
