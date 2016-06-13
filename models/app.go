package models

// App ...
type App struct {
	ID    string `json:"id"`    //
	Name  string `json:"name"`  //
	State string `json:"state"` // State is used to determine the current state of the app
	DevIP string // The Dev IP is an external IP to connect to services inside the dev container
}
