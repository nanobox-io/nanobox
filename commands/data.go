package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/util/data"
)

type (
	anything interface{

	}
)
	
var (


	// DataCmd ...
	DataCmd = &cobra.Command{
		Use:   "data",
		Short: "show element from the nanobox database",
		Long:  ``,
		Run:   dataFunc,
		Hidden: true,
	}
)

// dataFunc ...
func dataFunc(ccmd *cobra.Command, args []string) {
	switch {
	default:
		fmt.Println("I need to know some data starting point")

	case len(args) == 1:
		keys, err := data.Keys(args[0])
		if err != nil {
			fmt.Println(err)
		}
		for _, key := range keys {
			showData(args[0], key)
		}
	case len(args) == 2:
		showData(args[0], args[1])
	}
}

func showData(bucket, key string) {
	i := map[string]interface{}{}
	err := data.Get(bucket, key, &i)
	if err != nil {
		a := []interface{}{}
		err := data.Get(bucket, key, &a)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("key: %s, val: %+v\n", key, a)
		return
	}
	fmt.Printf("key: %s, val: %+v\n", key, i)
}