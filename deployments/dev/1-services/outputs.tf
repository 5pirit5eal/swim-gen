locals {
  outputs_location = try(pathexpand(var.outputs_location), "")
  tfvars = {
    backend_url  = try(google_cloud_run_v2_service.backend.uri, null)
    bff_url      = try(google_cloud_run_v2_service.bff.uri, null)
    frontend_url = try(google_cloud_run_v2_service.frontend.uri, null)
  }
}

output "tfvars" {
  description = "Terraform variable files for the following stages."
  value       = local.tfvars
}

resource "local_file" "tfvars" {
  for_each        = var.outputs_location == null ? {} : { 1 = 1 }
  file_permission = "0644"
  filename        = "${local.outputs_location}/services.auto.tfvars.json"
  content         = jsonencode(local.tfvars)
}
