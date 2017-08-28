// Package logvac handles the adding, removing, and writing to drains. It also
// defines the common types used accross logvac.
package logvac

import (
	"io"
	"sync"
	"time"

	"github.com/nanopack/logvac/config"
)

type (
	// Logger is a simple interface that's designed to be intentionally generic to
	// allow many different types of Logger's to satisfy its interface
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	// DEPRECATED: Use Message. OldMessage defines the structure of an old log message
	// I did what I hate most about docker, changed an exported struct definition. Sorry any client
	// using this.. at least I left it around?
	OldMessage struct {
		Time     time.Time `json:"time"`
		UTime    int64     `json:"utime"`
		Id       string    `json:"id"` // If setting multiple tags in id (syslog), set hostname first
		Tag      string    `json:"tag"`
		Type     string    `json:"type"` // Can be set if logs are submitted via http (deploy logs)
		Priority int       `json:"priority"`
		Content  string    `json:"message"`
	}

	// Message defines the structure of a log message
	Message struct {
		Time     time.Time `json:"time"`
		UTime    int64     `json:"utime"`
		Id       string    `json:"id"`   // ignoreifempty? // If setting multiple tags in id (syslog), set hostname first
		Tag      []string  `json:"tag"`  // ignoreifempty?
		Type     string    `json:"type"` // Can be set if logs are submitted via http (deploy logs)
		Priority int       `json:"priority"`
		Content  string    `json:"message"`
		Raw      []byte    `json:"raw,omitempty"`
	}

	// Logvac defines the structure for the default logvac object
	Logvac struct {
		drains map[string]drainChannels
	}

	// Drain defines a third party log drain endpoint (generally, only raw logs get drained)
	Drain struct {
		Type       string `json:"type"`             // type of service ("papertrail")
		URI        string `json:"endpoint"`         // uri of endpoint "log6.papertrailapp.com:199900"
		AuthKey    string `json:"key,omitempty"`    // key or user for authentication
		AuthSecret string `json:"secret,omitempty"` // password or secret for authentication
	}

	// DrainFunc is a function that "drains a Message"
	DrainFunc func(Message)

	drainChannels struct {
		send chan Message
		done chan bool
	}
)

// Vac is the default logvac object
var Vac Logvac

// Initializes a logvac object
func Init() error {
	Vac = Logvac{
		drains: make(map[string]drainChannels),
	}
	config.Log.Debug("Logvac initialized")
	return nil
}

// Close logvac and remove all drains
func Close() {
	Vac.close()
}

func (l *Logvac) close() {
	for tag := range l.drains {
		l.removeDrain(tag)
	}
}

// AddDrain adds a drain to the listeners and sets its logger
func AddDrain(tag string, drain DrainFunc) {
	Vac.addDrain(tag, drain)
}

func (l *Logvac) addDrain(tag string, drain DrainFunc) {
	channels := drainChannels{
		done: make(chan bool),
		send: make(chan Message),
	}

	go func() {
		for {
			select {
			case <-channels.done:
				return
			case msg := <-channels.send:
				// todo: ensure mist plays nice with goroutine
				go drain(msg)
			}
		}
	}()

	l.drains[tag] = channels
}

// RemoveDrain drops a drain
func RemoveDrain(tag string) {
	Vac.removeDrain(tag)
}

func (l *Logvac) removeDrain(tag string) {
	_, ok := l.drains[tag]
	if ok {
		close(l.drains[tag].done)
		delete(l.drains, tag)
	}
}

// WriteMessage broadcasts to all drains in seperate go routines
// Returns once all drains have received the message, but may not have processed
// the message yet
func WriteMessage(msg Message) {
	Vac.writeMessage(msg)
}

func (l *Logvac) writeMessage(msg Message) {
	// config.Log.Trace("Writing message - %s...", msg)
	group := sync.WaitGroup{}
	for _, drain := range l.drains {
		group.Add(1)
		go func(myDrain drainChannels) {
			select {
			case <-myDrain.done:
			case myDrain.send <- msg:
			}
			group.Done()
		}(drain)
	}
	group.Wait()
}

func (m Message) eof() bool {
	return len(m.Raw) == 0
}

func (m *Message) readByte() byte {
	// this function assumes that eof() check was done before
	b := m.Raw[0]
	m.Raw = m.Raw[1:]
	return b
}

func (m *Message) Read(p []byte) (n int, err error) {
	if m.eof() {
		err = io.EOF
		return
	}

	if c := cap(p); c > 0 {
		for n < c {
			p[n] = m.readByte()
			n++
			if m.eof() {
				break
			}
		}
	}
	return
}
