locals {
  secret_ids = {
    dbname          = google_secret_manager_secret.dbname.id
    dbuser          = google_secret_manager_secret.dbuser.id
    dbpassword_root = data.google_secret_manager_secret.dbpassword_root.id
    dbpassword_user = data.google_secret_manager_secret.dbpassword_user.id
  }
  secret_version_ids = {
    dbpassword_root = data.google_secret_manager_secret_version_access.dbpassword_root.id
    dbpassword_user = data.google_secret_manager_secret_version_access.dbpassword_user.id
  }
}

data "google_secret_manager_secret" "dbpassword_root" {
  secret_id = "db-password-root"
  project   = data.google_project.project.number
}

data "google_secret_manager_secret" "dbpassword_user" {
  secret_id = "db-password-${var.dbuser}"
  project   = data.google_project.project.number
}

# Manually created secret to keep the database password out of the tf state
data "google_secret_manager_secret_version_access" "dbpassword_root" {
  secret     = "db-password-root"
  version    = "latest"
  project    = data.google_project.project.number
  depends_on = [google_project_service.apis]
}

# Manually created secret to keep the database password out of the tf state
data "google_secret_manager_secret_version_access" "dbpassword_user" {
  secret     = "db-password-${var.dbuser}"
  version    = "latest"
  project    = data.google_project.project.number
  depends_on = [google_project_service.apis]
}


resource "google_secret_manager_secret" "dbuser" {
  secret_id = "db-user-secret"
  replication {
    auto {}
  }
  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "dbuser" {
  secret      = google_secret_manager_secret.dbuser.id
  secret_data = var.dbuser
}

resource "google_secret_manager_secret" "dbname" {
  secret_id = var.dbname
  replication {
    auto {}
  }
  depends_on = [google_project_service.apis]
}

resource "google_secret_manager_secret_version" "dbname" {
  secret      = google_secret_manager_secret.dbname.id
  secret_data = var.dbname
}
