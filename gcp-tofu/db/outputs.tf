locals {
  outputs_location = try(pathexpand(var.outputs_location), "")
  tfvars = {
    project_id = var.project_id
    region     = var.region
    csql_db = {
      id   = try(google_sql_database.swim-rag.id, null)
      name = var.dbname
      user = google_secret_manager_secret.dbuser.secret_data
      # password = data.google_secret_manager_secret_version_access.dbpassword.secret_data
    }
    csql_instance = {
      "connection_name" = try(google_sql_database_instance.main.connection_name, null)
      "uri"             = try(google_sql_database_instance.main.self_link, null)
      "public_ip"       = try(google_sql_database_instance.main.public_ip_address, null)
      "private_ip"      = try(google_sql_database_instance.main.private_ip_address, null)
    }
    secret_ids = [
      google_secret_manager_secret.dbuser.id,
      data.google_secret_manager_secret.dbpassword-root.id
    ]
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