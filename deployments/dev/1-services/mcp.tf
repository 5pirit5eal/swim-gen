locals {
  mcp_env_variables = {
    PROJECT_ID       = var.project_id
    REGION           = var.region
    LOG_LEVEL        = var.log_level
    SWIM_RAG_API_URL = google_cloud_run_v2_service.backend.uri
  }
}

data "google_artifact_registry_docker_image" "mcp_image" {
  location      = data.google_artifact_registry_repository.docker.location
  repository_id = data.google_artifact_registry_repository.docker.repository_id
  image_name    = "swim-gen-mcp:${var.mcp_server_image_tag}"
}

resource "google_cloud_run_v2_service" "mcp" {
  name     = "swim-gen-mcp"
  location = var.region

  # gcloud command used --no-allow-unauthenticated (public ingress but no allUsers binding)
  # So expose all ingress, rely on IAM to restrict.
  ingress              = "INGRESS_TRAFFIC_ALL"
  invoker_iam_disabled = false

  # Only allow authenticated invocations from other services with the "run.invoker" role
  # (e.g. the mcp service).
  # This is enforced via IAM bindings in the mcp.tf file.
  # See: https://cloud.google.com/run/docs/securing/service-identity#granting_other_identities_access_to_your_service
  deletion_protection = false

  template {
    service_account = var.iam.swim_gen_frontend.email
    timeout         = 600
    containers {
      image = data.google_artifact_registry_docker_image.mcp_image.self_link

      liveness_probe {
        http_get {
          path = "/health"
          port = 8080
        }
        initial_delay_seconds = 5
        period_seconds        = 10
        timeout_seconds       = 2
        failure_threshold     = 3
      }

      dynamic "env" {
        for_each = local.mcp_env_variables
        content {
          name  = env.key
          value = tostring(env.value)
        }
      }
    }
  }

  client     = "terraform"
  depends_on = [google_cloud_run_v2_service.backend]
}
