variable "github_token" {
  description = "GitHub token with repo and workflow permissions"
  type        = string
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

variable "model" {
  description = "The model name"
  type        = string
}

variable "small_model" {
  description = "The small model name"
  type        = string
  default     = "gemini-2.5-flash-lite"
}

variable "embedding_name" {
  description = "The embedding name"
  type        = string
}
variable "embedding_drill_name" {
  description = "The embedding drill name"
  type        = string
}
variable "embedding_model" {
  description = "The embedding model"
  type        = string
}
variable "embedding_size" {
  description = "The embedding size"
  type        = number
}

variable "log_level" {
  description = "The log level"
  type        = string
}
variable "domain_url" {
  description = "The domain URL for the frontend"
  type        = string
}

variable "outputs_location" {
  description = "The GCP output locations"
  type        = string
  default     = "../0-config"
}

variable "backend_image_tag" {
  type    = string
  default = "latest"
}
variable "bff_image_tag" {
  type    = string
  default = "latest"
}
variable "frontend_image_tag" {
  type    = string
  default = "latest"
}
