package provider

import(
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"os/exec"
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
		return nil
	}

	mountType := 
	switch config.Viper().GetString("mount-type") {

	case "netfs":
		// add netfs share
		// here we use the processor so we can do privilage exec
		if err := share.Add(local); err != nil  {
			return err
		}
		// add netfs mount
		if err := machine.addNetfsMount(local, host); err != nil {
			return err
		}

	default:
		
		// add share
		if err := machine.addShare(local, host); err != nil {
			return err
		}

		// add mount
		if err := machine.addNativeMount(local, host); err != nil {
			return err
		}
	}

	// switch between types
		// add netfs share
		// add netfs mount


	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	if !machine.HasMount(host) {

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

		process := exec.Command(cmd[0], cmd[1:]...)
		b, err := process.CombinedOutput()
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

		process = exec.Command(cmd[0], cmd[1:]...)
		b, err = process.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %s", b, err)
		}

		// todo: check output for failures

	}

	return nil
}

// RemoveMount removes a mount from the docker-machine vm
func (machine DockerMachine) RemoveMount(_, host string) error {

	if machine.HasMount(host) {

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

	}

	return nil
}



// HasShare checks to see if the share exists
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
	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	if !machine.HasShare(local, host) {

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

		process := exec.Command(cmd[0], cmd[1:]...)
		b, err := process.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %s", b, err)
		}

		// todo: check output for failures

	}

	return nil
}

// RemoveShare removes the provided path as a shareable filesystem; we don't care
// what the user has configured, we need to remove any shares that may have been
// setup previously
func (machine DockerMachine) removeShare(local, host string) error {
	h := sha256.New()
	h.Write([]byte(local))
	h.Write([]byte(host))
	name := hex.EncodeToString(h.Sum(nil))

	if machine.HasShare(local, host) {

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

	}

	return nil
}

