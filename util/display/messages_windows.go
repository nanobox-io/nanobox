// +build windows

package display

import (
  "fmt"
)

// prints a message to the user informing that only the command prompt is supported
func InvalidTerminal() {
  fmt.Println()
  fmt.Println("--------------------------------------------------------------")
  fmt.Println()
  fmt.Println("Oops, only the command prompt (cmd.exe) fully supports nanobox")
  fmt.Println()
  fmt.Println("--------------------------------------------------------------")
  fmt.Println()
}
