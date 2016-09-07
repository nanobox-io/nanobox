package display

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/nanobox-io/nanobox/util/config"
)

var (
	// Log - enable logging to a file
	Log = true

	// LogFile - the location of logfile
	LogFile = filepath.ToSlash(filepath.Join(config.GlobalDir(), "process.log"))

	// Summary - summarize the output and hide log details
	Summary = true

	// Interactive - re-draw the summary when updates occur
	Interactive = terminal.IsTerminal(int(os.Stderr.Fd()))

	// Level - info, warn, error, debug, trace
	Level = "info"

	// Mode - text, json
	Mode = "text"

	// Out - writer to send output to
	Out = os.Stderr

	// internal
	logFile *os.File // open file descriptor of the log file
	// context
	context    int // track the context level
	topContext int // track the number to toplevel contexts
	// task
	taskStarted bool          // track if we're running a task
	taskLog     *bytes.Buffer // track the log of the current task, in case it fails
	prefixer    *Prefixer     // use a prefixer to prefix logs
	summarizer  *Summarizer   // summarizer to summarize the current task
)

// OpenContext opens a context level and prints the header
func OpenContext(format string, args ...interface{}) error {
	label := fmt.Sprintf(format, args...)

	// if the current context is 0, let's increment the topContext
	if context == 0 {
		topContext++
	}

	// increment the context level counter
	context++

	// if this is a subsequent top-level context, let's prefix with a newline
	if topContext > 1 && context == 1 {
		if err := printAll("\n"); err != nil {
			return err
		}
	}

	prefix := ""

	if context > 0 {
		prefix = strings.Repeat("  ", context-1)
	}

	header := fmt.Sprintf("%s+ %s :\n", prefix, label)

	if err := printAll(header); err != nil {
		return err
	}

	return nil
}

// CloseContext closes the context level and prints a newline
func CloseContext() error {

	// decrement the context level counter
	context--

	// ensure the context doesn't drop below zero
	if context < 0 {
		context = 0
	}

	return nil
}

// StartTask starts a new task
func StartTask(format string, args ...interface{}) error {
	label := fmt.Sprintf(format, args...)

	// return an error if the current task has not ended
	if taskStarted {
		return errors.New("Current task has not been stopped")
	}

	// mark the task as started
	taskStarted = true

	// initialize the task log
	taskLog = bytes.NewBufferString("")

	// create a new prefixer
	prefixer = NewPrefixer(strings.Repeat("  ", context+1))

	// generate a header
	prefix := strings.Repeat("  ", context)
	header := fmt.Sprintf("%s+ %s :\n", prefix, label)

	// print the header to the logfile
	if err := printLogFile(header); err != nil {
		return err
	}

	if Summary {
		summarizer = NewSummarizer(label, prefix)
		summarizer.Start()
	} else {
		// print the header
		if err := printOut(header); err != nil {
			return err
		}
	}

	return nil
}

// PauseTask ...
func PauseTask() {
	// stop the task summarizer
	if Summary && summarizer != nil {
		summarizer.Pause()
		fmt.Println()
	}
}

// ResumeTask ...
func ResumeTask() {
	// resume task
	if Summary && summarizer != nil {
		fmt.Println()
		summarizer.Resume()
	}
}

// StopTask stops the current task
func StopTask() error {

	// stop the task summarizer
	if Summary && summarizer != nil {
		summarizer.Stop()
		summarizer = nil
	}

	// mark task as stopped
	taskStarted = false

	// reset the task log
	taskLog = nil

	// reset the prefixer
	prefixer = nil

	return nil
}

// ErrorTask errors the current task
func ErrorTask() error {

	// stop the task summarizer
	if Summary && summarizer != nil {
		summarizer.Error()
		summarizer = nil

		// print the task log
		Out.Write(taskLog.Bytes())
	}

	// mark task as stopped
	taskStarted = false

	// reset the task log
	taskLog = nil

	// reset the prefixer
	prefixer = nil

	return nil
}

// Info sends an info level message to the current task
func Info(message string, args ...interface{}) error {
	if len(args) != 0 {
		message = fmt.Sprintf(message, args...)
	}

	// short-circuit if our log-level isn't high enough
	if currentLogLevel() > 2 {
		return nil
	}

	if err := log(message); err != nil {
		return err
	}

	return nil
}

// Warn sends a warn level message to the current task
func Warn(message string, args ...interface{}) error {
	if len(args) != 0 {
		message = fmt.Sprintf(message, args...)
	}

	// short-circuit if our log-level isn't high enough
	if currentLogLevel() > 3 {
		return nil
	}

	if err := log(message); err != nil {
		return err
	}

	return nil
}

// Error sends an error level message to the current task
func Error(message string, args ...interface{}) error {
	if len(args) != 0 {
		message = fmt.Sprintf(message, args...)
	}

	// short-circuit if our log-level isn't high enough
	if currentLogLevel() > 4 {
		return nil
	}

	if err := log(message); err != nil {
		return err
	}

	return nil
}

// Debug sends a debug level message to the current task
func Debug(message string, args ...interface{}) error {
	if len(args) != 0 {
		message = fmt.Sprintf(message, args...)
	}

	// short-circuit if our log-level isn't high enough
	if currentLogLevel() > 1 {
		return nil
	}

	if err := log(message); err != nil {
		return err
	}

	return nil
}

// Trace sends a trace level message to the current task
func Trace(message string, args ...interface{}) error {
	if len(args) != 0 {
		message = fmt.Sprintf(message, args...)
	}

	// short-circuit if our log-level isn't high enough
	if currentLogLevel() > 0 {
		return nil
	}

	if err := log(message); err != nil {
		return err
	}

	return nil
}

// log logs a message to the current task
func log(message string) error {

	// run the message through prefixer
	if prefixer != nil {
		message = prefixer.Parse(message)
	}

	// append to the taskLog
	if taskLog != nil {
		taskLog.WriteString(message)
	}

	// print message to logfile
	if err := printLogFile(message); err != nil {
		return err
	}

	if Summary && summarizer != nil {
		summarizer.Log(message)
	} else {
		// print the message
		if err := printOut(message); err != nil {
			return err
		}
	}

	return nil
}

// printAll prints a message to the Out channel and the logfile
func printAll(message string) error {

	// print to the Out writer
	if err := printOut(message); err != nil {
		return err
	}

	// print to the log file
	if err := printLogFile(message); err != nil {
		return err
	}

	return nil
}

// printOut will print a message to the out stream
func printOut(message string) error {
	_, err := io.WriteString(Out, message)
	return err
}

// printLogFile prints a message to the log file
func printLogFile(message string) error {
	// short-circuit if Log is set to false
	if !Log {
		return nil
	}

	// make sure the logfile is opened
	if err := openLogFile(); err != nil {
		return err
	}

	// print to the logfile
	logFile.WriteString(message)

	return nil
}

// openLogFile opens the logFile for writes
func openLogFile() error {

	// short-circuit if the logFile is already open
	if logFile != nil {
		return nil
	}

	truncate := os.O_RDWR | os.O_CREATE | os.O_TRUNC

	f, err := os.OpenFile(LogFile, truncate, 0644)
	if err != nil {
		return err
	}

	logFile = f

	return nil
}

// currentLogLevel returns the current log level as an int
func currentLogLevel() int {
	switch Level {
	case "error":
		return 4
	case "warn":
		return 3
	case "info":
		return 2
	case "debug":
		return 1
	case "trace":
		return 0
	}

	return 0
}
