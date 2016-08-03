package print

import (
	"fmt"
	"io"
)

//
type progress struct {
	prefix string
}

// Write ...
func (s progress) Write(p []byte) (n int, err error) {
	fmt.Printf("%s%s", s.prefix, p)
	return len(p), nil
}

// NewProgress executes a pre-assembled command and streams the output with a prefix
func NewProgress(prefix string) io.Writer {
	return progress{prefix}
}
