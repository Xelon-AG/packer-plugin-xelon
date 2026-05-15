package builder

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

var (
	_ packer.Builder = (*Builder)(nil)
)

// PluginBuilderID is the unique ID for the builder.
const PluginBuilderID = "packer.builder.xelon"

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...any) (generatedVars []string, warnings []string, err error) {
	errs := b.config.Prepare(raws...)
	if errs != nil && len(errs.Errors) > 0 {
		return nil, nil, errs
	}

	packer.LogSecretFilter.Set(b.config.Token, b.config.AdminPassword)

	return []string{}, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	client := newXelonClient(b.config)

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
