resource "google_cloud_run_v2_service" "default" {
  name     = "cloudrun-service"
  location = var.region


  deletion_protection = false # set to "true" in production

  template {
    service_account = google_service_account.service_account.email
    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello:latest" # Image to deploy
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
            secret  = google_secret_manager_secret.dbuser.secret_id
            version = google_secret_manager_secret_version.dbuser_data.name
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
  depends_on = [google_service_account.service_account]
}