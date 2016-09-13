package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/nanopack/mist/clients"
)

var (

	// alias for publish
	messageCmd = &cobra.Command{
		Hidden: true,

		Use:           "message",
		Short:         "Publish a message",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: publish,
	}

	publishCmd = &cobra.Command{
		Use:           "publish",
		Short:         "Publish a message",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: publish,
	}

	// alias for publish
	sendCmd = &cobra.Command{
		Hidden: true,

		Use:           "send",
		Short:         "Publish a message",
		Long:          ``,
		SilenceErrors: true,
		SilenceUsage:  true,

		RunE: publish,
	}
)

var data string //

// init
func init() {
	publishCmd.Flags().StringVar(&host, "host", host, "The IP of a running mist server to connect to")
	messageCmd.Flags().StringVar(&host, "host", host, "The IP of a running mist server to connect to")
	sendCmd.Flags().StringVar(&host, "host", host, "The IP of a running mist server to connect to")

	publishCmd.Flags().StringVar(&data, "data", data, "The string data to publish")
	messageCmd.Flags().StringVar(&data, "data", data, "The string data to message")
	sendCmd.Flags().StringVar(&data, "data", data, "The string data to send")

	publishCmd.Flags().StringSliceVar(&tags, "tags", tags, "Tags to publish to")
	messageCmd.Flags().StringSliceVar(&tags, "tags", tags, "Tags to publish to")
	sendCmd.Flags().StringSliceVar(&tags, "tags", tags, "Tags to publish to")
}

// publish
func publish(ccmd *cobra.Command, args []string) error {

	// missing tags
	if tags == nil {
		fmt.Println("Unable to publish - Missing tags")
		return fmt.Errorf("")
	}

	// missing data
	if data == "" {
		fmt.Println("Unable to publish - Missing data")
		return fmt.Errorf("")
	}

	client, err := clients.New(host, viper.GetString("token"))
	if err != nil {
		fmt.Printf("Failed to connect to '%v' - %v\n", host, err)
		return err
	}

	err = client.Publish(tags, data)
	if err != nil {
		fmt.Printf("Failed to publish message - %v\n", err)
		return err
	}

	fmt.Println("success")

	return nil
}
