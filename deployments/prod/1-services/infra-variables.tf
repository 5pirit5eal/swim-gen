variable "project_id" {
  description = "The GCP Project ID"
  type        = string
}

variable "region" {
  description = "The GCP Region"
  type        = string
}

variable "artifactregistry" {
  description = "Artifact Registry settings"
  type = object({
    location   = string
    repository = string
  })
}

variable "bucket_name" {
  description = "Exported PDFs bucket name"
  type        = string
}

variable "csql_db" {
  description = "Cloud SQL database properties"
  type = object({
    id     = string
    name   = string
    tier   = string
    dbuser = string
  })
}

variable "csql_instance" {
  description = "Cloud SQL instance properties"
  type = object({
    connection_name = string
    private_ip      = string
    public_ip       = string
    uri             = string
  })
}

variable "iam" {
  description = "Service account identities"
  type = object({
    github_actions    = object({ email = string, id = string })
    pdf_export        = object({ email = string, id = string })
    swim_gen_backend  = object({ email = string, id = string })
    swim_gen_frontend = object({ email = string, id = string })
  })
}

variable "secret_ids" {
  description = "Secret resource IDs"
  type = object({
    dbname          = string
    dbpassword_root = string
    dbpassword_user = string
    dbuser          = string
  })
}

variable "secret_version_ids" {
  description = "Secret version resource IDs"
  type = object({
    dbpassword_root = string
    dbpassword_user = string
  })
}
