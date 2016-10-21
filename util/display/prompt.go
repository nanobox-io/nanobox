package display

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// get the username
func ReadUsername() (string, error) {
	return Ask("Username")
}

func Ask(question string) (string, error) {
	fmt.Printf("%s: ", question)

	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')

	return strings.TrimSpace(str), err

}
