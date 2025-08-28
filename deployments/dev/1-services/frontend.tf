locals {
  frontend_env_variables = {
    PROJECT_ID = var.project_id
    REGION     = var.region
    BFF_URL    = google_cloud_run_v2_service.bff.uri
  }
}

data "google_artifact_registry_docker_image" "frontend_image" {
  location      = data.google_artifact_registry_repository.docker.location
  repository_id = data.google_artifact_registry_repository.docker.repository_id
  image_name    = "swim-gen-frontend:${var.frontend_image_tag}"
}

resource "google_cloud_run_v2_service" "frontend" {
  name     = "swim-gen-frontend"
  location = var.region

  # gcloud command used --no-allow-unauthenticated (public ingress but no allUsers binding)
  # So expose all ingress, rely on IAM to restrict.
  ingress              = "INGRESS_TRAFFIC_ALL"
  invoker_iam_disabled = true

  # Only allow authenticated invocations from other services with the "run.invoker" role
  # (e.g. the frontend service).
  # This is enforced via IAM bindings in the frontend.tf file.
  # See: https://cloud.google.com/run/docs/securing/service-identity#granting_other_identities_access_to_your_service
  deletion_protection = false

  template {
    service_account                  = var.iam.swim_gen_frontend.email
    max_instance_request_concurrency = 200
    timeout                          = "60s"
    containers {
      image = data.google_artifact_registry_docker_image.frontend_image.self_link
      resources {
        limits = {
          cpu    = 1
          memory = "512Mi"
        }
        cpu_idle          = true
        startup_cpu_boost = true
      }
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
        for_each = local.frontend_env_variables
        content {
          name  = env.key
          value = tostring(env.value)
        }
      }
    }
  }

  client     = "terraform"
  depends_on = [google_cloud_run_v2_service.backend, google_cloud_run_v2_service.frontend]
}
