package bridge

import (
	"fmt"
	"path/filepath"
	// "runtime"

	"github.com/nanobox-io/nanobox/util/provider"
	"github.com/nanobox-io/nanobox/util/config"
)

var BridgeClient = "nanobox-vpn"
// var BridgeURL string

// func init() {
// 	switch runtime.GOOS {
// 	case "windows":
// 		// BridgeClient = BridgeClient + ".exe"
// 		BridgeURL = "https://s3.amazonaws.com/tools.nanobox.io/openvpn/windows/openvpn.exe"
// 	case "darwin":
// 		BridgeURL = "https://s3.amazonaws.com/tools.nanobox.io/openvpn/darwin/openvpn"
// 	case "linux":
// 		BridgeURL = "https://s3.amazonaws.com/tools.nanobox.io/openvpn/linux/openvpn"
// 	}
// }

func BridgeConfig() string {
	// node := ""
	// if runtime.GOOS == "windows" {
	// 	node = "dev-node MyTap"
	// }

	ip, _ := provider.HostIP()
	return fmt.Sprintf(`client

dev tap
proto udp
remote %s 1194
resolv-retry infinite
nobind
persist-key
persist-tun

ca "%s"
cert "%s"
key "%s"

cipher none
auth none
verb 3
`, ip, CaCrt(), ClientCrt(), ClientKey())
}

func ConfigFile() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "openvpn.conf"))
}

func CaCrt() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "ca.crt"))
}

func ClientKey() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "client.key"))
}

func ClientCrt() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "client.crt"))
}