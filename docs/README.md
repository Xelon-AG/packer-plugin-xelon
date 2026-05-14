The [Xelon](https://www.xelon.ch/) Packer plugin provides a builder for building templates in Xelon HQ.

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [
`packer init`](https://www.packer.io/docs/commands/init).

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
$ packer plugins install github.com/Xelon-AG/packer-plugin-xelon
```

### Components

The Scaffolding plugin is intended as a starting point for creating Packer plugins

#### Builders

- [builder](/packer/integrations/hashicorp/scaffolding/latest/components/builder/builder-name) - The scaffolding builder
  is used to create endless Packer
  plugins using a consistent plugin structure.

#### Provisioners

- [provisioner](/packer/integrations/hashicorp/scaffolding/latest/components/provisioner/provisioner-name) - The
  scaffolding provisioner is used to provisioner
  Packer builds.

#### Post-processors

- [post-processor](/packer/integrations/hashicorp/scaffolding/latest/components/post-processor/postprocessor-name) - The
  scaffolding post-processor is used to
  export scaffolding builds.

#### Data Sources

- [data source](/packer/integrations/hashicorp/scaffolding/latest/components/datasource/datasource-name) - The
  scaffolding data source is used to
  export scaffolding data.

