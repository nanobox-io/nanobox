// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	"fmt"
	"sync"

	"github.com/pagodabox/golang-hatchet"
)

//
const (
	DefaultPort = "1445"
	Version     = "0.1.0"
)

//
type (

	//
	Mist struct {
		sync.Mutex

		log           hatchet.Logger //
		port          string         //
		Subscriptions []Subscription //
	}

	// A Subscription has tags which are are used to match when publishing messages,
	// and a channel that receives those messages
	Subscription struct {
		Tags []string     `json:"tags"`    // the tags used to subscribe to published messages
		Sub  chan Message `json:"channel"` // the channel which published message data comes across
	}

	// A Message contains the tags used when subscribing, and the data that is being
	// published to the mist server
	Message struct {
		Tags []string `json:"tags"`        // the tags used to subscribe to updates
		Data string   `json:"data,string"` // the actual content of the message
	}
)

// creates a new mist, setting up a logger, and starting a mist server
func New(port string, logger hatchet.Logger) *Mist {

	//
	if logger == nil {
		logger = hatchet.DevNullLogger{}
	}

	// if no port is specified, use mists default port (1445)
	if port == "" {
		port = DefaultPort
	}

	//
	mist := &Mist{
		log:           logger,
		port:          port,
		Subscriptions: []Subscription{},
	}

	// start a mist server
	mist.start()

	return mist
}

// Publish takes a list of tags and iterates through mist's list of subscriptions,
// looking for matching subscriptions to publish messages too.
func (m *Mist) Publish(tags []string, data string) {

	// interate over each subscription that mist has, checking to see if there are
	// any that match the tags to publish on
	for _, s := range m.Subscriptions {

		// if the subscription contains the tags that are being published to, publish
		// the message across the subscriptions channel, this can be a lose check
		// because a message can still be published if the publish has more tags than
		// the subscription (but not the other way around)
		if contains(s.Tags, tags, false) {
			go func(ch chan Message, msg Message) { ch <- msg }(s.Sub, Message{Tags: tags, Data: data})
		}
	}
}

// Subscribe takes a subscription and appends it to mists list of subscriptions
func (m *Mist) Subscribe(sub Subscription) {
	m.Lock()
	m.Subscriptions = append(m.Subscriptions, sub)
	m.Unlock()
}

// Unsubscribe iterates through each of mists subscriptions keeping all subscriptions
// that aren't the specified subscription
func (m *Mist) Unsubscribe(sub Subscription) {
	m.Lock()

	// create a slice of subscriptions that are going to be kept
	keep := []Subscription{}

	// iterate over all of mists subscriptions looking for ones that match the
	// subscription to unsubscribe
	for _, s := range m.Subscriptions {

		// if the tags do not match add the subscrpition to the list of subscriptions
		// to be kept, this must be a strict match because a subscription should only
		// be unsubscribed if there is an exact tag match (not order)
		if !contains(s.Tags, sub.Tags, true) {
			keep = append(keep, sub)
		}
	}

	// set mists subscriptions equal to the remaining subscriptions
	m.Subscriptions = keep

	m.Unlock()
}

// List displays a list of mists current subscriptions
func (m *Mist) List() {
	fmt.Println(m.Subscriptions)
}

// contains takes two sets of tags, and compaires them to see if the first set
// (needle) is found in the second set (haystack)
func contains(needle, haystack []string, strict bool) bool {

	// if a strict comparison is desired, return false if the lengths dont match
	if strict {
		if len(needle) != len(haystack) {
			return false
		}
	}

	//
	tags := map[string]interface{}{}

	// create a map of all the tags that are in the set to be compaired against
	for _, t := range haystack {
		tags[t] = t
	}

	// interate through eatch tag in the 'needle' to see if its contained in the
	// 'haystack', if it is no found return false otherwise the loop eventually
	// completes and returns true
	for _, t := range needle {
		if _, ok := tags[t]; !ok {
			return false
		}
	}

	return true
}
