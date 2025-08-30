locals {
  bff_env_variables = {
    PROJECT_ID   = var.project_id
    REGION       = var.region
    LOG_LEVEL    = var.log_level
    BACKEND_URL  = google_cloud_run_v2_service.backend.uri
    FRONTEND_URL = var.domain_url
  }
}

data "google_artifact_registry_docker_image" "bff_image" {
  location      = data.google_artifact_registry_repository.docker.location
  repository_id = data.google_artifact_registry_repository.docker.repository_id
  image_name    = "swim-gen-bff:${var.bff_image_tag}"
}

resource "google_cloud_run_v2_service" "bff" {
  name     = "swim-gen-bff"
  location = var.region

  # gcloud command used --no-allow-unauthenticated (public ingress but no allUsers binding)
  # So expose all ingress, rely on IAM to restrict.
  ingress              = "INGRESS_TRAFFIC_ALL"
  invoker_iam_disabled = true

  # Only allow authenticated invocations from other services with the "run.invoker" role
  # (e.g. the BFF service).
  # This is enforced via IAM bindings in the bff.tf file.
  # See: https://cloud.google.com/run/docs/securing/service-identity#granting_other_identities_access_to_your_service
  deletion_protection = false

  # Set the number of maximum instances to control costs
  scaling {
    max_instance_count = 1
  }

  template {
    service_account                  = var.iam.swim_gen_frontend.email
    session_affinity                 = true
    max_instance_request_concurrency = 200
    timeout                          = "600s"

    containers {
      image = data.google_artifact_registry_docker_image.bff_image.self_link
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
        for_each = local.bff_env_variables
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

resource "google_cloud_run_v2_service_iam_member" "backend_invoker" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.backend.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${var.iam.swim_gen_frontend.email}"
}
