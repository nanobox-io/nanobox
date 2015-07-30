package lumber

import (
	"testing"
)

func TestIsStar(t *testing.T) {

	log := NewConsoleLogger(FATAL)
	if !log.IsFatal() {
		t.Fatal("Fatal should be logged")
	}
	if log.IsError() {
		t.Fatal("Error should not be logged")
	}
	if log.IsWarn() {
		t.Fatal("Warn should not be logged")
	}
	if log.IsInfo() {
		t.Fatal("Info should not be logged")
	}
	if log.IsDebug() {
		t.Fatal("Debug should not be logged")
	}
	if log.IsTrace() {
		t.Fatal("Trace should not be logged")

	}

	log.Level(ERROR)
	if !log.IsFatal() {
		t.Fatal("Fatal should be logged")
	}
	if !log.IsError() {
		t.Fatal("Error should be logged")
	}
	if log.IsWarn() {
		t.Fatal("Warn should not be logged")
	}
	if log.IsInfo() {
		t.Fatal("Info should not be logged")
	}
	if log.IsDebug() {
		t.Fatal("Debug should not be logged")
	}
	if log.IsTrace() {
		t.Fatal("Trace should not be logged")

	}

	log.Level(WARN)
	if !log.IsFatal() {
		t.Fatal("Fatal should be logged")
	}
	if !log.IsError() {
		t.Fatal("Error should be logged")
	}
	if !log.IsWarn() {
		t.Fatal("Warn should be logged")
	}
	if log.IsInfo() {
		t.Fatal("Info should not be logged")
	}
	if log.IsDebug() {
		t.Fatal("Debug should not be logged")
	}
	if log.IsTrace() {
		t.Fatal("Trace should not be logged")
	}

	log.Level(INFO)
	if !log.IsFatal() {
		t.Fatal("Fatal should be logged")
	}
	if !log.IsError() {
		t.Fatal("Error should be logged")
	}
	if !log.IsWarn() {
		t.Fatal("Warn should be logged")
	}
	if !log.IsInfo() {
		t.Fatal("Info should be logged")
	}
	if log.IsDebug() {
		t.Fatal("Debug should not be logged")
	}
	if log.IsTrace() {
		t.Fatal("Trace should not be logged")
	}

	log.Level(DEBUG)
	if !log.IsFatal() {
		t.Fatal("Fatal should be logged")
	}
	if !log.IsError() {
		t.Fatal("Error should be logged")
	}
	if !log.IsWarn() {
		t.Fatal("Warn should be logged")
	}
	if !log.IsInfo() {
		t.Fatal("Info should be logged")
	}
	if !log.IsDebug() {
		t.Fatal("Debug should be logged")
	}
	if log.IsTrace() {
		t.Fatal("Trace should not be logged")
	}

	log.Level(TRACE)
	if !log.IsFatal() {
		t.Fatal("Fatal should be logged")
	}
	if !log.IsError() {
		t.Fatal("Error should be logged")
	}
	if !log.IsWarn() {
		t.Fatal("Warn should be logged")
	}
	if !log.IsInfo() {
		t.Fatal("Info should be logged")
	}
	if !log.IsDebug() {
		t.Fatal("Debug should be logged")
	}
	if !log.IsTrace() {
		t.Fatal("Trace should be logged")
	}
}

func TestMultiIS(t *testing.T) {
	log := NewMultiLogger()
	log.AddLoggers(NewConsoleLogger(WARN))
	log.AddLoggers(NewConsoleLogger(INFO))

	if log.IsTrace() {
		t.Fatal("Logger should return trace")
	}
	if !log.IsInfo() {
		t.Fatal("Logger should return info")
	}
	if !log.IsWarn() {
		t.Fatal("Logger should return warn")
	}
	if !log.IsError() {
		t.Fatal("Logger should return error")
	}
	if !log.IsFatal() {
		t.Fatal("Logger should return fatal")
	}
}
