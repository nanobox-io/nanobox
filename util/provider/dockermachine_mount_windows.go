package provider

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// add mounts using cifs for windows development
func (machine DockerMachine) addNetfsMount(local, host string) error {
	appID := config.EnvID()
	user := os.Getenv("USERNAME")

	// pause the current task
	display.PauseTask()
	// wait a bit to ensure the output doesn't get messed up
	<-time.After(time.Second * 1)

	// fetch the password from the user
	fmt.Printf("%s's password is required to mount a Windows share. (must be your Windows Live password if linked)\n", user)
	pass, err := display.ReadPassword("Windows")
	if err != nil {
		return err
	}

	if len(pass) == 0 {
		return fmt.Errorf("currently we do not support passwordless windows accounts")
	}

	// resume the task
	display.ResumeTask()

	// ensure the destination directory exists
	cmd := []string{"sudo", "/bin/mkdir", "-p", host}
	if b, err := Run(cmd); err != nil {
		lumber.Debug("mkdir output: %s", b)
		return fmt.Errorf("mkdir:%s, %s", b, err.Error())
	}

	// ensure cifs/samba utilities are installed
	cmd = []string{"sh", "-c", setupCifsUtilsScript()}
	if b, err := Run(cmd); err != nil {
		lumber.Debug("cifs output: %s", b)
		return fmt.Errorf("cifs:%s", err.Error())
	}

	// mount!
	// mount -t cifs -o sec=ntlmssp,username=USER,password=PASSWORD,uid=1000,gid=1000 //192.168.99.1/<path to app> /<vm location>
	source := fmt.Sprintf("//192.168.99.1/nanobox-%s", appID)
	// mfsymlinks,
	config, _ := models.LoadConfig()
	additionalOptions := config.NetfsMountOpts
	opts := fmt.Sprintf("nodev,sec=ntlmssp,user='%s',password='%s',uid=1000,gid=1000", user, pass)
	if additionalOptions != "" {
		opts = fmt.Sprintf("%s,%s", additionalOptions, opts)
	}
	cmd = []string{
		"sudo",
		"/bin/mount",
		"-t",
		"cifs",
		"-o",
		opts,
		source,
		host,
	}
	lumber.Debug("cifs mount cmd: %v", cmd)
	if b, err := Run(cmd); err != nil {
		lumber.Debug("mount output: %s", b)
		return fmt.Errorf("mount: output: %s err:%s", b, err.Error())
	}

	return nil
}

// setupCifsUtilsScript returns a string containing the script to setup cifs
func setupCifsUtilsScript() string {
	script := `
		if [ ! -f /sbin/mount.cifs ]; then
			wget -O /mnt/sda1/tmp/tce/optional/samba-libs.tcz http://repo.tinycorelinux.net/7.x/x86_64/tcz/samba-libs.tcz &&
			wget -O /mnt/sda1/tmp/tce/optional/cifs-utils.tcz http://repo.tinycorelinux.net/7.x/x86_64/tcz/cifs-utils.tcz &&

			tce-load -i samba-libs &&
			tce-load -i cifs-utils;
		fi
	`

	return strings.Replace(script, "\n", "", -1)
}
