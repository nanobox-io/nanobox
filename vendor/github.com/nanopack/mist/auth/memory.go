package auth

import (
	"net/url"

	"github.com/deckarep/golang-set"
)

// memory is a new in-memory set (map) of token/tag combination
type memory map[string]mapset.Set

// add "memory" to the list of supported Authenticators
func init() {
	Register("memory", NewMemory)
}

// NewMemory creates a new in-memory Authenticator
func NewMemory(url *url.URL) (Authenticator, error) {
	return memory{}, nil
}

// Add Token
func (a memory) AddToken(token string) error {

	// look for an existing token
	if _, err := a.findMemoryToken(token); err == nil {
		return ErrTokenExist
	}

	// create a new token
	a[token] = mapset.NewSet()

	//
	return nil
}

// RemoveToken
func (a memory) RemoveToken(token string) error {
	delete(a, token)
	return nil
}

// AddTags
func (a memory) AddTags(token string, tags []string) error {

	// look for an existing token
	entry, err := a.findMemoryToken(token)
	if err != nil {
		return err
	}

	// add new tags individually to avoid duplication
	for _, tag := range tags {
		entry.Add(tag)
	}

	//
	return nil
}

// RemoveTags
func (a memory) RemoveTags(token string, tags []string) error {

	// look for an existing token
	entry, err := a.findMemoryToken(token)
	if err != nil {
		return err
	}

	// remove tags
	for _, tag := range tags {
		entry.Remove(tag)
	}

	//
	return nil
}

// GetTagsForToken
func (a memory) GetTagsForToken(token string) ([]string, error) {

	// look for an existing token
	entry, err := a.findMemoryToken(token)
	if err != nil {
		return nil, err
	}

	// convert tags from map to slice
	var tags []string
	for _, tag := range entry.ToSlice() {
		tags = append(tags, tag.(string))
	}

	//
	return tags, nil
}

// findMemoryToken attempts to find the desired token within memory
func (a memory) findMemoryToken(token string) (mapset.Set, error) {

	// look for existing token
	entry, ok := a[token]
	if !ok {
		return nil, ErrTokenNotFound
	}

	//
	return entry, nil
}
