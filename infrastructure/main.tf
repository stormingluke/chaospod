terraform {
  cloud {
    organization = "edgerx"

    workspaces {
      name = "controlplane"
    }
  }
}

variable "cloudbuild_trigger_name" {}

variable "artifact_registry_name" {}

variable "artifact_registry_description" {}

variable "GITHUB_OWNER" {}

variable "label_project" {}
