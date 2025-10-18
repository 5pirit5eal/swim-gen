locals {
  outputs_location = try(pathexpand(var.outputs_location), "")
  tfvars = {
    project_id = var.project_id
    region     = var.region
    supabase_pool = {
      id            = supabase_project.production.id
      name          = var.supabase.name
      backend_user  = "${var.dbusers.backend}.${supabase_project.production.id}"
      frontend_user = "${var.dbusers.frontend}.${supabase_project.production.id}"
      host          = "aws-1-${var.supabase.region}.pooler.supabase.com"
      port          = 6543
      dbname        = "postgres"
    }
    secret_ids         = local.secret_ids
    secret_version_ids = local.secret_version_ids
    bucket_name        = google_storage_bucket.exported_pdfs.name
    artifactregistry = {
      repository = google_artifact_registry_repository.docker.name
      location   = google_artifact_registry_repository.docker.location
    }
    iam = {
      github_actions = {
        email = try(google_service_account.github_actions_sa.email, null)
        id    = try(google_service_account.github_actions_sa.id, null)
      }
      swim_gen_backend = {
        email = google_service_account.swim_gen_backend_sa.email,
        id    = google_service_account.swim_gen_backend_sa.id,
      }
      swim_gen_frontend = {
        email = google_service_account.swim_gen_frontend_sa.email,
        id    = google_service_account.swim_gen_frontend_sa.id,
      }
      pdf_export = {
        email = google_service_account.pdf_export_sa.email
        id    = google_service_account.pdf_export_sa.id
      }
    }
  }
}


output "tfvars" {
  description = "Terraform variable files for the following stages."
  value       = local.tfvars
}


resource "local_file" "tfvars" {
  for_each        = var.outputs_location == null ? {} : { 1 = 1 }
  file_permission = "0644"
  filename        = "${local.outputs_location}/infra.auto.tfvars.json"
  content         = jsonencode(local.tfvars)
}
