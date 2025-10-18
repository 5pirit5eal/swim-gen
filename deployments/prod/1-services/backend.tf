locals {
  backend_env_variables = {
    PROJECT_ID       = var.project_id
    REGION           = var.region
    DB_NAME          = var.supabase.dbname
    DB_HOST          = var.supabase.host
    DB_PORT          = var.supabase.port
    DB_USER          = var.supabase.backend_user
    DB_PASS_LOCATION = var.secret_version_ids.dbpassword_user
    EMBEDDING_MODEL  = var.embedding_model
    EMBEDDING_NAME   = var.embedding_name
    EMBEDDING_SIZE   = var.embedding_size
    MODEL            = var.model
    SMALL_MODEL      = var.small_model # fixed key (was SMALl_MODEL)
    LOG_LEVEL        = var.log_level
    BUCKET_NAME      = var.bucket_name
    SIGNING_SA       = var.iam.pdf_export.email
  }
}

data "google_artifact_registry_docker_image" "backend_image" {
  location      = data.google_artifact_registry_repository.docker.location
  repository_id = data.google_artifact_registry_repository.docker.repository_id
  image_name    = "swim-gen-backend:${var.version_tag}"
}

resource "google_cloud_run_v2_service" "backend" {
  name     = "swim-gen-backend"
  location = var.region

  # gcloud command used --no-allow-unauthenticated (public ingress but no allUsers binding)
  # So expose all ingress, rely on IAM to restrict.
  ingress              = "INGRESS_TRAFFIC_ALL"
  invoker_iam_disabled = false

  # Only allow authenticated invocations from other services with the "run.invoker" role
  # (e.g. the BFF service).
  # This is enforced via IAM bindings in the bff.tf file.
  # See: https://cloud.google.com/run/docs/securing/service-identity#granting_other_identities_access_to_your_service
  deletion_protection = false

  # Set the number of maximum instances to control costs
  scaling {
    max_instance_count = 3
  }

  template {
    service_account                  = var.iam.swim_gen_backend.email
    max_instance_request_concurrency = 200
    timeout                          = "3600s"
    scaling {
      min_instance_count = 0
      max_instance_count = 15
    }
    containers {
      image = data.google_artifact_registry_docker_image.backend_image.self_link
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
        for_each = local.backend_env_variables
        content {
          name  = env.key
          value = tostring(env.value)
        }
      }
    }
  }

  client     = "terraform"
  depends_on = []
}
