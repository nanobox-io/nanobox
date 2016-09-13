package auth

import (
	"net/url"
)

// redis is an Authenticator that interfaces with the redis database
type redis struct{}

// add "redis" to the list of supported Authenticatores
func init() {
	Register("redis", NewRedis)
}

// NewRedis creates a new "redis" Authenticator
func NewRedis(url *url.URL) (Authenticator, error) {
	return &redis{}, nil
}

// AddToken
func (a *redis) AddToken(token string) error {
	return nil
}

// RemoveToken
func (a *redis) RemoveToken(token string) error {
	return nil
}

// AddTags
func (a *redis) AddTags(token string, tags []string) error {
	return nil
}

// RemoveTags
func (a *redis) RemoveTags(token string, tags []string) error {
	return nil
}

// GetTagsForToken
func (a *redis) GetTagsForToken(token string) ([]string, error) {
	return nil, nil
}
