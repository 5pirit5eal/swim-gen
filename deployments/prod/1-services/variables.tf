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

variable "version_tag" {
  description = "The version tag for all images"
  type        = string
  default     = "latest"
}
variable "outputs_location" {
  description = "The GCP output locations"
  type        = string
  default     = "../0-config"
}
