// Package mist ...
package mist

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jcelliott/lumber"
)

var (
	mutex       = &sync.RWMutex{}
	subscribers = make(map[uint32]*Proxy)
	uid         uint32
)

type (
	// A Message contains the tags used when subscribing, and the data that is being
	// published through mist
	Message struct {
		Command string   `json:"command"`
		Tags    []string `json:"tags,omitempty"`
		Data    string   `json:"data,omitempty"`
		Error   string   `json:"error,omitempty"`
	}

	// HandleFunc ...
	HandleFunc func(*Proxy, Message) error
)

// Subscribers is listall related
func Subscribers() string {
	subs := make(map[string]bool) // no duplicates

	// get tags all clients subscribed to
	for i := range subscribers {
		s := subscribers[i].subscriptions.ToSlice()
		for j := range s {
			for k := range s[j] {
				subs[s[j][k]] = true
			}
		}
	}

	// slice it
	subSlice := []string{}
	for k, _ := range subs {
		subSlice = append(subSlice, k)
	}

	return strings.Join(subSlice, " ")
}

// Who is who related
func Who() (int, int) {
	// subs := make(map[string]bool) // no duplicates
	subs := []string{}

	// get tags all clients subscribed to
	for i := range subscribers {
		subs = append(subs, fmt.Sprint(subscribers[i].id))
	}

	return len(subs), int(uid)
}

// todo: delete these 2. limiting what is a subscriber makes this not needed
// if they subscribe to a thing on a reused connection, they wanted to get updates.. hopefully
//
// Publish publishes to ALL subscribers. Usefull in client applications
// who reuse the publish connection for subscribing (publishes to self)
func Publish(tags []string, data string) error {
	lumber.Trace("Publishing...")
	return publish(0, tags, data)
}

// PublishAfter publishes to ALL subscribers. Usefull in client applications
// who reuse the publish connection for subscribing
func PublishAfter(tags []string, data string, delay time.Duration) error {
	go func() {
		<-time.After(delay)
		if err := Publish(tags, data); err != nil {
			// log this error and continue?
			lumber.Error("Failed to PublishAfter - %v", err)
		}
	}()

	return nil
}

// publish publishes to all subscribers except the one who issued the publish
func publish(pid uint32, tags []string, data string) error {

	if len(tags) == 0 {
		return fmt.Errorf("Failed to publish. Missing tags")
	}

	// if there are no subscribers, the message goes nowhere
	//
	// this could be more optimized, but it might not be an issue unless thousands
	// of clients are using mist.
	go func() {
		mutex.RLock()
		for _, subscriber := range subscribers {
			select {
			case <-subscriber.done:
				lumber.Trace("Subscriber done")
				// do nothing?

			default:

				// dont send this message to the publisher who just sent it
				if subscriber.id == pid {
					lumber.Trace("Subscriber is publisher, skipping publish")
					continue
				}

				// create message
				msg := Message{Command: "publish", Tags: tags, Data: data}

				// we don't want this operation blocking the range of other subscribers
				// waiting to get messages
				go func(p *Proxy, msg Message) {
					p.check <- msg
					lumber.Trace("Published message")
				}(subscriber, msg)
			}
		}
		mutex.RUnlock()
	}()

	return nil
}

// subscribe adds a proxy to the list of mist subscribers; we need this so that
// we can lock this process incase multiple proxies are subscribing at the same
// time
func subscribe(p *Proxy) {
	lumber.Trace("Adding proxy to subscribers...")

	mutex.Lock()
	subscribers[p.id] = p
	mutex.Unlock()
}

// unsubscribe removes a proxy from the list of mist subscribers; we need this
// so that we can lock this process incase multiple proxies are unsubscribing at
// the same time
func unsubscribe(pid uint32) {
	lumber.Trace("Removing proxy from subscribers...")

	mutex.Lock()
	delete(subscribers, pid)
	mutex.Unlock()
}
