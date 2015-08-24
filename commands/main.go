// Copyright (c) 2015 Pagoda Box Inc
//
// This Source Code Form is subject to the terms of the Mozilla Public License, v.
// 2.0. If a copy of the MPL was not distributed with this file, You can obtain one
// at http://mozilla.org/MPL/2.0/.
//

package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pagodabox/nanobox-cli/config"
	"github.com/pagodabox/nanobox-golang-stylish"
)

//
var (
	NanoboxCmd = &cobra.Command{
		Use:   "nanobox",
		Short: "",
		Long:  ``,

		//
		// Run: func(cmd *cobra.Command, args []string) {},
	}

	productionCmd = &cobra.Command{
		Use:   "nanobox production",
		Short: "",
		Long:  ``,

		//
		// Run: func(cmd *cobra.Command, args []string) {},
	}

	// flags
	fCount   int    //
	fDebug   bool   //
	fForce   bool   //
	fLevel   string //
	fRemove  bool   //
	fReset   bool   //
	fSandbox bool   //
	fStream  bool   //
	fTunnel  string //
	fVerbose bool   //
	fWatch   bool   //
	fWrite   bool   //
)

// init builds the list of available nanobox commands and sub commands
func init() {

	NanoboxCmd.PersistentFlags().BoolVarP(&fDebug, "debug", "d", false, "display debug output")
	// NanoboxCmd.SetUsageFunc(nanoHelp)

	// all available nanobox commands
	NanoboxCmd.AddCommand(bootstrapCmd)
	NanoboxCmd.AddCommand(buildCmd)
	NanoboxCmd.AddCommand(consoleCmd)
	NanoboxCmd.AddCommand(createCmd)
	NanoboxCmd.AddCommand(deployCmd)
	NanoboxCmd.AddCommand(destroyCmd)
	NanoboxCmd.AddCommand(execCmd)
	NanoboxCmd.AddCommand(fetchCmd)
	NanoboxCmd.AddCommand(haltCmd)
	NanoboxCmd.AddCommand(initCmd)
	NanoboxCmd.AddCommand(logCmd)
	NanoboxCmd.AddCommand(newCmd)
	NanoboxCmd.AddCommand(publishCmd)
	NanoboxCmd.AddCommand(reloadCmd)
	NanoboxCmd.AddCommand(resumeCmd)
	NanoboxCmd.AddCommand(sshCmd)
	NanoboxCmd.AddCommand(statusCmd)
	NanoboxCmd.AddCommand(suspendCmd)
	NanoboxCmd.AddCommand(tunnelCmd)
	NanoboxCmd.AddCommand(upCmd)
	NanoboxCmd.AddCommand(updateCmd)
	NanoboxCmd.AddCommand(upgradeCmd)
	NanoboxCmd.AddCommand(versionCmd)
	NanoboxCmd.AddCommand(watchCmd)

	NanoboxCmd.AddCommand(productionCmd)
	// productionCmd.AddCommand(deployCmd)
}

// runVagrantCommand provides a wrapper around a standard cmd.Run() in which
// all standard in/outputs are connected to the command, and the directory is
// changed to the corresponding app directory. This allows nanobox to run Vagrant
// commands w/o contaminating a users codebase
func runVagrantCommand(cmd *exec.Cmd) error {

	// run the command from ~/.nanobox/apps/<config.App>. if the directory doesn't
	// exist, simply return
	if err := os.Chdir(config.AppDir); err != nil {
		return err
	}

	// create a pipe that we can pipe the cmd standard output's too. The reason this
	// is done rather than just piping directly to os standard outputs and .Run()ing
	// the command (vs .Start()ing) is because the output needs to be modified
	// according to http://nanodocs.gopagoda.io/engines/style-guide
	//
	// NOTE: the reason it's done this way vs using the cmd.*Pipe's is so that all
	// the command output can be read from a single pipe, rather than having to create
	// a new pipe/scanner for each type of output
	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()

	// connect standard output
	cmd.Stdout = pw
	cmd.Stderr = pw

	//
	fmt.Printf(stylish.Bullet(fmt.Sprintf("running '%v'", strings.Trim(fmt.Sprint(cmd.Args), "[]"))))

	// scan the command output modifying it according to
	// http://nanodocs.gopagoda.io/engines/style-guide
	scanner := bufio.NewScanner(pr)
	go func() {
		for scanner.Scan() {
			fmt.Printf("   %s\n", strings.Replace(scanner.Text(), "\r", "\n", -1))
		}
	}()

	// start the command; we need this to 'fire and forget' so that we can manually
	// capture and modify the commands output
	if err := cmd.Start(); err != nil {
		return err
	}

	// wait for the command to complete/fail and exit, giving us a chance to scan
	// the output
	if err := cmd.Wait(); err != nil {
		return err
	}

	// switch back to project dir
	if err := os.Chdir(config.CWDir); err != nil {
		return err
	}

	return nil
}
