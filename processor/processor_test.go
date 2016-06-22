package processor_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/nanobox-io/nanobox/processor"
	_ "github.com/nanobox-io/nanobox/processor/code"
	_ "github.com/nanobox-io/nanobox/processor/platform"
	_ "github.com/nanobox-io/nanobox/processor/provider"
	_ "github.com/nanobox-io/nanobox/processor/service"
)

// testProccessor ...
type testProcessor struct {
	run bool
}

// TestMain ...
func TestMain(m *testing.M) {
	err := os.Chdir("../testing/")
	if err != nil {
		fmt.Println(err)
		return
	}
	processor.DefaultControl.Force = true
	processor.DefaultControl.Quiet = true
	// for testing we dont want to drop into a console
	// or hang on mist logging
	processor.Register("dev_console", testProcessBuilder)
	processor.Register("mist_log", testProcessBuilder)
	os.Exit(m.Run())
}

// Process ...
func (self testProcessor) Process() error {
	self.run = true
	return nil
}

// Results ....
func (self testProcessor) Results() processor.ProcessControl {
	return processor.ProcessControl{}
}

// TestRegister ...
func TestRegister(t *testing.T) {
	processor.Register("test", testProcessBuilder)
	err := processor.Run("test", processor.DefaultControl)
	if err != nil {
		t.Errorf("error from processor run", err)
	}
}

// TestBuild ...
func TestBuild(t *testing.T) {
	err := processor.Run("build", processor.DefaultControl)
	if err != nil {
		t.Errorf("error from build run", err)
	}
}

// TestDevDeploy ...
func TestDevDeploy(t *testing.T) {
	err := processor.Run("dev", processor.DefaultControl)
	if err != nil {
		t.Errorf("error from dev run", err)
	}
}

// TestDevDestroy ...
func TestDevDestroy(t *testing.T) {
	err := processor.Run("dev_destroy", processor.DefaultControl)
	if err != nil {
		t.Errorf("error from build run", err)
	}
}

// testProcessorBuilder ...
func testProcessBuilder(p processor.ProcessControl) (processor.Processor, error) {
	return testProcessor{}, nil
}
