data "google_secret_manager_secret" "dbpassword-root" {
  secret_id = "db-password-root"
  project   = var.project_id
}

# Manually created secret to keep the database password out of the tf state
data "google_secret_manager_secret_version_access" "dbpassword" {
  secret     = "db-password-root"
  version    = "latest"
  project    = var.project_id
  depends_on = [google_project_service.apis]
}


resource "google_secret_manager_secret" "dbuser" {
  secret_id = "db-user-secret"
  replication {
    auto {}
  }
  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "dbuser_data" {
  secret      = google_secret_manager_secret.dbuser.id
  secret_data = var.dbuser
}



