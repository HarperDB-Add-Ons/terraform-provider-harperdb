terraform {
  required_providers {
    harperdb = {
      source = "registry.terraform.io/harperdb/harperdb"
    }
  }
}

provider "harperdb" {
  endpoint = "https://terraform-test-moredhel.harperdbcloud.com"
  username = "moredhel"
  password = "question"
}
