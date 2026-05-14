package builder

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"

	"github.com/Xelon-AG/xelon-sdk-go/xelon"
)

var _ multistep.Step = (*stepCreateTemplate)(nil)

type stepCreateTemplate struct {
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

	return multistep.ActionContinue
}

func (s *stepCreateTemplate) Cleanup(_ multistep.StateBag) {
	return
}
