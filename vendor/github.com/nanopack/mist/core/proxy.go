package mist

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/jcelliott/lumber"
)

type (
	// Proxy ...
	Proxy struct {
		sync.RWMutex

		Authenticated bool
		Pipe          chan Message
		check         chan Message
		done          chan bool
		id            uint32
		subscriptions subscriptions
	}
)

// NewProxy ...
func NewProxy() (p *Proxy) {

	// create new proxy
	p = &Proxy{
		Pipe:          make(chan Message),
		check:         make(chan Message),
		done:          make(chan bool),
		id:            atomic.AddUint32(&uid, 1),
		subscriptions: newNode(),
	}

	p.connect()

	return
}

// connect
func (p *Proxy) connect() {

	lumber.Debug("Client connecting...")

	// this gofunc handles matching messages to subscriptions for the proxy
	go p.handleMessages()
}

//
func (p *Proxy) handleMessages() {

	defer func() {
		lumber.Trace("Got p.done, closing check and pipe")
		close(p.check)
		close(p.Pipe) // don't close pipe (response/pong messages need it), but leaving it unclosed leaves ram bloat on server even after client disconnects
	}()

	//
	for {
		select {

		// we need to ensure that this subscription actually has these tags before
		// sending anything to it; not doing this will cause everything to come
		// across the channel
		case msg := <-p.check:
			lumber.Trace("Got p.check")
			p.RLock()
			match := p.subscriptions.Match(msg.Tags)
			p.RUnlock()

			// if there is a subscription for the tags publish the message
			if match {
				lumber.Trace("Sending msg on pipe")
				p.Pipe <- msg
			}

		//
		case <-p.done:
			return
		}
	}
}

// Subscribe ...
func (p *Proxy) Subscribe(tags []string) {
	lumber.Trace("Proxy subscribing to '%v'...", tags)

	if len(tags) == 0 {
		return
	}

	// add proxy to subscribers list here so not all clients are 'subscribers'
	// since gets added to a map, there are no duplicates
	subscribe(p)

	// add tags to subscription
	p.Lock()
	p.subscriptions.Add(tags)
	p.Unlock()
}

// Unsubscribe ...
func (p *Proxy) Unsubscribe(tags []string) {
	lumber.Trace("Proxy unsubscribing from '%v'...", tags)

	if len(tags) == 0 {
		return
	}

	// remove tags from subscription
	p.Lock()
	p.subscriptions.Remove(tags)
	p.Unlock()
}

// Publish ...
func (p *Proxy) Publish(tags []string, data string) error {
	lumber.Trace("Proxy publishing to %v...", tags)

	return publish(p.id, tags, data)
}

// PublishAfter sends a message after [delay]
func (p *Proxy) PublishAfter(tags []string, data string, delay time.Duration) {
	go func() {
		<-time.After(delay)
		if err := publish(p.id, tags, data); err != nil {
			// log this error and continue
			lumber.Error("Proxy failed to PublishAfter - %v", err)
		}
	}()
}

// List returns a list of all current subscriptions
func (p *Proxy) List() (data [][]string) {
	lumber.Trace("Proxy listing subscriptions...")
	p.RLock()
	data = p.subscriptions.ToSlice()
	p.RUnlock()

	return
}

// Close ...
func (p *Proxy) Close() {
	lumber.Trace("Proxy closing...")

	if len(p.subscriptions.ToSlice()) != 0 {
		// remove the local p from mist's list of subscribers
		unsubscribe(p.id)
	}

	// this closes the goroutine that is matching messages to subscriptions
	close(p.done)
}
