package display

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// get the username
func ReadUsername() (string, error) {
	fmt.Print("Username: ")

	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')

	return strings.TrimSpace(str), err
}

