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

- [xelon](/packer/integrations/Xelon-AG/xelon/latest/components/builder/xelon) - Create Xelon templates by launching
  device from a source template and re-packaging it into a new template after provisioning.
