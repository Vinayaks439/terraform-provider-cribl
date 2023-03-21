terraform {
  required_providers {
    cribl = {
      source = "cribl.com/criblprovider/cribl"
      version = "1.0.5"
    }
  }
}

provider "cribl" {
  host = "http://127.0.0.1:19000"
  username = "admin"
  password = var.password
}