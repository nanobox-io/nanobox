package mist

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

var (
	testMsg     = "test"
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// BenchmarkPublish
func BenchmarkPublish(b *testing.B) {

	//
	p := NewProxy()
	defer p.Close()

	//
	p.Subscribe([]string{"a"})

	b.ResetTimer()

	//
	for i := 0; i < b.N; i++ {
		p.Publish([]string{"a"}, testMsg)
	}
}

// TestPublish tests that the publish Publish method publishes to all subscribers
func TestPublish(t *testing.T) {

	//
	p1 := NewProxy()
	defer p1.Close()

	p2 := NewProxy()
	defer p2.Close()

	//
	p1.Subscribe([]string{"a"})
	p2.Subscribe([]string{"a"})

	// have mist publish the message
	Publish([]string{"a"}, testMsg)

	//
	verifyMessage(testMsg, p1, t)
	verifyMessage(testMsg, p2, t)

	p1.Unsubscribe([]string{"a"})
	p2.Unsubscribe([]string{"a"})

	// have mist publish the message
	Publish([]string{"a"}, testMsg)

	// proxies should NOT get a message this time
	verifyNoMessage(p1, t)
	verifyNoMessage(p2, t)
}

// verifyMessage waits for a message to come to a proxy then tests to see if it's
// the expected message. After 1 second it assumes no message is coming and fails.
func verifyMessage(expected string, p *Proxy, t *testing.T) {
	select {
	case msg := <-p.Pipe:
		if msg.Data != expected {
			t.Fatalf("Incorrect data: Expected '%v' received '%v'\n", msg, msg.Data)
		}
		break
	case <-time.After(time.Second * 1):
		t.Errorf("Expecting messages, received none!")
	}
}

// verifyNoMessage waits for a message that should never come, assuming after 1
// second that no message is coming.
func verifyNoMessage(p *Proxy, t *testing.T) {
	select {
	case <-p.Pipe:
		t.Fatalf("Unexpected message!")
	case <-time.After(time.Second * 1):
		break
	}
}

// randKey
func randKey() string {
	b := make([]byte, 3)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// flattenSliceToString
func flattenSliceToString(list [][]string) (flat string) {
	for _, v := range list {
		flat += strings.Join(v, ",")
	}
	return
}
