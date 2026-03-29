//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package builder

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

var (
	_ packer.Builder = (*Builder)(nil)
)

// PluginBuilderID is the unique ID for the builder
const PluginBuilderID = "packer.builder.xelon"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`

	// The base URL endpoint for Xelon HQ. Default is `https://hq.xelon.ch/api/v2/`.
	// Alternatively, can be configured using the `XELON_BASE_URL` environment variable.
	BaseURL string `mapstructure:"base_url" required:"false"`
	// The client ID for IP ranges.
	// Alternatively, can be configured using the `XELON_CLIENT_ID` environment variable.
	ClientID string `mapstructure:"client_id" required:"true"`
	// The Xelon access token.
	// Alternatively, can be configured using the `XELON_TOKEN` environment variable.
	Token string `mapstructure:"token" required:"true"`

	// Skips template creation.
	SkipCreateTemplate bool `mapstructure:"skip_create_template" required:"false"`

	// The hostname and display name of the device.
	DeviceName string `mapstructure:"device_name" required:"false"`
	// The ID of the network configured for the device.
	NetworkID string `mapstructure:"network_id" required:"true"`
	// The ID of the Xelon template to launch device from.
	SourceTemplateID string `mapstructure:"source_template_id" required:"true"`
	// The ID of the Xelon tenant to whom the device and template belongs.
	TenantID string `mapstructure:"tenant_id" required:"true"`

	// The name of the resulting Xelon template that will appear when managing
	// templates in the Xelon HQ console or via APIs.
	TemplateName string `mapstructure:"template_name" required:"true"`
	// The description to set for the resulting template. By default, this description is empty.
	TemplateDescription string `mapstructure:"template_description" required:"false"`

	ctx interpolate.Context
}

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...any) (generatedVars []string, warnings []string, err error) {
	err = config.Decode(&b.config, &config.DecodeOpts{
		PluginType:         PluginBuilderID,
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
	}, raws...)
	if err != nil {
		return nil, nil, err
	}

	var errs *packer.MultiError
	errs = packer.MultiErrorAppend(errs, b.config.Comm.Prepare(&b.config.ctx)...)

	if errs != nil && len(errs.Errors) > 0 {
		return nil, nil, errs
	}

	packer.LogSecretFilter.Set(b.config.Token)

	return []string{}, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	client, err := newXelonClient(b.config)
	if err != nil {
		return nil, err
	}

	// set up the state bag
	state := new(multistep.BasicStateBag)
	state.Put("client", client)
	state.Put("config", &b.config)
	state.Put("hook", hook)
	state.Put("ui", ui)

	// build steps
	steps := []multistep.Step{
		&stepCreateSSHKey{
			Debug:        b.config.PackerDebug,
			DebugKeyPath: fmt.Sprintf("xelon_%s.pem", b.config.PackerBuildName),
		},
		&stepCreateDevice{},
		&communicator.StepConnect{
			Config:    &b.config.Comm,
			Host:      communicator.CommHost(b.config.Comm.Host(), "device_ip"),
			SSHConfig: b.config.Comm.SSHConfigFunc(),
		},
		&commonsteps.StepProvision{},
		&stepCreateTemplate{
			SkipCreateTemplate: b.config.SkipCreateTemplate,
		},
	}

	// run
	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	templateID := state.Get("template_id").(string)
	templateName := state.Get("template_name").(string)

	artifact := &Artifact{
		Client:       client,
		TemplateID:   templateID,
		TemplateName: templateName,
		StateData:    map[string]any{"one": "two"},
	}

	return artifact, nil
}
