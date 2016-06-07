// Package print ...
package print

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/mitchellh/colorstring"
)

// Color wraps a print message in 'colorstring' and passes it to fmt.Println
func Color(msg string, v ...interface{}) {
	fmt.Println(colorstring.Color(fmt.Sprintf(msg, v...)))
}

// Password prompts for a password but keeps the typed response hidden
func Password(p string) string {
	fmt.Printf(p)
	return string(gopass.GetPasswd())
}

// Prompt will prompt for input from the shell and return a trimmed response
func Prompt(p string, v ...interface{}) string {
	reader := bufio.NewReader(os.Stdin)

	//
	fmt.Print(colorstring.Color(fmt.Sprintf(p, v...)))

	//
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "[util/print] reader.ReadString() failed - %v", err.Error())
	}

	//
	return strings.TrimSpace(input)
}
