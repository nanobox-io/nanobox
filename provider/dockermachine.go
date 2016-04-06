package provider

import (

)

type (
	DockerMachine struct{

	}
)

func init() {
	Register("docker_machine", DockerMachine{})
}


func (self DockerMachine) Create() error {
	return nil
}


func (self DockerMachine) Reboot() error {
	return nil
}


func (self DockerMachine) Stop() error {
	return nil
}


func (self DockerMachine) Destroy() error {
	return nil
}


func (self DockerMachine) Start() error {
	return nil
}


func (self DockerMachine) AddIP(ip string) error {
	return nil
}


func (self DockerMachine) RemoveIP(ip string) error {
	return nil
}


func (self DockerMachine) AddNat(ip, ip string) error {
	return nil
}


func (self DockerMachine) RemoveNat(ip, ip string) error {
	return nil
}


func (self DockerMachine) AddMount(local, host string) error {
	return nil
}


func (self DockerMachine) RemoveMount(local, host string) error {
	return nil
}

