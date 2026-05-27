data "xelon-template" "ubuntu" {
  name        = "Ubuntu 20.04 64 Bit EN"
  most_recent = true
}

locals {
  template_id = data.xelon-template.basic-example.id
}
