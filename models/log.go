package models

// LogOpts are options for logging
type LogOpts struct {
	Follow bool   // Follow is whether or not to follow the log stream.
	Number int    // Number of logs to print.
	Raw    bool   // Raw will not strip out the timestamp from the log stream.
	Start  string // Start is where to start the logs from.
	End    string // End is where to end the logs.
	Limit  string // Limit is how many logs to show.
}
