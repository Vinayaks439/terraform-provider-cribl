terraform {
  required_providers {
    cribl = {
      source = "cribl.com/criblprovider/cribl"
      version = "1.0.4"
    }
  }
}

provider "cribl" {
  host = "localhost:19000"
  username = "admin"
  password = var.password
}