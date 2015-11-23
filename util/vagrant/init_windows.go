package vagrant

import (
	"os"
)

func sshLocation() string {
	return os.Getenv("HOMEDRIVE")+os.Getenv("HOMEDRIVE")+`\.ssh`
}