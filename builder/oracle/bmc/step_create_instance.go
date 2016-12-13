package oracle

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/elricL/oracle_bmc_sdk"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
	"golang.org/x/crypto/ssh"
)

type StepCreateInstance struct {
	availablityZone string
	compartmentId   string
	baseImageId     string
	shape           string
	subnetId        string
}

func (s *StepCreateInstance) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	computeApi := state.Get("api").(oraclebmc_sdk.ComputeApi)
	config := state.Get("config").(*Config)

	priv, publ := createPrivatePublicPair()

	state.Put("privateKey", priv)

	if config.OraclePublicKeys != nil {
		for _, element := range config.OraclePublicKeys {
			publ = publ + "\n" + element
		}
	}

	input := &oraclebmc_sdk.LaunchInstanceInput{
		CompartmentId:      s.compartmentId,
		AvailabilityDomain: s.availablityZone,
		DisplayName:        "packerCreateMachine",
		ImageId:            s.baseImageId,
		Metadata:           map[string]string{"ssh_authorized_keys": publ},
		Shape:              s.shape,
		SubnetId:           s.subnetId}
	ui.Say("Starting instance creation")
	ui.Say(fmt.Sprintf("%#v", input))
	ui.Say(fmt.Sprintf("%#v", computeApi.Config))
	instance, err := computeApi.CreateInstance(input)
	if err != nil {
		ui.Say(fmt.Sprintf("Encountered error: %s", err))
		return multistep.ActionHalt
	}
	if instance.Id == "" {
		ui.Say("Empty instance Id")
		return multistep.ActionHalt
	}
	ui.Say(instance.Id)
	ui.Say(fmt.Sprintf("Waiting %s for instance to start", instance.Id))

	err = computeApi.WaitForInstance(instance, "RUNNING")
	if err != nil {
		ui.Say(fmt.Sprintf("Encountered error while waiting for instance:%s", err))
	}
	state.Put("instance", instance)
	return multistep.ActionContinue
}

func createPrivatePublicPair() (string, string) {
	priv, _ := rsa.GenerateKey(rand.Reader, 2014)

	// ASN.1 DER encoded form
	priv_der := x509.MarshalPKCS1PrivateKey(priv)
	priv_blk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   priv_der,
	}
	// Set the private key in the statebag for later
	privateKey := string(pem.EncodeToMemory(&priv_blk))

	// Marshal the public key into SSH compatible format
	// TODO properly handle the public key error
	pub, _ := ssh.NewPublicKey(&priv.PublicKey)
	pub_sshformat := string(ssh.MarshalAuthorizedKey(pub))
	return privateKey, pub_sshformat
}

func (s *StepCreateInstance) Cleanup(multistep.StateBag) {}
