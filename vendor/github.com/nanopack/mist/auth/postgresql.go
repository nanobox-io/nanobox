package auth

import (
	"fmt"
	"net"
	"net/url"

	"database/sql"

	// this is how lib postgres is intended to work (blank import)
	_ "github.com/lib/pq"
)

type (
	postgres string
)

// add "postgres" to the list of supported Authenticators
func init() {
	Register("postgres", NewPostgres)
}

// NewPostgres creates a new "postgres" Authenticator
func NewPostgres(url *url.URL) (Authenticator, error) {

	//
	host, port, err := net.SplitHostPort(url.Host)
	db := url.Query().Get("db")
	user := url.User.Username()
	// pass, _ := url.User.Password()

	//
	pg := postgres(fmt.Sprintf("user=%v database=%v sslmode=disable host=%v port=%v", user, db, host, port))

	// create the tables needed to support mist authentication
	if _, err = pg.exec(`
CREATE TABLE IF NOT EXISTS tokens (
	token text NOT NULL,
	token_id SERIAL UNIQUE NOT NULL,
	PRIMARY KEY (token)
)`); err != nil {
		return pg, err
	}

	if _, err = pg.exec(`
CREATE TABLE IF NOT EXISTS tags (
  token_id integer NOT NULL REFERENCES tokens (token_id) ON DELETE CASCADE,
  tag text NOT NULL,
  PRIMARY KEY (token_id, tag)
)`); err != nil {
		return pg, err
	}

	return pg, nil
}

// AddToken
func (a postgres) AddToken(token string) error {
	_, err := a.exec("INSERT INTO tokens (token) VALUES ($1)", token)
	return err
}

// RemoveToken
func (a postgres) RemoveToken(token string) error {
	_, err := a.exec("DELETE FROM tokens WHERE token = $1", token)
	return err
}

// AddTags
func (a postgres) AddTags(token string, tags []string) error {

	// This could be optimized a LOT
	for _, tag := range tags {
		// errors are ignored, this may not be the best idea.
		a.exec("INSERT INTO tags (token_id,tag) VALUES ((SELECT token_id FROM tokens WHERE token = $1), $2)", token, tag)
	}
	return nil
}

// RemoveTags
func (a postgres) RemoveTags(token string, tags []string) error {
	for _, tag := range tags {
		a.exec("DELETE FROM tags INNER JOIN tokens ON (tags.token_id = tokens.token_id) WHERE token = $1 AND tag = $2", token, tag)
	}
	return nil
}

// GetTagsForToken
func (a postgres) GetTagsForToken(token string) ([]string, error) {

	//
	rows, err := a.query("SELECT tag FROM tags INNER JOIN tokens ON (tags.token_id = tokens.token_id) WHERE token = $1", token)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//
	var tags []string

	//
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	err = rows.Err()

	switch {
	case len(tags) == 0:
		return tags, ErrTokenNotFound
	default:
		return tags, err
	}
}

//
func (a postgres) connect() (*sql.DB, error) {
	return sql.Open("postgres", string(a))
}

// this could really be optimized a lot. instead of opening a new conenction for
// each query, it should reuse connections
func (a postgres) query(query string, args ...interface{}) (*sql.Rows, error) {
	client, err := a.connect()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.Query(query, args...)
}

// This could also be optimized a lot
func (a postgres) exec(query string, args ...interface{}) (sql.Result, error) {
	client, err := a.connect()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.Exec(query, args...)
}
