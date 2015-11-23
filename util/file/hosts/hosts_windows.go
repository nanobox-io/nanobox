// +build windows

package hosts

import "os"

var hostsPath = `C:\Windows\System32\etc\hosts`

func init() {
	hostsPath = os.Getenv("windir")+`\System32\etc\hosts`
}