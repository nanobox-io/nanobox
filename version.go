package main

// Version returns the current version of the Pagoda Box CLi
func (cli *CLI) Version() string {
	return cli.version
}

// CheckVersion checks to see if the current 'installed' version of the Pagoda Box
// CLI matches the most current release version.
func (cli *CLI) CheckVersion() {}

// Update updates the local Pagoda Box CLI to the most recent release
func (cli *CLI) Update() {}
