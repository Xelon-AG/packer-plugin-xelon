//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type NetworkConfig,NetworkDatasourceOutput
package datasource

import (
	"context"
	"fmt"
	"log"
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
	_ packer.Datasource = (*NetworkDatasource)(nil)
)

type NetworkDatasource struct {
	config NetworkConfig

	ctx interpolate.Context
}

type NetworkConfig struct {
	common.PackerConfig `mapstructure:",squash"`
	config.AccessConfig `mapstructure:",squash"`

	// The ID of the cloud.
	CloudID string `mapstructure:"cloud_id" required:"false"`
	// The network name.
	Name string `mapstructure:"name" required:"true"`
	// The ID of the Xelon tenant to whom the network belongs.
	TenantID string `mapstructure:"tenant_id" required:"false"`
}

func (d *NetworkDatasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *NetworkDatasource) Configure(raws ...any) error {
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

type NetworkDatasourceOutput struct {
	// The ID of the network.
	ID string `mapstructure:"id"`
	// The network type.
	Type string `mapstructure:"type"`
}

func (d *NetworkDatasource) OutputSpec() hcldec.ObjectSpec {
	return (&NetworkDatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *NetworkDatasource) Execute() (cty.Value, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client := xelonapi.NewXelonClient(d.config.AccessConfig)
	networkName := d.config.Name
	networkCloudID := d.config.CloudID
	networkTenantID := d.config.TenantID
	ctyNull := cty.NullVal(cty.EmptyObject)

	log.Printf("[DEBUG] Searching for network by name: %v", networkName)
	networks, _, err := client.Networks.List(ctx, &xelon.NetworkListOptions{Search: networkName})
	if err != nil {
		return ctyNull, fmt.Errorf("failed to list networks: %s", err)
	}
	log.Printf("[DEBUG] Got networks. Count: %d", len(networks))

	// TODO: add check if cloud_id and tenant_id are set
	if networkCloudID != "" || networkTenantID != "" {
		log.Printf("[WARN] Additional parameters are not evaluated yet: cloud_id=%s, tenant_id=%s", networkCloudID, networkTenantID)
	}

	if len(networks) == 0 {
		return ctyNull, fmt.Errorf("no networks was found matching name: %v", networkName)
	}
	if len(networks) > 1 {
		return ctyNull, fmt.Errorf("found more than one network, refine your search using cloud_id or/and tenant_id")
	}
	network := networks[0]

	output := &NetworkDatasourceOutput{
		ID:   network.ID,
		Type: network.Type,
	}
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
