# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name        = "Xelon"
  description = "The Xelon plugin can be used with HashiCorp Packer to create custom templates on Xelon Cloud."
  identifier  = "packer/Xelon-AG/xelon"
  component {
    type = "builder"
    name = "Xelon Cloud"
    slug = "xelon"
  }
}
