package bridge

import (
	"path/filepath"
	"runtime"

	"github.com/nanobox-io/nanobox/util/config"
)

var bridgeClient = filepath.ToSlash(filepath.Join(config.BinDir(), "openvpn"))
var bridgeURL string

func init() {
	switch runtime.GOOS {
	case "windows":
		bridgeClient = bridgeClient + ".exe"
		bridgeURL = "https://s3.amazonaws.com/tools.nanobox.io/openvpn/windows/openvpn.exe"
	case "darwin":
		bridgeURL = "https://s3.amazonaws.com/tools.nanobox.io/openvpn/darwin/openvpn"
	case "linux":
		bridgeURL = "https://s3.amazonaws.com/tools.nanobox.io/openvpn/linux/openvpn"
	}
}
