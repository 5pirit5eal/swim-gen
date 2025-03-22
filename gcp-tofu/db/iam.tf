# Description: This file is responsible for creating the service account and assigning the roles to the service account.
resource "google_project_service" "apis" {
  for_each = toset(var.apis)
  project  = var.project_id
  service  = each.key
}

