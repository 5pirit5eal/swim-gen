locals {
  outputs_location = try(pathexpand(var.outputs_location), "")
  tfvars = {
    project_id = var.project_id
    region     = var.region
    csql_db = {
      id   = try(google_sql_database.main_db.id, null)
      name = var.dbname
      tier = var.dbtier
    }
    csql_instance = {
      "connection_name" = try(google_sql_database_instance.main.connection_name, null)
      "uri"             = try(google_sql_database_instance.main.self_link, null)
      "public_ip"       = try(google_sql_database_instance.main.public_ip_address, null)
      "private_ip"      = try(google_sql_database_instance.main.private_ip_address, null)
    }
    secret_ids  = local.secret_ids
    bucket_name = google_storage_bucket.exported_pdfs.name
  }
}


output "tfvars" {
  description = "Terraform variable files for the following stages."
  value       = local.tfvars
}


resource "local_file" "tfvars" {
  for_each        = var.outputs_location == null ? {} : { 1 = 1 }
  file_permission = "0644"
  filename        = "${local.outputs_location}/db.auto.tfvars.json"
  content         = jsonencode(local.tfvars)
}