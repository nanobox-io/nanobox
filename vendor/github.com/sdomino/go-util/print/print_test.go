package print_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	printutil "github.com/sdomino/go-util/print"
)

// TestColor
func TestColor(t *testing.T) {
	raw := "\x1b[31mtest color output\x1b[0m\x1b[0m\n"
	out := stdoutToString(func() { printutil.Color("[red]test color output[reset]") })
	if out != raw {
		t.Error(fmt.Sprintf("Expected '%q' got '%q'", raw, out))
	}
}

// TestPrompt
func TestPrompt(t *testing.T) {
}

// TestPassword
func TestPassword(t *testing.T) {
}

// stdoutToString
func stdoutToString(f func()) string {

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	//
	f()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	return string(out)
}
