terraform {
  required_providers {
    harperdb = {
      source = "registry.terraform.io/harperdb/harperdb"
    }
  }
}

provider "harperdb" {
  endpoint = "https://my-intance.harperdbcloud.com"
  username = "ada"
  password = "lovelace"
}
