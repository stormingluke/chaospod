terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
        version = "~> 4"
    }
    github = {
      source  = "integrations/github"
      version = "~> 4"
    }
  }
}

variable "GOOGLE_CLOUD_KEYFILE_JSON" {}
variable "GCP_PROJECT" {}
variable "GITHUB_TOKEN" {}

provider "google" {
  project     = var.GCP_PROJECT
  region      = "europe-west4"
  zone        = "europe-west4"
  credentials = var.GOOGLE_CLOUD_KEYFILE_JSON
}

provider "github" {
  token = var.GITHUB_TOKEN
}
