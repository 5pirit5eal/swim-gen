variable "project_id" {
  description = "The GCP Project ID"
  type        = string
}

variable "region" {
  description = "The GCP Region"
  type        = string
}

variable "apis" {
  description = "The GCP APIs to enable"
  type        = list(string)
}

variable "outputs_location" {
  description = "The GCP output locations"
  type        = string
  default     = null
}

variable "dbname" {
  description = "The GCP Cloud SQL Database Name"
  type        = string
}

variable "dbuser" {
  description = "The GCP Cloud SQL Database User"
  type        = string
}

variable "dbtier" {
  description = "The GCP Cloud SQL Database Tier"
  type        = string
}