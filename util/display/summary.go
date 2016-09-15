package display

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/nanobox-io/nanobox/util"
)

// ...
var (
	EscSeqRegex   = regexp.MustCompile("\\x1b[[][?0123456789]*;?[?0123456789]*[ABCDEFGHJKRSTfminsulhp]")
	LogStripRegex = regexp.MustCompile("^[ \t-]*")
	termWidth     = 0
)

type (

	// Summarizer ...
	Summarizer struct {
		Label  string    // the task label to print as the header
		Prefix string    // the prefix to prepend to the summary
		Out    io.Writer // writer to send output to

		// internal
		chEvent     chan *sEventOp // channel to receive stop/error/tick events
		chLog       chan string    // channel to receive logs
		spinIdx     int            // track the current spinner index
		ticker      *time.Ticker   // timer to send tick events at an interval
		detail      string         // the current line of detail to print
		shutdown    bool           // toggle to inform the run loop to exit
		windowWidth int
	}

	// Sending events to the summarizer needs to block the caller until
	// the event handler has completed. To do this we need to create
	// a stateful event with an ad-hoc channel to receive the response
	sEventOp struct {
		action string    // the requested action (stop/error/tick/etc)
		res    chan bool // the response channel to syncronize state
	}
)

// NewSummarizer returns a new Summarizer
func NewSummarizer(label string, prefix string) *Summarizer {
	if termWidth == 0 {
		_, termWidth = util.GetTerminalSize()
	}

	return &Summarizer{
		Label:  label,
		Prefix: prefix,
		Out:    os.Stderr,

		chEvent:     make(chan *sEventOp),   // no buffering, block the sending process
		chLog:       make(chan string, 100), // buffer up to 100 log messages before blocking
		windowWidth: termWidth,
	}
}

// Start starts the summary process in a goroutine
func (s *Summarizer) Start() {
	go s.run()
}

// Pause ...
func (s *Summarizer) Pause() {
	// create the 'stop' event
	event := &sEventOp{
		action: "pause",
		res:    make(chan bool),
	}

	// send the event
	s.chEvent <- event

	// now wait until we get a response back
	<-event.res
}

// Resume prints the "complete" label and toggles shutdown
func (s *Summarizer) Resume() {
	// generate and print the complete header
	s.Label = s.Label

	// turn the ticker back on
	s.shutdown = false
	go s.run()
}

// Stop stops the summary process
func (s *Summarizer) Stop() {
	// create the 'stop' event
	event := &sEventOp{
		action: "stop",
		res:    make(chan bool),
	}

	// send the event
	s.chEvent <- event

	// now wait until we get a response back
	<-event.res
}

// Error will stop the summary process and print an error header
func (s *Summarizer) Error() {
	// create the 'error' event
	event := &sEventOp{
		action: "error",
		res:    make(chan bool),
	}

	// send the event
	s.chEvent <- event

	// now wait until we get a response back
	<-event.res
}

// Log sends a log message to the summary process
func (s *Summarizer) Log(msg string) {
	// send the message to the log channel
	s.chLog <- msg
}

// run runs the main loop to wait for events or timers
func (s *Summarizer) run() {

	// to kick things off, we need to print the first summary
	s.print()

	// start a ticker to send a "tick" event every 80 milliseconds
	s.startTicker()

	for !s.shutdown {
		select {
		case event := <-s.chEvent:
			s.handleEvent(event)
		case msg := <-s.chLog:
			s.handleLog(msg)
		}
	}

	// stop the ticker
	s.stopTicker()
}

// startTicker starts a ticker to send a tick event
func (s *Summarizer) startTicker() {
	s.ticker = time.NewTicker(time.Millisecond * 80)

	// send a "tick" event at each interval
	go func() {
		for range s.ticker.C {
			// create the 'tick' event
			event := &sEventOp{
				action: "tick",
				res:    make(chan bool),
			}

			// send the event
			s.chEvent <- event

			// now wait until we get a response back
			<-event.res
		}
	}()
}

// stopTicker stops the tick event ticker
func (s *Summarizer) stopTicker() {
	s.ticker.Stop()
}

// handleEvent dispatches actions for events received on the event channel
func (s *Summarizer) handleEvent(event *sEventOp) {

	switch event.action {
	case "tick":
		s.tick()
	case "stop":
		s.stop()
	case "pause":
		s.pause()
	case "error":
		s.error()
	}

	// send the response so the caller can continue on
	event.res <- true
}

// handleLog sets the detail line and refreshes the summary
func (s *Summarizer) handleLog(data string) {

	msg := ""

	// iterate through the lines, we'll keep the last line that has data
	for _, line := range strings.Split(data, "\n") {
		// first we need to remove escape sequences
		line = EscSeqRegex.ReplaceAllString(line, "")

		// then we need to remove any leading whitespace
		line = LogStripRegex.ReplaceAllString(line, "")

		if len(line) > 0 {
			msg = line
		}

		if s.windowWidth > 0 && len(line) > s.windowWidth {
			msg = msg[:s.windowWidth-(len(s.Label)+len(s.Prefix)+10)] + "..."
		}
	}

	if len(msg) > 0 {
		s.detail = msg
		s.reset()
		s.print()
	}
}

// tick updates the spinner index and refreshes the summary
func (s *Summarizer) tick() {
	// update the spinner index
	s.spinIdx++

	// reset the index back to 0 if we've reached the top
	if s.spinIdx == len(TaskSpinner) {
		s.spinIdx = 0
	}

	// reset and print the screen
	s.reset()
	s.print()
}

// stop prints the "complete" label and toggles shutdown
func (s *Summarizer) stop() {
	// reset the screen
	s.reset()

	// generate and print the complete header
	header := fmt.Sprintf("%s%s %s\n", s.Prefix, TaskComplete, s.Label)
	io.WriteString(s.Out, header)

	// set the shutdown flag to ensure the loop ends
	s.shutdown = true
}

// stop prints the "complete" label and toggles shutdown
func (s *Summarizer) pause() {
	// reset the screen
	s.reset()

	// generate and print the complete header
	header := fmt.Sprintf("%s%s %s\n", s.Prefix, TaskPause, s.Label)
	io.WriteString(s.Out, header)

	// set the shutdown flag to ensure the loop ends
	s.shutdown = true
}

// error prints the "error" label and toggles shutdown
func (s *Summarizer) error() {
	// reset the screen
	s.reset()

	// generate and print the complete header
	header := fmt.Sprintf("%s! %s\n", s.Prefix, s.Label)
	io.WriteString(s.Out, header)

	// set the shutdown flag to ensure the loop ends
	s.shutdown = true
}

// reset will reset the screen using escape sequences
func (s *Summarizer) reset() {
	// http://bluesock.org/~willg/dev/ansi.html

	// todo: make this conditional on the progress estimator
	lines := 2

	// create escape sequence to move up a line and clear the line for each line
	reset := strings.Repeat("\x1b[1A\x1b[K", lines)

	io.WriteString(s.Out, reset)
}

// print prints the current summary
func (s *Summarizer) print() {

	header := fmt.Sprintf("%s%s %s :\n", s.Prefix, TaskSpinner[s.spinIdx], s.Label)
	detail := fmt.Sprintf("%s  %s\n", s.Prefix, s.detail)

	// todo: add progress estimator

	io.WriteString(s.Out, header)
	io.WriteString(s.Out, detail)
}
