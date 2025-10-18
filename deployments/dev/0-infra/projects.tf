resource "supabase_project" "development" {
  organization_id   = var.supabase.organization_id
  name              = var.supabase.name
  database_password = data.google_secret_manager_secret_version_access.dbpassword_root.secret_data
  region            = var.supabase.region
  lifecycle {
    ignore_changes = [database_password]
  }
}

resource "supabase_settings" "development" {
  project_ref = supabase_project.development.id
  api = jsonencode({
    db_schema            = "public,storage,graphql_public"
    db_extra_search_path = "public,extensions"
    max_rows             = 1000
  })
}

data "google_project" "project" {
  project_id = var.project_id
}
