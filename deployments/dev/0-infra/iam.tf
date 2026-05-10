# Description: This file is responsible for creating the service account and assigning the roles to the service account.
resource "google_project_service" "apis" {
  for_each = toset(var.apis)
  project  = var.project_id
  service  = each.key
}

# Github Actions Service Account
resource "google_service_account" "github_actions_sa" {
  account_id                   = "github-actions-sa"
  display_name                 = "Github Actions Service Account"
  project                      = var.project_id
  create_ignore_already_exists = true
}

resource "google_project_iam_member" "github_actions_iam" {
  for_each = toset([
    "roles/iam.workloadIdentityUser",
    "roles/storage.admin",
    "roles/run.admin",
    "roles/logging.logWriter",
    "roles/iam.serviceAccountUser",
    "roles/iam.serviceAccountTokenCreator",
    "roles/artifactregistry.admin",
    "roles/secretmanager.secretAccessor",
    "roles/aiplatform.user",
    # Roles for Terraform to plan and apply
    "roles/serviceusage.serviceUsageAdmin",
    "roles/resourcemanager.projectIamAdmin",
    "roles/iam.serviceAccountAdmin",
    "roles/iam.workloadIdentityPoolAdmin",
    "roles/secretmanager.viewer"
  ])
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

resource "google_secret_manager_secret_iam_member" "github_actions_sa_secret_access" {
  for_each = local.secret_ids
  project  = data.google_project.project.number
  # Extract the short secret_id from the full resource id
  secret_id = each.value
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

# Backend Service Account
resource "google_service_account" "swim_gen_backend_sa" {
  account_id                   = "swim-gen-backend-sa"
  display_name                 = "Swim Gen Backend Service Account"
  project                      = var.project_id
  create_ignore_already_exists = true
}

resource "google_project_iam_member" "swim_gen_backend_iam" {
  for_each = toset([
    "roles/secretmanager.secretAccessor",
    "roles/storage.admin",
    "roles/aiplatform.user",
    "roles/iam.serviceAccountTokenCreator",
    "roles/monitoring.metricWriter",
    "roles/cloudtrace.agent",
  ])
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.swim_gen_backend_sa.email}"
}

resource "google_secret_manager_secret_iam_member" "swim_gen_backend_sa_secret_access" {
  for_each  = local.secret_ids
  secret_id = each.value
  project   = data.google_project.project.number
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.swim_gen_backend_sa.email}"
}

# Frontend Service Account
resource "google_service_account" "swim_gen_frontend_sa" {
  account_id                   = "swim-gen-frontend-sa"
  display_name                 = "Swim Gen Frontend Service Account"
  project                      = var.project_id
  create_ignore_already_exists = true
}

resource "google_service_account_iam_member" "swim_gen_frontend_token_creator_self" {
  service_account_id = google_service_account.swim_gen_frontend_sa.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${google_service_account.swim_gen_frontend_sa.email}"
}

resource "google_project_iam_member" "swim_gen_frontend_iam" {
  for_each = toset([
    "roles/monitoring.metricWriter",
    "roles/cloudtrace.agent",
  ])
  project = var.project_id
  role    = each.key
  member  = "serviceAccount:${google_service_account.swim_gen_frontend_sa.email}"
}

# Make the Github Actions service account a user of the Cloud Run service accounts
resource "google_service_account_iam_member" "github_actions_sa_user_backend" {
  service_account_id = google_service_account.swim_gen_backend_sa.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

resource "google_service_account_iam_member" "github_actions_sa_user_frontend" {
  service_account_id = google_service_account.swim_gen_frontend_sa.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.github_actions_sa.email}"
}

# Sign PDF Service Account
resource "google_service_account" "pdf_export_sa" {
  account_id                   = "pdf-export-sa"
  display_name                 = "Sign PDF Service Account"
  project                      = var.project_id
  create_ignore_already_exists = true
}


# Allow the Backend service account to impersonate the PDF export service account
resource "google_service_account_iam_member" "pdf_export_sa_user" {
  service_account_id = google_service_account.pdf_export_sa.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.swim_gen_backend_sa.email}"
}

# Give the PDF export service account access to the storage bucket
resource "google_storage_bucket_iam_member" "pdf_export_sa_access" {
  for_each = toset([
    "roles/storage.objectAdmin",
    "roles/storage.objectViewer",
  ])
  bucket = google_storage_bucket.exported_pdfs.name
  role   = each.key
  member = "serviceAccount:${google_service_account.pdf_export_sa.email}"
}

# =============================================================================
# Observability Analytics IAM
# =============================================================================
# Grants the deploying user (identified via gcloud application-default login)
# the roles required to:
#   - Query log views in the custom log bucket via Observability Analytics
#   - Query trace spans via the _Trace._AllSpans view
#   - Run and save SQL queries, and pin charts to Cloud Monitoring dashboards
#   - Read linked BigQuery datasets from Looker Studio
#
# roles/observability.viewAccessor  — read access to observability views (logs + traces)
# roles/observability.analyticsUser — run/save queries in Observability Analytics UI
# roles/logging.viewAccessor        — read access to log views in the log bucket
# roles/bigquery.dataViewer         — read linked BQ datasets from Looker Studio / BigQuery Studio

data "google_client_openid_userinfo" "me" {}

resource "google_project_iam_member" "observability_analytics" {
  for_each = toset([
    "roles/observability.viewAccessor",
    "roles/observability.analyticsUser",
    "roles/logging.viewAccessor",
    "roles/bigquery.dataViewer",
    "roles/bigquery.jobUser",
  ])
  project = var.project_id
  role    = each.key
  member  = "user:${data.google_client_openid_userinfo.me.email}"
}
