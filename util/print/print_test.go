package print_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/nanobox-io/nanobox/config"
	printutil "github.com/nanobox-io/nanobox/util/print"
)

// TestVerboseOn
func TestVerboseOn(t *testing.T) {
	config.Verbose = true
	out := stdoutToString(func() { printutil.Verbose("test verbose output") })
	if out != "test verbose output" {
		t.Error(fmt.Sprintf("Expected 'test verbose output' got '%s'", out))
	}
}

// TestVerboseOff
func TestVerboseOff(t *testing.T) {
	config.Verbose = false
	out := stdoutToString(func() { printutil.Verbose("test verbose output") })
	if out != "" {
		t.Error(fmt.Sprintf("Expected nothing got '%s'", out))
	}
}

// TestSilenceOn
func TestSilenceOn(t *testing.T) {
	config.Silent = true
	out := stdoutToString(func() { printutil.Silence("test silence output") })
	if out != "" {
		t.Error(fmt.Sprintf("Expected nothing got '%s'", out))
	}
}

// TestSilenceOff
func TestSilenceOff(t *testing.T) {
	config.Silent = false
	out := stdoutToString(func() { printutil.Silence("test silence output") })
	if out != "test silence output" {
		t.Error(fmt.Sprintf("Expected 'test silence output' got '%s'", out))
	}
}

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
