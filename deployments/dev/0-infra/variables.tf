variable "project_id" {
  description = "The GCP Project ID"
  type        = string
}

variable "region" {
  description = "The GCP Region"
  type        = string
}

variable "github_token" {
  description = "GitHub token with repo and workflow permissions"
  type        = string
  sensitive   = true
}

variable "github_owner" {
  description = "GitHub organization or user that owns the repository"
  type        = string
  default     = "5pirit5eal"
}

variable "github_repository" {
  description = "Short repository name (without owner)"
  type        = string
  default     = "swim-gen"
}

variable "supabase_access_token" {
  description = "Supabase access token with full permissions"
  type        = string
  sensitive   = true
}

variable "apis" {
  description = "The GCP APIs to enable"
  type        = list(string)
}

variable "outputs_location" {
  description = "The GCP output locations"
  type        = string
  default     = "../0-config"
}

variable "dbusers" {
  description = "The SQL Database Users"
  type = object({
    backend  = string
    frontend = string
  })
}

variable "supabase" {
  description = "Supabase project configuration"
  type = object({
    organization_id = string
    name            = string
    region          = string
  })
}
