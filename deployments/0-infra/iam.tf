# Description: This file is responsible for creating the service account and assigning the roles to the service account.
resource "google_project_service" "apis" {
  for_each = toset(var.apis)
  project  = var.project_id
  service  = each.key
}

# Cloud Build Service Account
resource "google_service_account" "cloud_build_sa" {
  account_id                   = "cloud-build-sa"
  display_name                 = "Cloud Build Service Account"
  project                      = var.project_id
  create_ignore_already_exists = true
}

resource "google_project_iam_member" "cloud_build_iam" {
  for_each = toset([
    "roles/cloudbuild.builds.editor",
    "roles/storage.admin",
    "roles/run.developer",
    "roles/logging.logWriter",
    "roles/iam.serviceAccountUser",
    "roles/iam.serviceAccountTokenCreator",
    "roles/artifactregistry.writer",
    "roles/secretmanager.secretAccessor"
  ])
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.cloud_build_sa.email}"
}

resource "google_secret_manager_secret_iam_member" "cloud_build_sa_secret_access" {
  for_each  = local.secret_ids
  secret_id = each.value
  project   = var.project_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_build_sa.email}"
}

# Cloud Run Service Account
resource "google_service_account" "cloud_run_sa" {
  account_id                   = "cloud-run-sa"
  display_name                 = "Cloud Run Service Account"
  project                      = var.project_id
  create_ignore_already_exists = true
}

resource "google_project_iam_member" "service_account_iam" {
  for_each = toset([
    "roles/secretmanager.secretAccessor",
    "roles/cloudsql.client",
    "roles/cloudsql.editor",
    "roles/storage.admin"
  ])
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.cloud_run_sa.email}"
}

resource "google_secret_manager_secret_iam_member" "cloud_run_sa_secret_access" {
  for_each  = local.secret_ids
  secret_id = each.value
  project   = var.project_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_run_sa.email}"
}

# Make the Cloud Build service account a user of the Cloud Run service account
resource "google_service_account_iam_member" "cloud_build_sa_user" {
  service_account_id = google_service_account.cloud_run_sa.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.cloud_build_sa.email}"
}

# Make the default cloud build service account accessor of the github token secret
resource "google_secret_manager_secret_iam_member" "cloud_build_sa_github_token_access" {
  secret_id = google_secret_manager_secret.github_token_secret.id
  project   = var.project_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-cloudbuild.iam.gserviceaccount.com"
}