//
package main

import (
	"github.com/nanopack/mist/commands"
)

func main() {
	if err := commands.MistCmd.Execute(); err != nil {
		return
	}
}
