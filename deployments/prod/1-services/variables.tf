variable "github_token" {
  description = "GitHub token with repo and workflow permissions"
  type        = string
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
variable "mcp_server_image_tag" {
  type    = string
  default = "latest"
}
