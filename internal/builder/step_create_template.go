package builder

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/packerbuilderdata"

	"github.com/Xelon-AG/xelon-sdk-go/xelon"
)

var _ multistep.Step = (*stepCreateTemplate)(nil)

type stepCreateTemplate struct {
	GeneratedData      *packerbuilderdata.GeneratedData
	SkipCreateTemplate bool
}

func (s *stepCreateTemplate) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("client").(*xelon.Client)
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	deviceID := state.Get("device_id").(string)

	if s.SkipCreateTemplate {
		ui.Say("Skipping template creation...")
		return multistep.ActionContinue
	}

	ui.Sayf("Creating template from device %s", deviceID)

	template, _, err := client.Templates.Create(ctx, &xelon.TemplateCreateRequest{
		Description: config.TemplateDescription,
		DeviceID:    deviceID,
		Name:        config.TemplateName,
		SendEmail:   false,
		TenantID:    config.TenantID,
	})
	if err != nil {
		state.Put("error", err)
		return multistep.ActionHalt
	}

	state.Put("template_id", template.ID)
	state.Put("template_name", template.Name)

	// provision generated_data from declared in Builder.Prepare func
	// see doc https://www.packer.io/docs/extending/custom-builders#build-variables for details
	s.GeneratedData.Put("TemplateID", template.ID)
	s.GeneratedData.Put("TemplateName", template.Name)

	return multistep.ActionContinue
}

func (s *stepCreateTemplate) Cleanup(_ multistep.StateBag) {}
