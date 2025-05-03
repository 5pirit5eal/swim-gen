locals {
  iam_roles = [
    "roles/run.invoker",
    "roles/secretmanager.secretAccessor",
    "roles/cloudsql.client",
    "roles/cloudsql.editor",
    "roles/storage.admin"
  ]
}

resource "google_service_account" "service_account" {
  account_id                   = "swim-rag-cr-sa"
  display_name                 = "Swim RAG Cloud Run Service Account"
  project                      = var.project_id
  create_ignore_already_exists = true
}

resource "google_project_iam_member" "service_account_iam" {
  for_each = toset(local.iam_roles)
  project  = var.project_id
  role     = each.key
  member   = "serviceAccount:${google_service_account.service_account.email}"
}

resource "google_secret_manager_secret_iam_member" "sa_secret_access" {
  for_each  = toset(var.secret_ids)
  secret_id = each.key
  project   = var.project_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.service_account.email}"
}