resource "google_cloud_run_v2_service" "backend" {
  name     = "swim-rag-backend-go"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_INTERNAL_ONLY"

  deletion_protection = false # set to "true" in production

  template {
    service_account = google_service_account.swim_gen_backend_sa.email
    containers {
      image = "${google_artifact_registry_repository.docker.location}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.docker.repository_id}/swim-rag-backend:latest"
      liveness_probe {
        http_get {
          path = "/health"
          port = 8080 # Must match the container port
        }
        initial_delay_seconds = 5  # Optional: Delay before the first probe
        period_seconds        = 10 # Optional: How often (in seconds) to perform the probe
        timeout_seconds       = 2  # Optional: Probe timeout in seconds
        failure_threshold     = 3  # Optional: Number of consecutive failures before considering the container unhealthy
      }


      env {
        name = "DB_USER"
        value_source {
          secret_key_ref {
            secret  = data.google_secret_manager_secret.dbuser.secret_id
            version = data.google_secret_manager_secret_version.dbuser_data.name
          }
        }
      }
      # Mount secrets as environment variables
      env {
        name = "DB_PASS"
        value_source {
          secret_key_ref {
            secret  = data.google_secret_manager_secret_version_access.basic.secret
            version = "latest"
          }
        }
      }
      env {
        name  = "DB_NAME"
        value = "swim-db"
      }

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }
    }
    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.default.connection_name]
      }
      # Configure the liveness probe

    }
  }
  client     = "terraform"
  depends_on = [google_service_account.swim_gen_backend_sa]
}

resource "google_cloud_run_v2_service_iam_member" "bff_invoker" {
  project  = var.project_id
  location = var.region
  name     = google_cloud_run_v2_service.backend.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.swim_gen_frontend_sa.email}"
}

resource "google_cloud_run_v2_service" "bff" {
  name     = "swim-rag-bff"
  location = var.region

  deletion_protection = false

  template {
    service_account = google_service_account.swim_gen_frontend_sa.email
    containers {
      image = "${google_artifact_registry_repository.docker.location}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.docker.repository_id}/swim-rag-bff:latest"
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
      env {
        name  = "BACKEND_URL"
        value = google_cloud_run_v2_service.backend.uri
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }

  depends_on = [
    google_cloud_run_v2_service.backend,
    google_service_account.swim_gen_frontend_sa,
    google_cloud_run_v2_service_iam_member.bff_invoker
  ]
}

resource "google_cloud_run_v2_service" "frontend" {
  name     = "swim-rag-frontend"
  location = var.region

  deletion_protection = false

  template {
    service_account = google_service_account.swim_gen_frontend_sa.email
    containers {
      image = "${google_artifact_registry_repository.docker.location}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.docker.repository_id}/swim-rag-frontend:latest"
      liveness_probe {
        http_get {
          path = "/health"
          port = 80 # Nginx serves on port 80
        }
        initial_delay_seconds = 5
        period_seconds        = 10
        timeout_seconds       = 2
        failure_threshold     = 3
      }
      env {
        name  = "VITE_APP_API_URL"
        value = google_cloud_run_v2_service.bff.uri # Use the bff service URI
      }
    }
  }

  traffic {
    type    = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
    percent = 100
  }

  depends_on = [google_cloud_run_v2_service.bff] # Ensure bff is deployed first
}
