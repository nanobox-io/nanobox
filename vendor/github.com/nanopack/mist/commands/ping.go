package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nanopack/mist/clients"
)

var (

	//
	pingCmd = &cobra.Command{
		Use:           "ping",
		Short:         "Ping a running mist server",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: ping,
	}
)

func init() {
	pingCmd.Flags().StringVar(&host, "host", host, "The IP of a running mist server to connect to")
}

// ping
func ping(ccmd *cobra.Command, args []string) error {

	client, err := clients.New(host, viper.GetString("token"))
	if err != nil {
		fmt.Printf("Failed to connect to '%v' - %v\n", host, err)
		return err
	}

	err = client.Ping()
	if err != nil {
		fmt.Printf("Failed to ping - %v\n", err)
		return err
	}

	msg := <-client.Messages()
	fmt.Println(msg.Data)

	return nil
}
