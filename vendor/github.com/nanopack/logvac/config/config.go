// Package config is a central location for configuration options. It also contains
// config file parsing logic.
package config

import (
	"path/filepath"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// collectors
	ListenHttp = "127.0.0.1:6360" // address the api and http log collectors listen on
	ListenUdp  = "127.0.0.1:514"  // address the udp log collector listens on
	ListenTcp  = "127.0.0.1:6361" // address the tcp log collector listens on

	// drains
	PubAddress = ""                             // publisher address // mist://127.0.0.1:1445
	PubAuth    = ""                             // publisher auth token
	DbAddress  = "boltdb:///var/db/logvac.bolt" // database address

	// authenticator
	AuthAddress = "boltdb:///var/db/log-auth.bolt" // address or file location of auth backend ('boltdb:///var/db/logvac.bolt' or 'postgresql://127.0.0.1')

	// other
	CorsAllow = "*"            // sets `Access-Control-Allow-Origin` header
	LogKeep   = `{"app":"2w"}` // LogType and expire (X(m)in, (h)our,  (d)ay, (w)eek, (y)ear) (1, 10, 100 == keep up to that many) // todo: maybe map[string]interface
	LogType   = "app"          // default incoming log type when not set
	LogLevel  = "info"         // level which logvac will log at
	Token     = "secret"       // token to connect to logvac's api
	Log       lumber.Logger    // logger to write logs
	Insecure  = false          // whether or not to start insecure
	Server    = false          // whether or not to start logvac as a server
	Version   = false          // whether or not to print version info and exit
	CleanFreq = 60             // how often to clean log database
)

// AddFlags adds cli flags to logvac
func AddFlags(cmd *cobra.Command) {
	// collectors
	cmd.Flags().StringVarP(&ListenHttp, "listen-http", "a", ListenHttp, "API listen address (same endpoint for http log collection)")
	cmd.Flags().StringVarP(&ListenUdp, "listen-udp", "u", ListenUdp, "UDP log collection endpoint")
	cmd.Flags().StringVarP(&ListenTcp, "listen-tcp", "t", ListenTcp, "TCP log collection endpoint")

	// drains
	cmd.Flags().StringVarP(&PubAddress, "pub-address", "p", PubAddress, "Log publisher (mist) address (\"mist://127.0.0.1:1445\")")
	cmd.Flags().StringVarP(&PubAuth, "pub-auth", "P", PubAuth, "Log publisher (mist) auth token")
	cmd.Flags().StringVarP(&DbAddress, "db-address", "d", DbAddress, "Log storage address")

	// authenticator
	cmd.PersistentFlags().StringVarP(&AuthAddress, "auth-address", "A", AuthAddress, "Address or file location of authentication db. ('boltdb:///var/db/logvac.bolt' or 'postgresql://127.0.0.1')")

	// other
	cmd.Flags().StringVarP(&CorsAllow, "cors-allow", "C", CorsAllow, "Sets the 'Access-Control-Allow-Origin' header")
	cmd.Flags().StringVarP(&LogKeep, "log-keep", "k", LogKeep, "Age or number of logs to keep per type '{\"app\":\"2w\", \"deploy\": 10}' (int or X(m)in, (h)our,  (d)ay, (w)eek, (y)ear)")
	cmd.Flags().StringVarP(&LogLevel, "log-level", "l", LogLevel, "Level at which to log")
	cmd.Flags().StringVarP(&LogType, "log-type", "L", LogType, "Default type to apply to incoming logs (commonly used: app|deploy)")
	cmd.Flags().StringVarP(&Token, "token", "T", Token, "Administrative token to add/remove 'X-USER-TOKEN's used to pub/sub via http")
	cmd.Flags().BoolVarP(&Server, "server", "s", Server, "Run as server")
	cmd.Flags().BoolVarP(&Insecure, "insecure", "i", Insecure, "Don't use TLS (used for testing)")
	cmd.Flags().BoolVarP(&Version, "version", "v", Version, "Print version info and exit")
	cmd.Flags().IntVar(&CleanFreq, "clean-frequency", CleanFreq, "How often to clean log database")
	cmd.Flags().MarkHidden("clean-frequency")

	Log = lumber.NewConsoleLogger(lumber.LvlInt("ERROR"))
}

// ReadConfigFile reads in the config file, if any
func ReadConfigFile(configFile string) error {
	if configFile == "" {
		return nil
	}

	// Set defaults to whatever might be there already
	viper.SetDefault("listen-http", ListenHttp)
	viper.SetDefault("listen-udp", ListenUdp)
	viper.SetDefault("listen-tcp", ListenTcp)
	viper.SetDefault("pub-address", PubAddress)
	viper.SetDefault("pub-auth", PubAuth)
	viper.SetDefault("db-address", DbAddress)
	viper.SetDefault("auth-address", AuthAddress)
	viper.SetDefault("cors-allow", CorsAllow)
	viper.SetDefault("log-keep", LogKeep)
	viper.SetDefault("log-level", LogLevel)
	viper.SetDefault("log-type", LogType)
	viper.SetDefault("token", Token)
	viper.SetDefault("server", Server)
	viper.SetDefault("insecure", Insecure)

	filename := filepath.Base(configFile)
	viper.SetConfigName(filename[:len(filename)-len(filepath.Ext(filename))])
	viper.AddConfigPath(filepath.Dir(configFile))

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Set values. Config file will override commandline
	ListenHttp = viper.GetString("listen-http")
	ListenUdp = viper.GetString("listen-udp")
	ListenTcp = viper.GetString("listen-tcp")
	PubAddress = viper.GetString("pub-address")
	PubAuth = viper.GetString("pub-auth")
	DbAddress = viper.GetString("db-address")
	AuthAddress = viper.GetString("auth-address")
	CorsAllow = viper.GetString("cors-allow")
	LogKeep = viper.GetString("log-keep")
	LogLevel = viper.GetString("log-level")
	LogType = viper.GetString("log-type")
	Token = viper.GetString("token")
	Server = viper.GetBool("server")
	Insecure = viper.GetBool("insecure")

	return nil
}
