package models

type App struct {
	ID   		string `json:"id"`
	Name 		string `json:"name"`
	// State is used to determine the current state of the app
	State		string `json:"state"`
	// The Dev IP is an external IP to connect
	// to services inside the dev container
	DevIP		string
}
