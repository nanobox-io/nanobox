// Package commands ...
package commands

import (
	"fmt"
	"path/filepath"

	"github.com/jcelliott/lumber"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nanopack/mist/auth"
	"github.com/nanopack/mist/server"
)

var (
	host = "127.0.0.1:1445" // host clients will connect to
	tags []string           // tags to publish and [un]subscribe to/from

	config   string // location of the config file
	showVers bool   // whether to show version info and exit or not

	// to be populated by linker
	version string
	commit  string

	// MistCmd ...
	MistCmd = &cobra.Command{
		Use:           "mist",
		Short:         "Mist is a simple pub/sub for tagged messages",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		// parse the config if one is provided, or use the defaults
		PersistentPreRunE: readConfig,

		// print version or help, or continue, depending on flag settings
		PreRunE: preFlight,

		// either run as a server, or run as a CLI depending on what flags
		// are provided
		RunE: start,
	}
)

func readConfig(ccmd *cobra.Command, args []string) error {
	// if --version is passed print the version info
	if showVers {
		fmt.Printf("mist %s (%s)\n", version, commit)
		return fmt.Errorf("")
	}

	// if --config is passed, attempt to parse the config file
	if config != "" {
		filename := filepath.Base(config)
		viper.SetConfigName(filename[:len(filename)-len(filepath.Ext(filename))])
		viper.AddConfigPath(filepath.Dir(config))

		err := viper.ReadInConfig()
		if err != nil {
			fmt.Printf("ERROR: Failed to read config file: %s\n", err.Error())
			return err
		}
	}

	return nil
}

func preFlight(ccmd *cobra.Command, args []string) error {
	// if --server is not passed, print help
	if !viper.GetBool("server") {
		ccmd.HelpFunc()(ccmd, args)
		return fmt.Errorf("") // no error, just exit
	}

	return nil
}

func start(ccmd *cobra.Command, args []string) error {
	// configure the logger
	lumber.Prefix("[mist]")
	lumber.Level(lumber.LvlInt(viper.GetString("log-level")))

	if err := auth.Start(viper.GetString("authenticator")); err != nil {
		lumber.Fatal("Failed to start authenticator - %v", err)
		return err
	}

	if err := server.Start(viper.GetStringSlice("listeners"), viper.GetString("token")); err != nil {
		lumber.Fatal("One or more servers failed to start - %v", err)
		return err
	}

	return nil
}

func init() {

	// persistent config flags
	MistCmd.PersistentFlags().String("log-level", "INFO", "Output level of logs (TRACE, DEBUG, INFO, WARN, ERROR, FATAL)")
	viper.BindPFlag("log-level", MistCmd.PersistentFlags().Lookup("log-level"))

	MistCmd.PersistentFlags().String("token", "", "Auth token for connections")
	viper.BindPFlag("token", MistCmd.PersistentFlags().Lookup("token"))

	// local flags;
	MistCmd.Flags().String("authenticator", "", "Setting enables authentication, storing tokens in the authenticator provided")
	viper.BindPFlag("authenticator", MistCmd.Flags().Lookup("authenticator"))

	MistCmd.Flags().StringSlice("listeners", []string{"tcp://127.0.0.1:1445", "ws://127.0.0.1:8888"}, "A comma delimited list of servers to start")
	viper.BindPFlag("listeners", MistCmd.Flags().Lookup("listeners")) // no reason to have "http://127.0.0.1:8080" too, it only has /ping

	MistCmd.Flags().StringVar(&config, "config", config, "Path to config file")
	viper.BindPFlag("config", MistCmd.Flags().Lookup("config")) // seems this is only bound to access in server/server.go(BUG)

	MistCmd.Flags().Bool("server", false, "Run mist as a server")
	viper.BindPFlag("server", MistCmd.Flags().Lookup("server"))

	MistCmd.Flags().BoolVarP(&showVers, "version", "v", false, "Display the current version of this CLI")

	// commands
	MistCmd.AddCommand(pingCmd)
	MistCmd.AddCommand(subscribeCmd)
	MistCmd.AddCommand(publishCmd)

	// hidden/aliased commands
	MistCmd.AddCommand(listCmd)
	MistCmd.AddCommand(whoCmd)
	MistCmd.AddCommand(messageCmd)
	MistCmd.AddCommand(sendCmd)
}
