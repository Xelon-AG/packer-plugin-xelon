Type: `xelon`
Artifact BuilderId: `packer.builder.xelon`

The `xelon` Packer builder is able to create
Xelon [templates](https://www.xelon.ch/en/docs/working-with-templates-in-hq)
for use with [Xelon HQ](https://www.xelon.ch/hq) based on existing templates.

This builder builds a Xelon template by launching a device from a source template,
provisioning that running machine, and then creating a template from that machine.
This is all done in your own Xelon account. The builder will create temporary SSH
keys, etc. that provide it temporary access to the device while the template is
being created. This simplifies configuration quite a bit.

The builder does _not_ manage Xelon templates. Once it creates a template and
stores it in your account, it is up to you to use, delete, etc. the template.

-> **Note:** Temporary resources are, by default, all created with the
prefix `packer`.

## Configuration Reference

Configuration options are organized below into two categories: required and
optional. Within each category, the available options are alphabetized and
described.

### Access Configuration

**Required:**

<!-- Code generated from the comments of the AccessConfig struct in internal/builder/config.go; DO NOT EDIT MANUALLY -->

- `client_id` (string) - The client ID for IP ranges.
  Alternatively, can be configured using the `XELON_CLIENT_ID` environment variable.

- `token` (string) - The Xelon access token.
  Alternatively, can be configured using the `XELON_TOKEN` environment variable.

<!-- End of code generated from the comments of the AccessConfig struct in internal/builder/config.go; -->


**Optional:**

<!-- Code generated from the comments of the AccessConfig struct in internal/builder/config.go; DO NOT EDIT MANUALLY -->

- `base_url` (string) - The base URL endpoint for Xelon HQ. Default is `https://hq.xelon.ch/api/v2/`.
  Alternatively, can be configured using the `XELON_BASE_URL` environment variable.

<!-- End of code generated from the comments of the AccessConfig struct in internal/builder/config.go; -->


### Device Configuration

**Required:**

<!-- Code generated from the comments of the DeviceConfig struct in internal/builder/config.go; DO NOT EDIT MANUALLY -->

- `network_id` (string) - The ID of the network the builder device is launched into. The network must exist
  before the build runs and must allow outbound internet access, the SSH/WinRM endpoint
  Packer connects to, and any package mirrors the provisioners use.

- `source_template_id` (string) - The ID of the Xelon template to use to create the new template from.

<!-- End of code generated from the comments of the DeviceConfig struct in internal/builder/config.go; -->


**Optional:**

<!-- Code generated from the comments of the DeviceConfig struct in internal/builder/config.go; DO NOT EDIT MANUALLY -->

- `device_cpu_core_count` (int) - The number of CPU cores to allocate to the builder device. Defaults to `2`.

- `device_memory_gb` (int) - The amount of RAM in GB to allocate to the builder device. Defaults to `4`.

- `boot_disk_size_gb` (int) - The size of the builder device's boot disk in GB. The disk is created from the source
  template and is what the OS runs on during the build. Defaults to `10`.

- `swap_disk_size_gb` (int) - The size of the builder device's swap disk in GB. Defaults to `1`.

- `device_name` (string) - The hostname and display name of the builder device. Defaults to `packer-{{ uuid }}`.

- `admin_password` (string) - Password for root or administrative user, set at device creation. Must satisfy password
  complexity requirements (6+ chars, mixed case and digit).

<!-- End of code generated from the comments of the DeviceConfig struct in internal/builder/config.go; -->


### Template Configuration

**Required:**

<!-- Code generated from the comments of the Config struct in internal/builder/config.go; DO NOT EDIT MANUALLY -->

- `tenant_id` (string) - The ID of the Xelon tenant to whom the device and template belongs.

<!-- End of code generated from the comments of the Config struct in internal/builder/config.go; -->


**Optional:**

<!-- Code generated from the comments of the Config struct in internal/builder/config.go; DO NOT EDIT MANUALLY -->

- `skip_create_template` (bool) - If true, Packer will not create the Xelon template. Useful for setting to `true`
  during a build test stage. Defaults to `false`.

<!-- End of code generated from the comments of the Config struct in internal/builder/config.go; -->

<!-- Code generated from the comments of the TemplateConfig struct in internal/builder/config.go; DO NOT EDIT MANUALLY -->

- `template_name` (string) - The name of the resulting Xelon template that will appear when managing
  templates in the Xelon HQ console or via APIs. Defaults to `packer-{{ timestamp }}`.

- `template_description` (string) - The description to set for the resulting template. Defaults to `Created by Packer`.

<!-- End of code generated from the comments of the TemplateConfig struct in internal/builder/config.go; -->


### Example Usage

**HCL2**

```hcl
source "xelon" "basic" {
  client_id = var.xelon_client_id
  token     = var.xelon_token

  source_template_id = "<source-template-id>"
  network_id         = "<network-id>"
  admin_password     = var.admin_password
  ssh_username       = "ubuntu"

  tenant_id     = "<tenant-id>"
  template_name = "packer_Xelon_example_{{timestamp}}"
}

build {
  sources = ["source.xelon.basic"]
}
```
