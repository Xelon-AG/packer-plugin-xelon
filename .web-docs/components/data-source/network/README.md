Type: `xelon-network`

The Xelon Network data source provides information about a network in Xelon HQ.

-> **Note:** Data sources is a feature exclusively available to HCL2 templates.

Basic example of usage:

```hcl
data "xelon-network" "basic-example" {
  name = "packer-builder-wan"
}

# usage example of the data source output
locals {
  network_id   = data.xelon-network.basic-example.id
  network_type = data.xelon-network.basic-example.type
}
```

## Configuration Reference

Configuration options are organized below into two categories: required and optional. Within each category, the
available options are alphabetized and described.

**Required:**

<!-- Code generated from the comments of the NetworkConfig struct in internal/datasource/data_network.go; DO NOT EDIT MANUALLY -->

- `name` (string) - The network name.

<!-- End of code generated from the comments of the NetworkConfig struct in internal/datasource/data_network.go; -->


**Optional:**

<!-- Code generated from the comments of the NetworkConfig struct in internal/datasource/data_network.go; DO NOT EDIT MANUALLY -->

- `cloud_id` (string) - The ID of the cloud.

- `tenant_id` (string) - The ID of the Xelon tenant to whom the network belongs.

<!-- End of code generated from the comments of the NetworkConfig struct in internal/datasource/data_network.go; -->


## Output Data

<!-- Code generated from the comments of the NetworkDatasourceOutput struct in internal/datasource/data_network.go; DO NOT EDIT MANUALLY -->

- `id` (string) - The ID of the network.

- `type` (string) - The network type.

<!-- End of code generated from the comments of the NetworkDatasourceOutput struct in internal/datasource/data_network.go; -->


## Authentication

This data source uses the same authentication method as the main Xelon plugin to connect to Xelon HQ. See the of
the [authentication section](/packer/integrations/Xelon-AG/xelon#authentication) plugin’s documentation for details.
