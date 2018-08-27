package oracle

import (
	"fmt"
	"github.com/elricL/oracle_bmc_sdk"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"golang.org/x/crypto/ssh"
)

func commHost(state multistep.StateBag) (string, error) {
	instance := state.Get("instance").(*oraclebmc_sdk.Instance)
	computeApi := state.Get("api").(oraclebmc_sdk.ComputeApi)

	ui := state.Get("ui").(packer.Ui)
	vnicAttachments, err := computeApi.ListVnicAttachments(instance.CompartmentId, instance.Id)
	if err != nil {
		ui.Say("Unable to get VnicAttachemnts")
		return "", err
	}

	vnic, err := computeApi.GetVnic((*vnicAttachments)[0].VnicId)
	if err != nil {
		ui.Say("Unable to get Vnic")
		return "", err
	}
	return vnic.PublicIp, nil
}

func sshConfig(state multistep.StateBag) (*ssh.ClientConfig, error) {
	config := state.Get("config").(*Config)
	privateKey := state.Get("privateKey").(string)

	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, fmt.Errorf("Error setting up SSH config: %s", err)
	}

	return &ssh.ClientConfig{
		User: config.Comm.SSHUsername,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}, nil
}
