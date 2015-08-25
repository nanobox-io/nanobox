// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	set "github.com/deckarep/golang-set"
	"sync/atomic"
)

type (

	//
	Client interface {
		List() ([][]string, error)
		Subscribe(tags []string) error
		Unsubscribe(tags []string) error
		Publish(tags []string, data string) error
		Ping() error
		Messages() <-chan Message
		Close() error
	}

	//
	Mist struct {
		subscribers map[uint32]localSubscriber
		next        uint32
	}

	// A Message contains the tags used when subscribing, and the data that is being
	// published through mist
	Message struct {
		tags set.Set
		Tags []string `json:"tags"`
		Data string   `json:"data"`
	}
)

// creates a new mist
func New() *Mist {

	return &Mist{
		subscribers: make(map[uint32]localSubscriber),
	}
}

// Publish takes a list of tags and iterates through mist's list of subscribers,
// sending to each if they are available.
func (mist *Mist) Publish(tags []string, data string) {

	message := Message{
		Tags: tags,
		tags: makeSet(tags),
		Data: data,
	}

	for _, localSubscriber := range mist.subscribers {
		select {
		case <-localSubscriber.done:
		case localSubscriber.check <- message:
			// default:
			// do we drop the message? enqueue it? pull one off the front and then add this one?
		}
	}
}

//
func (mist *Mist) nextId() uint32 {
	return atomic.AddUint32(&mist.next, 1)
}
