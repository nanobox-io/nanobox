package display

// Streamer ...
type Streamer struct {
	logLevel      string
	prefixer      *Prefixer
	captureOutput bool
	message       []byte
}

// NewStreamer returns a new Streamer
func NewStreamer(logLevel string) *Streamer {
	return &Streamer{
		logLevel: logLevel,
	}
}

// NewPrefixedStreamer returns a new Streamer with a Prefixer
func NewPrefixedStreamer(logLevel string, prefix string) Streamer {
	return Streamer{
		logLevel: logLevel,
		prefixer: NewPrefixer(prefix),
	}
}

// CaptureOutput will write messages to the message var.
func (s *Streamer) CaptureOutput(output bool) {
	s.captureOutput = output
}

// Write implements the io.Writer interface to write bytes on a writer
func (s *Streamer) Write(p []byte) (n int, err error) {
	msg := string(p)

	// todo: likely want to goroutine this in case there isn't a channel receiver
	if s.captureOutput {
		s.message = append(s.message, p...)
	}

	// if we have a prefixer run the message through it
	if s.prefixer != nil {
		msg = s.prefixer.Parse(msg)
	}

	switch s.logLevel {
	case "info":
		err = Info(msg)
	case "warn":
		err = Warn(msg)
	case "error":
		err = Error(msg)
	case "debug":
		err = Debug(msg)
	case "trace":
		err = Trace(msg)
	}

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (s Streamer) Output() string {
	return string(s.message)
}
