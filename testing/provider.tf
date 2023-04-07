terraform {
  required_providers {
    harperdb = {
      source = "registry.terraform.io/harperdb/harperdb"
    }
  }
}

provider "harperdb" {
  endpoint = var.db_endpoint
  username = var.db_username
  password = var.db_password
}
