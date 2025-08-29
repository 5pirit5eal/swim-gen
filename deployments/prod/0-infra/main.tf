resource "random_id" "db_name_suffix" {
  byte_length = 4
}

resource "google_sql_database_instance" "main" {
  name             = "main-instance-${random_id.db_name_suffix.hex}"
  database_version = "POSTGRES_17"
  region           = var.region
  project          = var.project_id
  root_password    = data.google_secret_manager_secret_version_access.dbpassword_root.secret_data

  settings {
    tier              = var.dbtier
    disk_size         = 10
    edition           = "ENTERPRISE"
    activation_policy = "ALWAYS"
    database_flags {
      name  = "max_connections"
      value = "200"
    }
    # enable_google_ml_integration = true
    password_validation_policy {
      min_length                  = 6
      reuse_interval              = 5
      complexity                  = "COMPLEXITY_DEFAULT"
      disallow_username_substring = true
      password_change_interval    = "30s"
      enable_password_policy      = true
    }
    availability_type = "ZONAL"
    ip_configuration {
      ssl_mode = "ENCRYPTED_ONLY"
      # Add optional authorized networks
      # Update to match the customer's networks
      #   authorized_networks {
      #     name  = "test-net-3"
      #     value = "203.0.113.0/24"
      #   }
      # Enable public IP
      ipv4_enabled = true
    }
    backup_configuration {
      enabled    = true
      start_time = "20:55"
      backup_retention_settings {
        retained_backups = 42
        retention_unit   = "COUNT"
      }
    }
  }
  deletion_protection = true
  depends_on          = [google_project_service.apis]
}

resource "google_sql_database" "main_db" {
  name            = var.dbname
  instance        = google_sql_database_instance.main.name
  deletion_policy = "ABANDON"
}

resource "google_sql_user" "dbuser" {
  name     = google_secret_manager_secret_version.dbuser.secret_data
  instance = google_sql_database_instance.main.name
  password = data.google_secret_manager_secret_version_access.dbpassword_user.secret_data
  password_policy {
    allowed_failed_attempts = 5
  }
}


# Storage bucket for the exported pdfs from the backend
resource "google_storage_bucket" "exported_pdfs" {
  name     = "${var.project_id}-swim-gen-exported-pdfs"
  location = var.region
  project  = var.project_id

  lifecycle {
    prevent_destroy = true
  }

  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "Delete"
    }
  }

  depends_on = [google_project_service.apis]
}
