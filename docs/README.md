The [Xelon](https://www.xelon.ch/) Packer plugin provides a builder for building templates in Xelon HQ.

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run
[`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    name = {
      source  = "github.com/Xelon-AG/xelon"
      version = "~> 1"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/Xelon-AG/xelon
```

### Components

#### Builders

- [xelon](/packer/integrations/Xelon-AG/xelon/latest/components/builder/xelon) - The xelon builder creates new templates
  from existing ones by launching a device based on a source template, provisioning that device, and exporting it as a
  reusable template.

### Authentication

Authentication with Xelon HQ requires a Client ID and an Access Token. You can obtain both by creating a Service Token
in your [user profile](https://api-v2-developers.xelon.ch/#section/Introduction/Authorization).

The following options are available for the `xelon` builder and the `xelon-network` data source:

#### Required

@include 'internal/AccessConfig-required.mdx'

#### Optional

@include 'internal/AccessConfig-not-required.mdx'
