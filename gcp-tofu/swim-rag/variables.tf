variable "project_id" {
  description = "The GCP Project ID"
  type        = string
}

variable "region" {
  description = "The GCP Region"
  type        = string
}

variable "outputs_location" {
  description = "The GCP output locations"
  type        = string
  default     = null
}



