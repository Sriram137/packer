package oracle

import (
	"github.com/elricL/oracle_bmc_sdk"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
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

type Builder struct {
	config *Config
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
	//decoded_str, err := strconv.Unquote(b.config.PrivateKeyContents)
	//if err != nil {
	//	return nil, err
	//}
	//a, _ := ioutil.ReadFile("/Users/sriramg/.ssh/id_rsa.pub")
	//log.Println("&&&&&&&&&&&&&&&&")
	//log.Println(a)
	//log.Println("****************")
	//b.config.PrivateKeyContents = decoded_str
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
			compartmentId:   config.AvailabilityDomain,
			shape:           config.Shape,
			subnetId:        config.SubnetId,
		},
		&StepCreateImage{
			DisplayName:   config.ImageDisplayName,
			CompartmentId: config.CompartmentId,
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
