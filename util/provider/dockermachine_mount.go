package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/processors/env/share"
	"github.com/nanobox-io/nanobox/util/config"
)

// HasMount checks to see if the mount exists in the vm
func (machine DockerMachine) HasMount(mount string) bool {

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"nanobox",
		"sudo",
		"cat",
		"/proc/mounts",
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), mount)
}

// AddMount adds a virtualbox mount into the docker-machine vm
func (machine DockerMachine) AddMount(local, host string) error {

	// stop early if already mounted
	if machine.HasMount(host) {
		lumber.Info("mount exists for %s", host)
		return nil
	}

	switch config.Viper().GetString("mount-type") {

	case "netfs":

		// add netfs share
		// here we use the processor so we can do privilage exec
		lumber.Info("adding share in netfs for %s", local)
		if err := share.Add(local); err != nil {
			return err
		}
		// add netfs mount

		lumber.Info("adding mount in for netfs %s -> %s", local, host)
		if err := machine.addNetfsMount(local, host); err != nil {
			return err
		}

	default:

		// add share
		lumber.Info("adding share for native %s", local)
		if err := machine.addShare(local, host); err != nil {
			return err
		}

		// add mount
		lumber.Info("adding mount for  native %s -> %s", local, host)
		if err := machine.addNativeMount(local, host); err != nil {
			return err
		}
	}

	return nil
}

// RemoveMount removes a mount from the docker-machine vm
func (machine DockerMachine) RemoveMount(local, host string) error {
	if !machine.HasMount(host) {
		return nil
	}

	// all mounts are removed as if they are native
	if err := machine.removeNativeMount(local, host); err != nil {
		return err
	}

	// if we are supposed to keep the shares return here
	if registry.GetBool("keep-share") {
		return nil
	}

	// remove any netfs shares
	if err := share.Remove(local); err != nil {
		return err
	}
	// remove any native shares
	if err := machine.removeShare(local, host); err != nil {
		return err
	}

	return nil
}

// hasShare checks to see if the share exists
func (machine DockerMachine) hasShare(local, host string) bool {
	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	cmd := []string{
		vboxManageCmd,
		"showvminfo",
		"nanobox",
		"--machinereadable",
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	output, err := process.CombinedOutput()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), name)
}

// AddShare adds the provided path as a shareable filesystem
func (machine DockerMachine) addShare(local, host string) error {

	if machine.hasShare(local, host) {
		return nil
	}

	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	cmd := []string{
		vboxManageCmd,
		"sharedfolder",
		"add",
		"nanobox",
		"--name",
		name,
		"--hostpath",
		local,
		"--transient",
	}

	lumber.Info("add share native cmd: %v", cmd)
	process := exec.Command(cmd[0], cmd[1:]...)
	b, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// todo: check output for failures
	return nil
}

// RemoveShare removes the provided path as a shareable filesystem; we don't care
// what the user has configured, we need to remove any shares that may have been
// setup previously
func (machine DockerMachine) removeShare(local, host string) error {

	if !machine.hasShare(local, host) {
		return nil
	}

	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	cmd := []string{
		vboxManageCmd,
		"sharedfolder",
		"remove",
		"nanobox",
		"--name",
		name,
		"--transient",
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	b, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// todo: check output for failures

	return nil
}

func (machine DockerMachine) addNativeMount(local, host string) error {
	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	// create folder

	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"nanobox",
		"sudo",
		"mkdir",
		"-p",
		host,
	}

	lumber.Info("add mount native cmd (mkdir): %v", cmd)
	process := exec.Command(cmd[0], cmd[1:]...)
	b, err := process.CombinedOutput()
	lumber.Info("mkdir data: %s", b)
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// mount
	cmd = []string{
		dockerMachineCmd,
		"ssh",
		"nanobox",
		"sudo",
		"mount",
		"-t",
		"vboxsf",
		"-o",
		"uid=1000,gid=1000",
		name,
		host,
	}

	lumber.Info("add mount native cmd: %v", cmd)
	process = exec.Command(cmd[0], cmd[1:]...)
	b, err = process.CombinedOutput()
	lumber.Info("mount data: %s", b)
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// todo: check output for failures

	return nil
}

func (machine DockerMachine) removeNativeMount(local, host string) error {
	cmd := []string{
		dockerMachineCmd,
		"ssh",
		"nanobox",
		"sudo",
		"umount",
		host,
	}

	process := exec.Command(cmd[0], cmd[1:]...)
	b, err := process.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", b, err)
	}

	// todo: check output for failures

	return nil
}
