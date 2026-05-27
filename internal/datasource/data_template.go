//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type TemplateConfig,TemplateDatasourceOutput
package datasource

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	packercfg "github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/zclconf/go-cty/cty"

	"github.com/Xelon-AG/packer-plugin-xelon/internal/config"
	"github.com/Xelon-AG/packer-plugin-xelon/internal/xelonapi"
	"github.com/Xelon-AG/xelon-sdk-go/xelon"
)

var (
	_ packer.Datasource = (*TemplateDatasource)(nil)
)

type TemplateDatasource struct {
	config TemplateConfig

	ctx interpolate.Context
}

type TemplateConfig struct {
	common.PackerConfig `mapstructure:",squash"`
	config.AccessConfig `mapstructure:",squash"`

	// If true, the most recent OS template will be returned. If false (default),
	// an error will be returned if more than one template matches the filters.
	MostRecent bool `mapstructure:"most_recent" required:"false"`
	// The template name.
	Name string `mapstructure:"name" required:"true"`
}

func (d *TemplateDatasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *TemplateDatasource) Configure(raws ...any) error {
	if err := packercfg.Decode(&d.config, &packercfg.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &d.ctx,
	}, raws...); err != nil {
		return err
	}

	var errs *packer.MultiError
	errs = packer.MultiErrorAppend(errs, d.config.Prepare(&d.ctx, raws))

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

type TemplateDatasourceOutput struct {
	// The ID of the template.
	ID string `mapstructure:"id"`
	// The name of the template.
	Name string `mapstructure:"name"`
	// The date of creation of the template.
	CreationDate string `mapstructure:"creation_date"`
}

func (d *TemplateDatasource) OutputSpec() hcldec.ObjectSpec {
	return (&TemplateDatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *TemplateDatasource) Execute() (cty.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := xelonapi.NewXelonClient(d.config.AccessConfig)
	templateName := d.config.Name
	ctyNull := cty.NullVal(cty.EmptyObject)

	log.Printf("[DEBUG]  Searching for template by name: %v", templateName)
	templates, _, err := client.Templates.List(ctx, &xelon.TemplateListOptions{Search: templateName})
	if err != nil {
		return ctyNull, fmt.Errorf("failed to list templates: %s", err)
	}
	log.Printf("[DEBUG]  Got templates with matching criteria: %v", templates)

	var template *xelon.Template

	if d.config.MostRecent {
		log.Printf("[INFO] Use most recent template")
		slices.SortFunc(templates, func(first, second xelon.Template) int {
			if first.CreatedAt == nil && second.CreatedAt != nil {
				return 1
			} else if first.CreatedAt != nil && second.CreatedAt == nil {
				return -1
			} else if first.CreatedAt == nil && second.CreatedAt == nil {
				return 0
			}
			return second.CreatedAt.Compare(*first.CreatedAt)
		})
	} else {
		if len(templates) == 0 {
			return ctyNull, fmt.Errorf("no templates was found matching name: %v", templateName)
		}
		if len(templates) > 1 {
			return ctyNull, fmt.Errorf("found more than one network, refine your search or use most_recent option")
		}
	}

	template = &templates[0]

	log.Printf("[INFO] Found template: %v", template)

	output := &TemplateDatasourceOutput{
		ID:           template.ID,
		Name:         template.Name,
		CreationDate: template.CreatedAt.Format(time.RFC3339),
	}
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
