// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package mist

import (
	set "github.com/deckarep/golang-set"
	"sync"
)

type (
	localSubscriber struct {
		sync.Mutex

		check chan Message
		done  chan bool
		pipe  chan Message

		subscriptions []set.Set
		mist          *Mist
		id            uint32
	}
)

//
func NewLocalClient(mist *Mist, buffer int) *localSubscriber {
	client := &localSubscriber{
		check: make(chan Message, buffer),
		done:  make(chan bool),
		pipe:  make(chan Message),
		mist:  mist,
		id:    mist.nextId()}

	// this gofunc handles matching messages to subscriptions for the client
	go func(client *localSubscriber) {

		defer func() {
			close(client.check)
			close(client.pipe)
		}()

		for {
			select {
			case msg := <-client.check:
				// we do this so that we don't need a mutex
				subscriptions := client.subscriptions
				for _, subscription := range subscriptions {
					if subscription.IsSubset(msg.tags) {
						client.pipe <- msg
					}
				}
			case <-client.done:
				return
			}
		}
	}(client)

	// add the local client to mists list of subscribers
	mist.subscribers[client.id] = *client

	return client
}

//
func (client *localSubscriber) List() ([][]string, error) {
	subscriptions := make([][]string, len(client.subscriptions))
	for i, subscription := range client.subscriptions {
		sub := make([]string, subscription.Cardinality())
		for j, tag := range subscription.ToSlice() {
			sub[j] = tag.(string)
		}
		subscriptions[i] = sub
	}
	return subscriptions, nil
}

//
func (client *localSubscriber) Subscribe(tags []string) {
	subscription := makeSet(tags)

	client.Lock()
	client.subscriptions = append(client.subscriptions, subscription)
	client.Unlock()
}

// Unsubscribe iterates through each of mist clients subscriptions keeping all subscriptions
// that aren't the specified subscription
func (client *localSubscriber) Unsubscribe(tags []string) {
	client.Lock()

	//create a set for quick comparison
	test := makeSet(tags)

	// create a slice of subscriptions that are going to be kept
	keep := []set.Set{}

	// iterate over all of mist clients subscriptions looking for ones that match the
	// subscription to unsubscribe
	for _, subscription := range client.subscriptions {

		// if they are not the same set (meaning they are a different subscription) then add them
		// to the keep set
		if !test.Equal(subscription) {
			keep = append(keep, subscription)
		}
	}

	client.subscriptions = keep

	client.Unlock()
}

// Sends a message across mist
func (client *localSubscriber) Publish(tags []string, data string) error {
	client.mist.Publish(tags, data)
	return nil
}

//
func (client *localSubscriber) Ping() error {
	return nil
}

// Returns all messages that have sucessfully matched the list of subscriptions that this
// client has subscribed to
func (client *localSubscriber) Messages() <-chan Message {
	return client.pipe
}

//
func (client *localSubscriber) Close() error {
	// this closes the goroutine that is matching messages to subscriptions
	close(client.done)

	// remove the local client from mists list of subscribers
	delete(client.mist.subscribers, client.id)

	return nil
}

//
func makeSet(tags []string) set.Set {
	set := set.NewThreadUnsafeSet()
	for _, i := range tags {
		set.Add(i)
	}

	return set
}
