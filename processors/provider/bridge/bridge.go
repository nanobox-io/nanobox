package bridge

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/nanobox-io/nanobox/util/provider"
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

func bridgeConfig() string {
	node := ""
	if runtime.GOOS == "windows" {
		node = "dev-node MyTap"
	}

	ip, _ := provider.HostIP()
	return fmt.Sprintf(`client
%s
dev tap
proto udp
remote %s 1194
resolv-retry infinite
nobind
persist-key
persist-tun

ca %s
cert %s
key %s

cipher none
auth none
verb 3
`, node, ip, caCrt(), clientCrt(), clientKey())
}

func configFile() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "openvpn.conf"))
}

func caCrt() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "ca.crt"))
}

func clientKey() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "client.key"))
}

func clientCrt() string {
	return filepath.ToSlash(filepath.Join(config.EtcDir(), "openvpn", "client.crt"))
}