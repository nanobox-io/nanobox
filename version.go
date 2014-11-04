package main

// Version returns the current version of the Nanobox CLi
func (cli *CLI) Version() string {
	return cli.version
}

// CheckVersion checks to see if the current 'installed' version of the Nanobox
// CLI matches the most current release version.
func (cli *CLI) CheckVersion() {}

// Update updates the local Nanobox CLI to the most recent release
func (cli *CLI) Update() {}
