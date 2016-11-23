package oracle

import (
	"github.com/elricL/oracle_bmc_sdk"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
	"log"
)

const BuilderId = "sriramg.oraclebmc"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	ImageConfig         `mapstructure:",squash"`
	OracleConfig        `mapstructure:",squash"`
	InstanceConfig      `mapstructure:",squash"`
	ctx                 interpolate.Context
}

type ImageConfig struct {
	imageDisplayName string `mapstructure:"image_name"`
}

type OracleConfig struct {
	compartmentId      string `mapstructure:"compartmentId"`
	userId             string `mapstructure:"userId"`
	fingerprint        string `mapstructure:"fingerprint"`
	privateKeyContents string `mapstructure:"privateKeyContents"`
	tenantId           string `mapstructure:"tenantId"`
}

type InstanceConfig struct {
	baseImageId        string `mapstructure:"baseImageId"`
	availabilityDomain string `mapstructure:"availabilityDomain"`
	subnetId           string `mapstructure:"subnetId"`
	shape              string `mapstructure:"shape"`
}

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	return nil, nil
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	config := b.config

	oracleConfig := oraclebmc_sdk.NewConfig(config.userId, config.tenantId, config.fingerprint, config.privateKeyContents)

	computeApi := oraclebmc_sdk.ComputeApi{Config: oracleConfig}

	// Setup the state bag and initial state for the steps
	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("api", computeApi)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		&StepCreateInstance{
			availablityZone: config.availabilityDomain,
			baseImageId:     config.baseImageId,
			compartmentId:   config.availabilityDomain,
			shape:           config.shape,
			subnetId:        config.subnetId,
		},
		&StepCreateImage{
			DisplayName:   config.imageDisplayName,
			CompartmentId: config.compartmentId,
		},
		&StepTerminateInstance{},
	}

	b.runner = common.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(state)
	image := state.Get("image").(*oraclebmc_sdk.Image)

	artifact := &ImageArtifact{
		builderId:   BuilderId,
		displayName: image.DisplayName,
		imageId:     image.Id,
	}

	return artifact, nil
}
