package processor_test

import (
	"testing"
	"os"
	"fmt"

	"github.com/nanobox-io/nanobox/processor"
	_ "github.com/nanobox-io/nanobox/processor/code"
	_ "github.com/nanobox-io/nanobox/processor/nanopack"
	_ "github.com/nanobox-io/nanobox/processor/provider"
	_ "github.com/nanobox-io/nanobox/processor/service"
)

type testProcessor struct {
	run bool
}

func TestMain(m *testing.M) {
	err := os.Chdir("../testing/")
	if err != nil {
		fmt.Println(err)
		return
	}
	processor.DefaultConfig.Force = true
	// for testing we dont want to drop into a console
	// or hang on mist logging
	processor.Register("Console", testProcessBuilder)
	processor.Register("mist_log", testProcessBuilder)
	os.Exit(m.Run())
}


func (self testProcessor) Process() error {
	self.run = true
	return nil
}

func (self testProcessor) Results() processor.ProcessConfig {
	return processor.ProcessConfig{}
}

func testProcessBuilder(p processor.ProcessConfig) (processor.Processor, error) {
	return testProcessor{}, nil
}

func TestRegister(t *testing.T) {
	processor.Register("test", testProcessBuilder)
	err := processor.Run("test", processor.DefaultConfig)
	if err != nil {
		t.Errorf("error from processor run", err)
	}
}

func TestBuild(t *testing.T) {
	err := processor.Run("build", processor.DefaultConfig)
	if err != nil {
		t.Errorf("error from build run", err)
	}
}

func TestDevDeploy(t *testing.T) {
	err := processor.Run("dev", processor.DefaultConfig)
	if err != nil {
		t.Errorf("error from dev run", err)
	}
}

// func TestDevDestroy(t *testing.T) {
// 	err := processor.Run("dev_destroy", processor.DefaultConfig)
// 	if err != nil {
// 		t.Errorf("error from build run", err)
// 	}
// }
