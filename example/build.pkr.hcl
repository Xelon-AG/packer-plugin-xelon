packer {
  required_plugins {
    xelon = {
      source  = "github.com/Xelon-AG/xelon"
      version = ">= 1"
    }
  }
}

source "xelon" "example" {
  client_id = "YOUR CLIENT ID"
  token     = "YOUR API TOKEN"

  tenant_id          = "YOUR TENANT ID"
  source_template_id = "SOURCE TEMPLATE ID"
  network_id         = "NETWORK ID"
  admin_password     = "<secure-password-for-admin-user>"

  ssh_username = "root"
}

build {
  sources = ["source.xelon.example"]
}
