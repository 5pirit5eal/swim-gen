# Variables from the db stage

variable "csql_instance" {
  type = optional(object({
    connection_name = optional(string)
    public_ip       = optional(string)
    private_ip      = optional(string)
    uri             = optional(string)
  }))
  description = "Database instance details including connection name, IP, and URI"
}

variable "csql_db" {
  type = object({
    name     = string
    user     = string
    password = optional(string)
  })
  description = "Name of the database"
}

variable "project_id" {
  type        = string
  description = "GCP project ID"
}

variable "region" {
  type        = string
  description = "GCP region"
}

variable "secret_ids" {
  type        = list(string)
  description = "List of secret IDs"
}