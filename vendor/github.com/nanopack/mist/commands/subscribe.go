package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nanopack/mist/clients"
)

var (
	subscribeCmd = &cobra.Command{
		Use:           "subscribe",
		Short:         "Subscribe tags",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: subscribe,
	}
)

func init() {
	subscribeCmd.Flags().StringVar(&host, "host", host, "The IP of a running mist server to connect to")
	subscribeCmd.Flags().StringSliceVar(&tags, "tags", tags, "Tags to subscribe to")
}

// subscribe
func subscribe(ccmd *cobra.Command, args []string) error {

	// missing tags
	if tags == nil {
		fmt.Println("Unable to subscribe - Missing tags")
		return fmt.Errorf("")
	}

	client, err := clients.New(host, viper.GetString("token"))
	if err != nil {
		fmt.Printf("Failed to connect to '%v' - %v\n", host, err)
		return err
	}

	//
	if err := client.Subscribe(tags); err != nil {
		fmt.Printf("Unable to subscribe - %v\n", err.Error())
		return fmt.Errorf("")
	}

	// listen for messages on tags
	fmt.Printf("Listening on tags '%v'\n", tags)
	for msg := range client.Messages() {

		// skip handler messages
		if msg.Data != "success" {
			if viper.GetString("log-level") == "DEBUG" {
				fmt.Printf("Message: %#v\n", msg)
			} else {
				fmt.Println(msg.Data)
			}
		}
	}

	return nil
}
