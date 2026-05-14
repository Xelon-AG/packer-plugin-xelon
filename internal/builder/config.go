//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package builder

import (
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	Comm                communicator.Config `mapstructure:",squash"`

	AccessConfig   `mapstructure:",squash"`
	DeviceConfig   `mapstructure:",squash"`
	TemplateConfig `mapstructure:",squash"`

	// If true, Packer will not create the Xelon template. Useful for setting to `true`
	// during a build test stage. Defaults to `false`.
	SkipCreateTemplate bool `mapstructure:"skip_create_template" required:"false"`
	// The ID of the Xelon tenant to whom the device and template belongs.
	TenantID string `mapstructure:"tenant_id" required:"true"`

	ctx interpolate.Context
}

// AccessConfig is for common configuration related to Xelon HQ API access.
type AccessConfig struct {
	// The base URL endpoint for Xelon HQ. Default is `https://hq.xelon.ch/api/v2/`.
	// Alternatively, can be configured using the `XELON_BASE_URL` environment variable.
	BaseURL string `mapstructure:"base_url" required:"false"`
	// The client ID for IP ranges.
	// Alternatively, can be configured using the `XELON_CLIENT_ID` environment variable.
	ClientID string `mapstructure:"client_id" required:"true"`
	// The Xelon access token.
	// Alternatively, can be configured using the `XELON_TOKEN` environment variable.
	Token string `mapstructure:"token" required:"true"`
}

func (c *AccessConfig) Prepare(_ *interpolate.Context, _ ...any) *packer.MultiError {
	return nil
}

// DeviceConfig contains configuration for running a Xelon device from a source template.
type DeviceConfig struct {
	// The number of CPU cores to allocate to the builder device. Defaults to `2`.
	DeviceCPUCores int `mapstructure:"device_cpu_core_count" required:"false"`
	// The amount of RAM in GB to allocate to the builder device. Defaults to `4`.
	DeviceMemoryGB int `mapstructure:"device_memory_gb" required:"false"`
	// The size of the builder device's boot disk in GB. The disk is created from the source
	// template and is what the OS runs on during the build. Defaults to `10`.
	BootDiskSizeGB int `mapstructure:"boot_disk_size_gb" required:"false"`
	// The size of the builder device's swap disk in GB. Defaults to `1`.
	SwapDiskSizeGB int `mapstructure:"swap_disk_size_gb" required:"false"`
	// The hostname and display name of the builder device. Defaults to `packer-{{ uuid }}`.
	DeviceName string `mapstructure:"device_name" required:"false"`
	// The ID of the network the builder device is launched into. The network must exist
	// before the build runs and must allow outbound internet access, the SSH/WinRM endpoint
	// Packer connects to, and any package mirrors the provisioners use.
	NetworkID string `mapstructure:"network_id" required:"true"`
	// Password for root or administrative user, set at device creation. Must satisfy password
	// complexity requirements (6+ chars, mixed case and digit).
	AdminPassword string `mapstructure:"admin_password" required:"false"`
	// The ID of the Xelon template to use to create the new template from.
	SourceTemplateID string `mapstructure:"source_template_id" required:"true"`
}

func (c *DeviceConfig) Prepare(_ *interpolate.Context, _ ...any) *packer.MultiError {
	var errs *packer.MultiError

	if c.DeviceCPUCores == 0 {
		c.DeviceCPUCores = 2
	}

	if c.DeviceMemoryGB == 0 {
		c.DeviceMemoryGB = 4
	}

	if c.BootDiskSizeGB == 0 {
		c.BootDiskSizeGB = 10
	}

	if c.SwapDiskSizeGB == 0 {
		c.SwapDiskSizeGB = 1
	}

	if c.DeviceName == "" {
		c.DeviceName = fmt.Sprintf("packer-%s", uuid.TimeOrderedUUID())
	}

	return errs
}

// TemplateConfig is for common configuration related to creating Xelon templates.
type TemplateConfig struct {
	// The name of the resulting Xelon template that will appear when managing
	// templates in the Xelon HQ console or via APIs. Defaults to `packer-{{ timestamp }}`.
	TemplateName string `mapstructure:"template_name" required:"false"`
	// The description to set for the resulting template. Defaults to `Created by Packer`.
	TemplateDescription string `mapstructure:"template_description" required:"false"`
}

func (c *TemplateConfig) Prepare(ctx *interpolate.Context, raws ...any) *packer.MultiError {
	var errs *packer.MultiError

	if c.TemplateName == "" {
		templateName, err := interpolate.Render("packer-{{ timestamp }}", ctx)
		if err != nil {
			errs = packer.MultiErrorAppend(errs, fmt.Errorf("failed to render default template name: %v", err))
		} else {
			c.TemplateName = templateName
		}
	}

	if c.TemplateDescription == "" {
		c.TemplateDescription = "Created by Packer"
	}

	return errs
}

func (c *Config) Prepare(raws ...any) *packer.MultiError {
	var errs *packer.MultiError
	err := config.Decode(c, &config.DecodeOpts{
		PluginType:         PluginBuilderID,
		Interpolate:        true,
		InterpolateContext: &c.ctx,
	}, raws...)
	if err != nil {
		errs = packer.MultiErrorAppend(errs, err)
		return errs
	}

	errs = packer.MultiErrorAppend(errs, c.AccessConfig.Prepare(&c.ctx, raws))
	errs = packer.MultiErrorAppend(errs, c.DeviceConfig.Prepare(&c.ctx, raws))
	errs = packer.MultiErrorAppend(errs, c.TemplateConfig.Prepare(&c.ctx, raws))
	errs = packer.MultiErrorAppend(errs, c.Comm.Prepare(&c.ctx)...)

	return errs
}
