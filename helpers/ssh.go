package helpers

import (
	"code.google.com/p/go.crypto/ssh"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"

	"github.com/nanobox-core/cli/ui"
)

// SSHOptions represents all the options needed when attempting an SSH action.
// These actions include 'run', 'ssh', and 'tunnel'
type SSHOptions struct {
	Command     string            // The command to run when using the 'run' action
	Config      *ssh.ClientConfig // Determines the user and auth method when connecting
	LocalIP     string            // The forward IP when tunneling (localhost)
	LocalPort   int               // The forward port when tunneling
	RemoteIP    string            //
	RemotePort  int               //
	RemoteUser  string            //
	ServerIP    string            //
	ServerPort  int               //
	ServiceApp  string            //
	ServiceUser string            //
	ServicePass string            //
}

// GetKeyFile attempts to read a ~/.ssh/id_rsa file. If none is found it prompts
// for a path to the file, or provides information on how to generate an ssh key.
// If a key is found it attempts to decode the key, prompting for a password if
// the key is encrypted, and returns the key for use with SSH.
func GetKeyFile(path string) (key ssh.Signer, err error) {

	if path == "" {
		homeDir, err := homedir.Dir()
		if err != nil {
			fmt.Println("Unable to access your home directory...\n")
			return nil, err
		}

		path = filepath.Clean(homeDir + "/.ssh/id_rsa")
	}

	//
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf(`
We were unable to find an SSH key at ` + path + `. If it's located somewhere else
you can run:
  pagoda ssh -a app-name -s service -i path/to/file

If you don't have an SSH key you can create one by running:
  ssh-keygen -t rsa -C "your_email@example.com"
    `)
		os.Exit(1)
	}

	//
	b, _ := pem.Decode(buf)

	//
	if x509.IsEncryptedPEMBlock(b) {
		password := ui.PPrompt("It appears your key is password protected. Provide your password to continue: ")

		buf, err = x509.DecryptPEMBlock(b, []byte(password))
		if err != nil {
			fmt.Println("Unable to decrypt key, please ensure correct password.\n")
			return nil, err
		}

		k, err := x509.ParsePKCS1PrivateKey(buf)
		if err != nil {
			fmt.Println("Unable to parse key...\n")
			return nil, err
		}

		key, err := ssh.NewSignerFromKey(k)
		if err != nil {
			fmt.Println("Unable to sign key...\n")
			return nil, err
		}

		return key, nil

		//
	} else {
		key, err := ssh.ParsePrivateKey(buf)
		if err != nil {
			fmt.Println("Unable to parse key...\n")
			return nil, err
		}

		return key, nil
	}
}
