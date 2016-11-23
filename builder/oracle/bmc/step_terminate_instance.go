package oracle

import (
	"fmt"
	"github.com/elricL/oracle_bmc_sdk"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type StepTerminateInstance struct {
}

func (s *StepTerminateInstance) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	computeApi := state.Get("api").(oraclebmc_sdk.ComputeApi)
	instance := state.Get("image").(*oraclebmc_sdk.Instance)

	ui.Say("Starting instance creation")
	err := computeApi.TerminateInstance(instance)

	if err != nil {
		ui.Say(fmt.Sprintf("Encountered error while terminating instnace: %s", err))
		return multistep.ActionHalt
	}
	ui.Say(fmt.Sprintf("Waiting %s insatnce to be terminated", instance.Id))
	err = computeApi.WaitForInstance(instance, "TERMINATED")
	if err != nil {
		ui.Say(fmt.Sprintf("Encountered error while waiting for image creation:%s", err))
	}

	ui.Say("Instance Termination is done")
	return multistep.ActionContinue
}

func (s *StepTerminateInstance) Cleanup(multistep.StateBag) {}
