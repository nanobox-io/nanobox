package server

import (
	"fmt"
	"strings"

	"github.com/nanopack/mist/core"
)

// GenerateHandlers ...
func GenerateHandlers() map[string]mist.HandleFunc {
	return map[string]mist.HandleFunc{
		"auth":        handleAuth,
		"ping":        handlePing,
		"subscribe":   handleSubscribe,
		"unsubscribe": handleUnsubscribe,
		"publish":     handlePublish,
		// "publishAfter":     handlePublishAfter,
		"list":    handleList,
		"listall": handleListAll, // listall related
		"who":     handleWho,     // who related
	}
}

// handleAuth only exists to avoid getting the message "Unknown command" when
// authing with a authenticated server
func handleAuth(proxy *mist.Proxy, msg mist.Message) error {
	return nil
}

// handlePing
func handlePing(proxy *mist.Proxy, msg mist.Message) error {
	// goroutining any of these would allow a client to spam and overwhelm the server. clients don't need the ability to ping indefinitely
	proxy.Pipe <- mist.Message{Command: "ping", Tags: []string{}, Data: "pong"}
	return nil
}

// handleSubscribe
func handleSubscribe(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Subscribe(msg.Tags)
	return nil
}

// handleUnsubscribe
func handleUnsubscribe(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Unsubscribe(msg.Tags)
	return nil
}

// handlePublish
func handlePublish(proxy *mist.Proxy, msg mist.Message) error {
	proxy.Publish(msg.Tags, msg.Data)
	return nil
}

// handlePublishAfter - how do we get the [delay] here?
// func handlePublishAfter(proxy *mist.Proxy, msg mist.Message) error {
// 	proxy.PublishAfter(msg.Tags, msg.Data, ???)
// 	go func() {
// 		proxy.Pipe <- mist.Message{Command: "publish after", Tags: msg.Tags, Data: "success"}
// 	}()
// 	return nil
// }

// handleList
func handleList(proxy *mist.Proxy, msg mist.Message) error {
	var subscriptions string
	for _, v := range proxy.List() {
		subscriptions += strings.Join(v, ",")
	}
	proxy.Pipe <- mist.Message{Command: "list", Tags: msg.Tags, Data: subscriptions}
	return nil
}

// handleListAll - listall related
func handleListAll(proxy *mist.Proxy, msg mist.Message) error {
	subscriptions := mist.Subscribers()
	proxy.Pipe <- mist.Message{Command: "listall", Tags: msg.Tags, Data: subscriptions}
	return nil
}

// handleWho - who related
func handleWho(proxy *mist.Proxy, msg mist.Message) error {
	who, max := mist.Who()
	subscribers := fmt.Sprintf("Lifetime  connections: %v\nSubscribers connected: %v", max, who)
	proxy.Pipe <- mist.Message{Command: "who", Tags: msg.Tags, Data: subscribers}
	return nil
}
