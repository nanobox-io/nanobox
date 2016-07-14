package models

// App ...
type App struct {
	ID        string
	Directory string
	Name      string
	// State is used to ensure we don't setup the app multiple times
	State  string
	Status string
	// There are certain global ips that need to be reserved across container
	// lifetimes. The dev ip and preview ip are examples. We'll store those here.
	GlobalIPs map[string]string
	// There are also certain platform service ips that need to 1) remain constant
	// even if the component were repaired and 2) be available even before the
	// component is. logvac and mist ips are examples. We'll store those here.
	LocalIPs map[string]string
}
