package main

import (
	pagodaAPI "github.com/nanobox-core/api-client-go"
	"github.com/nanobox-core/cli/commands"
)

// Commands represents a map of all the available commands that the Pagoda Box
// CLI can run
var Commands map[string]Command

// Command represents a Pagoda Box CLI command. Every command must have a Help()
// and Run() function
type Command interface {
	Help()                                                 // Prints the help text associated with this command
	Run(fApp string, opts []string, api *pagodaAPI.Client) // Houses the logic that will be run upon calling this command
}

// init builds the list of available Pagoda Box CLI commands
func init() {

	// the map of all available commands the Pagoda Box CLI can run
	Commands = map[string]Command{
		"create":          &commands.AppCreateCommand{},
		"destroy":         &commands.AppDestroyCommand{},
		"env":             &commands.EVarListCommand{},
		"info":            &commands.AppInfoCommand{},
		"list":            &commands.AppListCommand{},
		"log":             &commands.AppLogCommand{},
		"open":            &commands.AppOpenCommand{},
		"rebuild":         &commands.AppRebuildCommand{},
		"rollback":        &commands.AppRollbackCommand{},
		"run":             &commands.ServiceRunCommand{},
		"ssh":             &commands.ServiceSSHCommand{},
		"tunnel":          &commands.ServiceTunnelCommand{},
		"app:create":      &commands.AppCreateCommand{},
		"app:destroy":     &commands.AppDestroyCommand{},
		"app:info":        &commands.AppInfoCommand{},
		"app:list":        &commands.AppListCommand{},
		"app:log":         &commands.AppLogCommand{},
		"app:open":        &commands.AppOpenCommand{},
		"app:rebuild":     &commands.AppRebuildCommand{},
		"app:rollback":    &commands.AppRollbackCommand{},
		"evar:create":     &commands.EVarCreateCommand{},
		"evar:destroy":    &commands.EVarDestroyCommand{},
		"evar:list":       &commands.EVarListCommand{},
		"service:info":    &commands.ServiceInfoCommand{},
		"service:list":    &commands.ServiceListCommand{},
		"service:reboot":  &commands.ServiceRebootCommand{},
		"service:repair":  &commands.ServiceRepairCommand{},
		"service:restart": &commands.ServiceRestartCommand{},
		"service:run":     &commands.ServiceRunCommand{},
		"service:ssh":     &commands.ServiceSSHCommand{},
		"service:tunnel":  &commands.ServiceTunnelCommand{},
	}

}
