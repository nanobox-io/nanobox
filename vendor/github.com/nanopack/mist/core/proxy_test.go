package mist

import "testing"

// TestSameSubscriber tests to ensure that mist will not send message to the
// same proxy who publishes them
func TestSameSubscriber(t *testing.T) {

	// create a new proxy
	sender := NewProxy()
	defer sender.Close()

	// sender subscribes to tags and then tries to publish to those same tags,
	// verifying that a message is receieved...
	sender.Subscribe([]string{"a"})
	defer sender.Unsubscribe([]string{"a"})

	// sender published to the same tags it subscribed to then waits for a message
	// that should never come because mist shouldnt send a message to the same proxy
	// that publishes them. After 1 second assume no message is coming.
	sender.Publish([]string{"a"}, testMsg)
	verifyNoMessage(sender, t)
}

// TestDifferentSubscriber tests to ensure that mist will send messages
// to another subscribed proxy, and then not send when unsubscribed.
func TestDifferentSubscriber(t *testing.T) {

	//
	sender := NewProxy()
	defer sender.Close()

	//
	receiver := NewProxy()
	defer receiver.Close()

	// receiver subscribes to tags and then sender publishes to those tags,
	// verifying that a message is receieved...
	receiver.Subscribe([]string{"a"})
	sender.Publish([]string{"a"}, testMsg)
	verifyMessage(testMsg, receiver, t)

	// receiver unsubscribes from the tags and sender publishes again to the same
	// tags, verifying that NO message is receieved...
	receiver.Unsubscribe([]string{"a"})
	sender.Publish([]string{"a"}, testMsg)
	verifyNoMessage(receiver, t)
}

// TestManySubscribers tests to ensure that mist will send messages to many
// subscribers of the same tags, and then not send once unsubscribed
func TestManySubscribers(t *testing.T) {

	//
	sender := NewProxy()
	defer sender.Close()

	//
	r1 := NewProxy()
	defer r1.Close()

	//
	r2 := NewProxy()
	defer r2.Close()

	//
	r3 := NewProxy()
	defer r3.Close()

	// receivers subscribe to tags and then sender publishes to those tags, verifying
	// that messages are received...
	r1.Subscribe([]string{"a"})
	r2.Subscribe([]string{"a"})
	r3.Subscribe([]string{"a"})
	sender.Publish([]string{"a"}, testMsg)
	verifyMessage(testMsg, r1, t)
	verifyMessage(testMsg, r2, t)
	verifyMessage(testMsg, r3, t)

	// receiver unsubscribes from the tags and sender publishes again to the same
	// tags, veriftying that messages are NOT receieved...
	r1.Unsubscribe([]string{"a"})
	r2.Unsubscribe([]string{"a"})
	r3.Unsubscribe([]string{"a"})
	sender.Publish([]string{"a"}, testMsg)
	verifyNoMessage(r1, t)
	verifyNoMessage(r2, t)
	verifyNoMessage(r3, t)
}

// TestListSubscriptions tests to ensure that mist will list only the current
// subscriptions
func TestListSubscriptions(t *testing.T) {

	//
	sender := NewProxy()
	defer sender.Close()

	var list string

	// sender subscribes to a single tag and then lists it's tags verifying it
	sender.Subscribe([]string{"a"})
	list = flattenSliceToString(sender.List())
	if list != "a" {
		t.Errorf("Unexpected tags - Expecting '%v' received %v", "a", list)
	}

	// sender subscribes to another tag and then lists it's tags verifying multiple
	// tags; we test both configurations here because maps are unordered
	sender.Subscribe([]string{"b"})
	list = flattenSliceToString(sender.List())
	switch list {
	case "ab", "ba":
		// pass
	default:
		t.Errorf("Unexpected tags - Expecting '%v' received %v", "ab OR ba", list)
	}

	// sender subscribes to the multiple tags and then lists it's tags verifying a
	// "compound" subscription; we test multiple configurations here because maps
	// are unordered
	sender.Subscribe([]string{"a", "b"})
	list = flattenSliceToString(sender.List())
	switch list {
	case "aba,b", "baa,b", "a,bab", "a,bba", "aa,bb", "ba,ba", "abb,a":
		// pass
	default:
		t.Errorf("Unexpected tags - Expecting '%v' received %v", "aba,b OR baa,b OR a,bab OR a,bba", list)
	}

	// sender subscribes to the same multiple tags (unordered) and then lists it's
	// tags verifying no additional tags; we test multiple configurations here because
	// maps are unordered
	sender.Subscribe([]string{"b", "a"})
	list = flattenSliceToString(sender.List())
	switch list {
	case "aba,b", "baa,b", "a,bab", "a,bba", "aa,bb", "ba,ba", "abb,a":
		// pass
	default:
		t.Errorf("Unexpected tags - Expecting '%v' received %v", "aba,b OR baa,b OR a,bab OR a,bba", list)
	}
}

// TestTags tests to ensure that mist will send messages to single or multiple
// tags. It tests that multiple tags are received only as subscribed, and that
// multiple tags aren't receieved once unsubscribed.
func TestTags(t *testing.T) {

	//
	sender := NewProxy()
	defer sender.Close()

	//
	receiver := NewProxy()
	defer receiver.Close()

	// receiver subscribes to a single tag and then sender publishes to that tag,
	// verifying that a message is receieved...
	receiver.Subscribe([]string{"a"})
	sender.Publish([]string{"a"}, testMsg)
	verifyMessage(testMsg, receiver, t)

	// receiver subscribes to multiple tags and then sender publishes to those tags
	// verifying that a message is receieved...
	receiver.Subscribe([]string{"a", "b"})
	sender.Publish([]string{"a", "b"}, testMsg)
	verifyMessage(testMsg, receiver, t)

	// sender then published to those tags again (unordered), verifying that a
	// message is received...
	sender.Publish([]string{"b", "a"}, testMsg)
	verifyMessage(testMsg, receiver, t)

	// receiver unsubscribes from the single tag and sender publishes again, verifying
	// there is no message.
	receiver.Unsubscribe([]string{"a"})
	sender.Publish([]string{"a"}, testMsg)
	verifyNoMessage(receiver, t)

	// receiver unsubscribes from multiple tags (unordered) and sender publishes
	// again, verifying there is no message.
	receiver.Unsubscribe([]string{"b", "a"})
	sender.Publish([]string{"a", "b"}, testMsg)
	verifyNoMessage(receiver, t)
}
