package auth

import "github.com/nanopack/mist/core"

// GenerateHandlers ...
func GenerateHandlers() map[string]mist.HandleFunc {
	return map[string]mist.HandleFunc{
		"register":   handleRegister,
		"unregister": handleUnregister,
		"set":        handleSet,
		"unset":      handleUnset,
		"tags":       handleTags,
	}
}

// handleRegister
func handleRegister(proxy *mist.Proxy, msg mist.Message) error {

	//
	if err := DefaultAuth.AddToken(msg.Data); err != nil {
		return err
	}

	//
	if err := DefaultAuth.AddTags(msg.Data, msg.Tags); err != nil {
		return err
	}

	//
	return nil
}

// handleUnregister
func handleUnregister(proxy *mist.Proxy, msg mist.Message) error {

	//
	if err := DefaultAuth.RemoveToken(msg.Data); err != nil {
		return err
	}

	//
	return nil
}

// handleSet
func handleSet(proxy *mist.Proxy, msg mist.Message) error {

	//
	if err := DefaultAuth.AddTags(msg.Data, msg.Tags); err != nil {
		return err
	}

	//
	return nil
}

// handleUnset
func handleUnset(proxy *mist.Proxy, msg mist.Message) error {

	//
	if err := DefaultAuth.RemoveTags(msg.Data, msg.Tags); err != nil {
		return err
	}

	//
	return nil
}

// handleTags
func handleTags(proxy *mist.Proxy, msg mist.Message) error {

	//
	tags, err := DefaultAuth.GetTagsForToken(msg.Data)
	if err != nil {
		return err
	}

	proxy.Pipe <- mist.Message{Command: "tags", Tags: tags}

	return nil
}
