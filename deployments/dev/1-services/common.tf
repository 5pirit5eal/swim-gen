data "google_artifact_registry_repository" "docker" {
  location      = var.artifactregistry.location
  repository_id = var.artifactregistry.repository
}
