terraform {
  required_providers {
    tfe = {
      source  = "hashicorp/tfe"
      version = "0.42.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.4.3"
    }
  }
  required_version = ">= 1.1.0"
}
