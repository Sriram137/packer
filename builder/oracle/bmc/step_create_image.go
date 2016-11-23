package oracle

import (
	"fmt"
	"github.com/elricL/oracle_bmc_sdk"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type StepCreateImage struct {
	CompartmentId string
	DisplayName   string
}

func (s *StepCreateImage) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	computeApi := state.Get("api").(oraclebmc_sdk.ComputeApi)
	instance := state.Get("instance").(*oraclebmc_sdk.Instance)

	createImageInput := &oraclebmc_sdk.CreateImageInput{
		CompartmentId: s.CompartmentId,
		DisplayName:   s.DisplayName,
		InstanceId:    instance.Id,
	}

	ui.Say("Starting image creation")
	image, err := computeApi.CreateImage(createImageInput)

	if err != nil {
		ui.Say(fmt.Sprintf("Encountered error while capturing image: %s", err))
		return multistep.ActionHalt
	}
	ui.Say(fmt.Sprintf("Waiting %s image to be ready", instance.Id))
	err = computeApi.WaitForImage(image, "AVAILABLE")
	if err != nil {
		ui.Say(fmt.Sprintf("Encountered error while waiting for image creation:%s", err))
	}

	state.Put("image", image)
	ui.Say(fmt.Sprintf("Image is ready %s", image.Id))
	return multistep.ActionContinue
}

func (s *StepCreateImage) Cleanup(multistep.StateBag) {}
