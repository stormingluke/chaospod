variable "artifact_registry_labels" {
  type = map(string)
  default = {
    finance   = "artifact_registry"
    technical = "docker_artifact_registry"
  }
}
# Cannot connect the github to cloudbuild without pressing the 'I Agree to the Terms' button in Cloudbuild
resource "google_artifact_registry_repository" "artifact_registry_controlplane" {
  labels = merge({"project": var.label_project}, var.artifact_registry_labels)
  location      = "europe-west4"
  repository_id = var.artifact_registry_name
  description   = var.artifact_registry_description
  format        = "DOCKER"
}

resource "google_cloudbuild_trigger" "include-build-logs-trigger" {
  name     = "controlplane-podchaos"
  filename = "cloudbuild.yaml"

  github {
    owner = var.GITHUB_OWNER
    name  = "podchaosmonkey"
    push {
      branch = "^main$"
    }
  }

  include_build_logs = "INCLUDE_BUILD_LOGS_WITH_STATUS"
}
