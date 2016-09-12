This ssh package contains helpers for working with ssh in go.  The `client.go` file
is a modified version of `docker/machine/libmachine/ssh/client.go` that only
uses golang's native ssh client. It has also been improved to resize the tty as
needed. The key functions are meant to be used by either client or server
and will generate/store keys if not found.

## Usage:

```go
package main

import (
	"fmt"

	"github.com/glinton/ssh"
)

func main() {
	err := connect()
	if err != nil {
		fmt.Printf("Failed to connect - %s\n", err)
	}
}

func connect() error {
	nanPass := ssh.Auth{Passwords: []string{"pass"}}
	client, err := ssh.NewNativeClient("user", "localhost", "SSH-2.0-MyCustomClient-1.0", 2222, &nanPass)
	if err != nil {
		return fmt.Errorf("Failed to create new client - %s", err)
	}

	err = client.Shell()
	if err != nil && err.Error() != "exit status 255" {
		return fmt.Errorf("Failed to request shell - %s", err)
	}

	return nil
}
```
