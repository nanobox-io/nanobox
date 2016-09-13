package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nanopack/mist/clients"
)

var (

	//
	listCmd = &cobra.Command{
		Hidden:        true,
		Use:           "listall",
		Short:         "List all subscriptions",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: list,
	}
)

// init
func init() {
	listCmd.Flags().StringVar(&host, "host", host, "The IP of a running mist server to connect to")
}

// list shows a unique list of all subscriptions subscribers are subscribed to
func list(ccmd *cobra.Command, args []string) error {

	// create new mist client
	client, err := clients.New(host, viper.GetString("token"))
	if err != nil {
		fmt.Printf("Failed to connect to '%v' - %v\n", host, err)
		return err
	}

	// listall related
	err = client.ListAll()
	if err != nil {
		fmt.Printf("Failed to list - %v\n", err)
		return err
	}

	msg := <-client.Messages()
	if msg.Data == "" {
		fmt.Printf("No subscribers connected to mist at '%v'\n", host)
	} else {
		fmt.Printf("Subscribers are subscribing on the following tags: %v\n", msg.Data)
	}

	return nil
}
