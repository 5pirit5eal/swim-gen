locals {
  outputs_location = try(pathexpand(var.outputs_location), "")
  tfvars = {
    project_id = var.project_id
    region     = var.region
    csql_db = {
      id     = try(google_sql_database.main_db.id, null)
      name   = var.dbname
      tier   = var.dbtier
      dbuser = var.dbuser
    }
    csql_instance = {
      connection_name = try(google_sql_database_instance.main.connection_name, null)
      uri             = try(google_sql_database_instance.main.self_link, null)
      public_ip       = try(google_sql_database_instance.main.public_ip_address, null)
      private_ip      = try(google_sql_database_instance.main.private_ip_address, null)
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
