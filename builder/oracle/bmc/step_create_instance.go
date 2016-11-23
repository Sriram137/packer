package oracle

import (
	"fmt"
	"github.com/elricL/oracle_bmc_sdk"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
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

	input := &oraclebmc_sdk.LaunchInstanceInput{
		CompartmentId:      s.compartmentId,
		AvailabilityDomain: s.availablityZone,
		DisplayName:        "packerCreateMachine",
		ImageId:            s.baseImageId,
		Metadata:           map[string]string{},
		Shape:              s.shape,
		SubnetId:           s.subnetId}
	ui.Say("Starting instance creation")
	instance, err := computeApi.CreateInstance(input)
	if err != nil {
		ui.Say(fmt.Sprintf("Encountered error: %s", err))
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

func (s *StepCreateInstance) Cleanup(multistep.StateBag) {}
