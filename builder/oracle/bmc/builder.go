package oracle

import (
	"github.com/elricL/oracle_bmc_sdk"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/helper/config"
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
	DebugConfig         `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`
	ctx                 interpolate.Context
}

type ImageConfig struct {
	ImageDisplayName string `mapstructure:"image_name"`
}

type OracleConfig struct {
	CompartmentId      string `mapstructure:"compartmentId"`
	UserId             string `mapstructure:"userId"`
	Fingerprint        string `mapstructure:"fingerprint"`
	PrivateKeyContents string `mapstructure:"privateKeyContents"`
	TenantId           string `mapstructure:"tenantId"`
}

type InstanceConfig struct {
	BaseImageId        string `mapstructure:"baseImageId"`
	AvailabilityDomain string `mapstructure:"availabilityDomain"`
	SubnetId           string `mapstructure:"subnetId"`
	Shape              string `mapstructure:"shape"`
}

type DebugConfig struct {
	OraclePublicKeys []string `mapstructure:"oracle_ssh_public_keys"`
}

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
	}, raws...)
	if err != nil {
		return nil, err
	}
	c := &b.config
	if c.Comm.SSHUsername == "" {
		c.Comm.SSHUsername = "opc"
	}
	var errs *packer.MultiError
	if es := c.Comm.Prepare(&c.ctx); len(es) > 0 {
		errs = packer.MultiErrorAppend(errs, es...)
	}
	if errs != nil && len(errs.Errors) > 0 {
		return nil, errs
	}
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

	oracleConfig := oraclebmc_sdk.NewConfig(config.UserId, config.TenantId, config.Fingerprint, config.PrivateKeyContents)

	computeApi := oraclebmc_sdk.ComputeApi{Config: oracleConfig}

	// Setup the state bag and initial state for the steps
	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("api", computeApi)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		&StepCreateInstance{
			availablityZone: config.AvailabilityDomain,
			baseImageId:     config.BaseImageId,
			compartmentId:   config.CompartmentId,
			shape:           config.Shape,
			subnetId:        config.SubnetId,
		},
		&communicator.StepConnect{
			Config:    &b.config.Comm,
			Host:      commHost,
			SSHConfig: sshConfig,
		},
		&StepProvision{},
		&StepCreateImage{
			DisplayName:   config.ImageDisplayName,
			CompartmentId: config.CompartmentId,
		},
		&StepTerminateInstance{},
	}

	b.runner = common.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	image := state.Get("image").(*oraclebmc_sdk.Image)

	artifact := &ImageArtifact{
		builderId:   BuilderId,
		displayName: image.DisplayName,
		imageId:     image.Id,
	}

	return artifact, nil
}
