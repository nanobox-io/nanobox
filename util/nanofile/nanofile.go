package nanofile

import (
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/nanobox-io/nanobox/util"
)

var vip *viper.Viper

func Viper() *viper.Viper {
	if vip != nil {
		return vip
	}

	vip = viper.New()
	vip.SetDefault("external-network-space", "192.168.99.50/24")
	vip.SetDefault("internal-network-space", "192.168.0.50/16")
	vip.SetDefault("cpu-cap",  50)
	vip.SetDefault("cpus",  2)
	vip.SetDefault("host-dns", "off")
	vip.SetDefault("mount-nfs",  true)
	vip.SetDefault("provider", "virtualbox") // this may change in the future (adding additional hosts such as vmware
	vip.SetDefault("ram",  1024)
	vip.SetDefault("use-proxy",  false)

	vip.SetConfigFile(filepath.Join(util.GlobalDir(), "nanofile.yml"))
	vip.MergeInConfig() // using merge because it starts from existing config

	// we no longer use the local nanofile 
	// this is now in the boxfile under the 'dev' node
	// vip.SetConfigFile(filepath.Join(util.LocalDir(), "nanofile.yml"))
	// vip.MergeInConfig()
	return vip
}


