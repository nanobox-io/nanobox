package print

import (
	"fmt"
	"io"
)

//
type streamer struct {
	prefix string
}

// Write ...
func (s streamer) Write(p []byte) (n int, err error) {
	fmt.Printf("%s%s", s.prefix, p)
	return len(p), nil
}

// NewStreamer executes a pre-assembled command and streams the output with a prefix
func NewStreamer(prefix string) io.Writer {
	return streamer{prefix}
}
