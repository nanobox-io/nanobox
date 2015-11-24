// +build windows

package hosts

import "os"

var hostsPath = `C:\Windows\System32\drivers\etc\hosts`

func init() {
	hostsPath = os.Getenv("windir") + `\System32\drivers\etc\hosts`
}
